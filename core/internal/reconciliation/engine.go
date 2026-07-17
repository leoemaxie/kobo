package reconciliation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/billing"
	"github.com/leoemaxie/kobo/internal/monnify"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Engine interface {
	ProcessWebhook(ctx context.Context, payload *monnify.WebhookPayload) error
}

type MonnifyTransactionFetcher interface {
	FetchSingleTransaction(ctx context.Context, transactionRef string) (*monnify.TransactionResult, error)
}

type engine struct {
	q           sqlc.Querier
	idemRepo    IdempotencyRepository
	recorder    *billing.UsageRecorder
	monnifyClient MonnifyTransactionFetcher
}

func NewEngine(q sqlc.Querier, idemRepo IdempotencyRepository, recorder *billing.UsageRecorder, monnifyClient MonnifyTransactionFetcher) Engine {
	return &engine{q: q, idemRepo: idemRepo, recorder: recorder, monnifyClient: monnifyClient}
}

func (e *engine) ProcessWebhook(ctx context.Context, payload *monnify.WebhookPayload) error {
	if payload.EventType != "payment_success" {
		return nil // Ignore other events
	}

	if payload.Data.Transaction.Type == "online_checkout" {
		return e.handleCheckoutWebhook(ctx, payload)
	}

	if payload.Data.Transaction.AliasAccountType != "VIRTUAL" {
		return nil // Ignore non-virtual account payments
	}

	accountNumStr := payload.Data.Transaction.AliasAccountNumber
	amountKobo := int64(math.Round(payload.Data.Transaction.TransactionAmount * 100))

	var pgAccountNum pgtype.Text
	if accountNumStr != "" {
		pgAccountNum = pgtype.Text{String: accountNumStr, Valid: true}
	}

	account, err := e.q.GetVirtualAccountByAccountNumber(ctx, pgAccountNum)
	if err != nil {
		return fmt.Errorf("virtual account not found: %s", accountNumStr)
	}

	transactionID := payload.Data.Transaction.TransactionID

	if e.monnifyClient != nil {
		txn, err := e.monnifyClient.FetchSingleTransaction(ctx, transactionID)
		if err != nil {
			return fmt.Errorf("failed to verify transaction %s with Monnify: %w", transactionID, err)
		}
		if txn.Status != "SUCCESS" && txn.Status != "PAYMENT_SUCCESSFUL" {
			return fmt.Errorf("transaction %s is not successful according to Monnify, status: %s", transactionID, txn.Status)
		}
	}

	// Idempotency check and insert
	isDuplicate, err := e.idemRepo.CheckOrSetIdempotency(ctx, transactionID, "webhook", func() (uuid.UUID, error) {
		entryID := uuid.New()

		var occurredAt time.Time
		if t, err := time.Parse(time.RFC3339, payload.Data.Transaction.Time); err == nil {
			occurredAt = t
		} else {
			occurredAt = time.Now()
		}

		narrationStr := payload.Data.Transaction.Narration
		var pgNarration pgtype.Text
		if narrationStr != "" {
			pgNarration = pgtype.Text{String: narrationStr, Valid: true}
		}

		senderNameStr := payload.Data.Customer.SenderName
		var pgSenderName pgtype.Text
		if senderNameStr != "" {
			pgSenderName = pgtype.Text{String: senderNameStr, Valid: true}
		}

		_, err := e.q.InsertLedgerEntry(ctx, sqlc.InsertLedgerEntryParams{
			ID:               entryID,
			VirtualAccountID: account.ID,
			IdentityID:       account.IdentityID,
			AmountKobo:       amountKobo,
			Direction:        "inbound",
			Status:           "matched",
			MonnifyReference:   transactionID,
			Source:           "webhook",
			Narration:        pgNarration,
			SenderName:       pgSenderName,
			OccurredAt:       occurredAt,
		})
		return entryID, err
	})

	if err != nil {
		return err
	}

	if isDuplicate {
		return nil // Already processed, return success
	}

	ident, err := e.q.GetIdentityByVirtualAccountID(ctx, account.ID)
	if err == nil {
		env := sqlc.ConsoleEnvironmentProduction
		if ident.CredentialEnvironment == sqlc.ConsoleEnvironmentSandbox {
			env = sqlc.ConsoleEnvironmentSandbox
		}
		e.recorder.RecordAsync(ident.IntegratorID, env, "transaction_processed", transactionID, 200) // ₦2 fee
	}

	return nil
}

func (e *engine) handleCheckoutWebhook(ctx context.Context, payload *monnify.WebhookPayload) error {
	if payload.Data.TokenizedCardData != nil && payload.Data.TokenizedCardData.TokenKey != "" && payload.Data.TokenizedCardData.TokenKey != "N/A" {
		integratorID, err := uuid.Parse(payload.Data.Order.CustomerId)
		if err != nil {
			return fmt.Errorf("invalid customer ID (integrator ID) in webhook: %w", err)
		}

		// 1. Unset any existing default payment methods for this integrator
		if err := e.q.UnsetDefaultPaymentMethods(ctx, integratorID); err != nil {
			return fmt.Errorf("failed to unset default payment methods: %w", err)
		}

		cardLast4 := ""
		if len(payload.Data.TokenizedCardData.CardPan) >= 4 {
			pan := payload.Data.TokenizedCardData.CardPan
			cardLast4 = pan[len(pan)-4:]
		}

		// 2. Insert the new payment method
		_, err = e.q.InsertPaymentMethod(ctx, sqlc.InsertPaymentMethodParams{
			IntegratorID:  integratorID,
			MonnifyTokenKey: payload.Data.TokenizedCardData.TokenKey,
			CardLast4:     pgtype.Text{String: cardLast4, Valid: cardLast4 != ""},
			CardBrand:     pgtype.Text{String: payload.Data.TokenizedCardData.CardType, Valid: payload.Data.TokenizedCardData.CardType != ""},
			IsDefault:     true,
		})
		if err != nil {
			return fmt.Errorf("failed to save payment method: %w", err)
		}
	}
	return nil
}
