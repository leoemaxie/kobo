package reconciliation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Engine interface {
	ProcessWebhook(ctx context.Context, payload *nomba.WebhookPayload) error
}

type engine struct {
	q      sqlc.Querier
	idemRepo IdempotencyRepository
}

func NewEngine(q sqlc.Querier, idemRepo IdempotencyRepository) Engine {
	return &engine{q: q, idemRepo: idemRepo}
}

func (e *engine) ProcessWebhook(ctx context.Context, payload *nomba.WebhookPayload) error {
	if payload.EventType != "payment_success" {
		return nil // Ignore other events
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
		// Log error, account not found
		return fmt.Errorf("virtual account not found: %s", accountNumStr)
	}

	transactionID := payload.Data.Transaction.TransactionID

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
			NombaReference:   transactionID,
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

	return nil
}
