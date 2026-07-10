package payout

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

const (
	MinPayoutKobo         int64 = 100_000 // 1,000 NGN
	TransferFeeBufferKobo int64 = 5_000   // 50 NGN buffer
	PlatformFeeKobo       int64 = 0
)

var (
	ErrNoBankAccount       = errors.New("no active bank account configured")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrBelowMinimum        = errors.New("amount below minimum threshold")
	ErrPayoutInProgress    = errors.New("a payout is already in progress")
)

type Service struct {
	pool        *pgxpool.Pool
	q           sqlc.Querier
	nombaClient *nomba.Client
}

func NewService(pool *pgxpool.Pool, q sqlc.Querier, nombaClient *nomba.Client) *Service {
	return &Service{
		pool:        pool,
		q:           q,
		nombaClient: nombaClient,
	}
}

func (s *Service) SaveBankAccount(ctx context.Context, integratorID, userID uuid.UUID, accountNumber, bankCode, bankName string) (sqlc.ConsolePayoutBankAccount, error) {
	// 1. Verify via Nomba
	lookupResp, err := s.nombaClient.LookupBankAccount(ctx, accountNumber, bankCode)
	if err != nil {
		return sqlc.ConsolePayoutBankAccount{}, fmt.Errorf("bank account verification failed: %w", err)
	}

	var savedAccount sqlc.ConsolePayoutBankAccount

	// 2. Transaction to deactivate existing and insert new
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		qtx := s.q.(*sqlc.Queries).WithTx(tx)

		err := qtx.DeactivatePayoutBankAccounts(ctx, integratorID)
		if err != nil {
			return err
		}

		savedAccount, err = qtx.InsertPayoutBankAccount(ctx, sqlc.InsertPayoutBankAccountParams{
			IntegratorID:  integratorID,
			AccountNumber: accountNumber,
			AccountName:   lookupResp.AccountName,
			BankCode:      bankCode,
			BankName:      bankName,
			CreatedBy:     userID,
		})
		return err
	})

	return savedAccount, err
}

func (s *Service) GetBankAccount(ctx context.Context, integratorID uuid.UUID) (sqlc.ConsolePayoutBankAccount, error) {
	account, err := s.q.GetActivePayoutBankAccount(ctx, integratorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.ConsolePayoutBankAccount{}, ErrNoBankAccount
		}
		return sqlc.ConsolePayoutBankAccount{}, err
	}
	return account, nil
}

func (s *Service) RequestPayout(ctx context.Context, integratorID, initiatedByUserID uuid.UUID, requestedAmountKobo int64) (sqlc.ConsolePayout, error) {
	if requestedAmountKobo < MinPayoutKobo {
		return sqlc.ConsolePayout{}, ErrBelowMinimum
	}

	bankAccount, err := s.GetBankAccount(ctx, integratorID)
	if err != nil {
		return sqlc.ConsolePayout{}, err
	}

	inProgressCount, err := s.q.CountInProgressPayouts(ctx, integratorID)
	if err != nil {
		return sqlc.ConsolePayout{}, err
	}
	if inProgressCount > 0 {
		return sqlc.ConsolePayout{}, ErrPayoutInProgress
	}

	var payout sqlc.ConsolePayout
	var totalDeduction int64
	var merchantTxRef string
	var netAmountKobo int64

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return sqlc.ConsolePayout{}, fmt.Errorf("failed to start tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.(*sqlc.Queries).WithTx(tx)

	_, err = qtx.LockIntegratorRow(ctx, integratorID)
	if err != nil {
		return sqlc.ConsolePayout{}, fmt.Errorf("failed to lock integrator row: %w", err)
	}

	totalDeduction = requestedAmountKobo + TransferFeeBufferKobo + PlatformFeeKobo

	_, err = qtx.DeductIntegratorBalance(ctx, sqlc.DeductIntegratorBalanceParams{
		ID:                integratorID,
		WalletBalanceKobo: totalDeduction,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.ConsolePayout{}, ErrInsufficientBalance
		}
		return sqlc.ConsolePayout{}, err
	}

	netAmountKobo = requestedAmountKobo - PlatformFeeKobo
	merchantTxRef = "payout_" + uuid.New().String()

	payout, err = qtx.CreatePayout(ctx, sqlc.CreatePayoutParams{
		IntegratorID:          integratorID,
		BankAccountID:         bankAccount.ID,
		RequestedAmountKobo:   requestedAmountKobo,
		PlatformFeeKobo:       PlatformFeeKobo,
		TransferFeeBufferKobo: TransferFeeBufferKobo,
		NetAmountKobo:         netAmountKobo,
		MerchantTxRef:         merchantTxRef,
		InitiatedBy:           initiatedByUserID,
	})
	if err != nil {
		return sqlc.ConsolePayout{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.ConsolePayout{}, fmt.Errorf("failed to commit tx: %w", err)
	}

	req := nomba.TransferToBankRequest{
		Amount:        float64(netAmountKobo) / 100,
		AccountNumber: bankAccount.AccountNumber,
		AccountName:   bankAccount.AccountName,
		BankCode:      bankAccount.BankCode,
		MerchantTxRef: merchantTxRef,
		SenderName:    "Kobo",
		Narration:     "Kobo Payout",
	}

	resp, transferErr := s.nombaClient.TransferToBank(ctx, req)

	if transferErr != nil {
		// REVERSAL
		log.Printf("payout request to nomba failed, reversing... %v", transferErr)
		err := s.q.CreditIntegratorBalance(ctx, sqlc.CreditIntegratorBalanceParams{
			ID:                integratorID,
			WalletBalanceKobo: totalDeduction,
		})
		if err != nil {
			log.Printf("CRITICAL: failed to credit integrator balance after nomba transfer error: %v", err)
		}

		failureReasonStr := transferErr.Error()
		failureReason := pgtype.Text{String: failureReasonStr, Valid: true}
		payout, _ = s.q.UpdatePayoutStatus(ctx, sqlc.UpdatePayoutStatusParams{
			ID:            payout.ID,
			Status:        "failed",
			FailureReason: failureReason,
		})
		return payout, transferErr
	}

	var newStatus string
	var nombaTransferID pgtype.Text
	var actualFee pgtype.Int8

	if resp.IsAsync {
		newStatus = "processing"
	} else {
		newStatus = "successful"
		nombaTransferID = pgtype.Text{String: resp.NombaTransferID, Valid: true}
		feeKobo := int64(resp.FeeNaira * 100)
		actualFee = pgtype.Int8{Int64: feeKobo, Valid: true}
	}

	payout, err = s.q.UpdatePayoutStatus(ctx, sqlc.UpdatePayoutStatusParams{
		ID:                    payout.ID,
		Status:                newStatus,
		NombaTransferID:       nombaTransferID,
		ActualTransferFeeKobo: actualFee,
	})

	if err != nil {
		log.Printf("failed to update payout status: %v", err)
	}

	return payout, nil
}

func (s *Service) HandleTransferWebhook(ctx context.Context, merchantTxRef, nombaStatus, nombaTransferID string) error {
	payout, err := s.q.GetPayoutByMerchantTxRef(ctx, merchantTxRef)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // not a payout transfer
		}
		return err
	}

	if payout.Status == "successful" || payout.Status == "failed" {
		return nil
	}

	if nombaStatus == "SUCCESS" {
		_, err := s.q.UpdatePayoutStatus(ctx, sqlc.UpdatePayoutStatusParams{
			ID:              payout.ID,
			Status:          "successful",
			NombaTransferID: pgtype.Text{String: nombaTransferID, Valid: true},
		})
		return err
	}

	// Any non-SUCCESS is considered failed
	totalDeduction := payout.RequestedAmountKobo + payout.TransferFeeBufferKobo + payout.PlatformFeeKobo

	// Reversal
	err = s.q.CreditIntegratorBalance(ctx, sqlc.CreditIntegratorBalanceParams{
		ID:                payout.IntegratorID,
		WalletBalanceKobo: totalDeduction,
	})
	if err != nil {
		log.Printf("CRITICAL: failed to credit integrator balance on webhook reversal: %v", err)
	}

	_, err = s.q.UpdatePayoutStatus(ctx, sqlc.UpdatePayoutStatusParams{
		ID:            payout.ID,
		Status:        "failed",
		FailureReason: pgtype.Text{String: nombaStatus, Valid: true},
	})
	return err
}

func (s *Service) ListPayouts(ctx context.Context, integratorID uuid.UUID, limit, offset int32) ([]sqlc.ConsolePayout, error) {
	return s.q.ListPayoutsForIntegrator(ctx, sqlc.ListPayoutsForIntegratorParams{
		IntegratorID: integratorID,
		Limit:        limit,
		Offset:       offset,
	})
}
