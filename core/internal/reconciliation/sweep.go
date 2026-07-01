package reconciliation

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Sweeper interface {
	RunSweep(ctx context.Context) error
}

type sweeper struct {
	q        sqlc.Querier
	idemRepo IdempotencyRepository
	client   *nomba.Client
}

func NewSweeper(q sqlc.Querier, idemRepo IdempotencyRepository, client *nomba.Client) Sweeper {
	return &sweeper{q: q, idemRepo: idemRepo, client: client}
}

func (s *sweeper) RunSweep(ctx context.Context) error {
	log.Println("Running reconciliation sweep...")

	accounts, err := s.q.ListActiveVirtualAccounts(ctx)
	if err != nil {
		return err
	}

	dateTo := time.Now()
	// Reconcile the last 2 hours (covers the max 53m retry delay of nomba + buffer)
	dateFrom := dateTo.Add(-2 * time.Hour) 

	for _, acc := range accounts {
		if !acc.AccountNumber.Valid {
			continue
		}
		
		txns, err := s.client.FetchTransactions(ctx, acc.AccountNumber.String, dateFrom, dateTo)
		if err != nil {
			log.Printf("Failed to fetch transactions for %s: %v", acc.AccountNumber.String, err)
			continue
		}

		for _, txn := range txns {
			if txn.Type == "vact_transfer" && txn.EntryType == "CREDIT" && txn.RecipientAccountType == "VIRTUAL" {
				amountFloat, err := strconv.ParseFloat(txn.Amount, 64)
				if err != nil {
					log.Printf("Invalid amount %s: %v", txn.Amount, err)
					continue
				}

				amountKobo := int64(math.Round(amountFloat * 100))

				_, _ = s.idemRepo.CheckOrSetIdempotency(ctx, txn.ID, "sweep", func() (uuid.UUID, error) {
					entryID := uuid.New()
					
					var occurredAt time.Time
					if t, err := time.Parse(time.RFC3339, txn.TimeCreated); err == nil {
						occurredAt = t
					} else {
						occurredAt = time.Now()
					}

					var pgNarration pgtype.Text
					if txn.Narration != "" {
						pgNarration = pgtype.Text{String: txn.Narration, Valid: true}
					}

					var pgSenderName pgtype.Text
					if txn.SenderName != "" {
						pgSenderName = pgtype.Text{String: txn.SenderName, Valid: true}
					}

					_, err := s.q.InsertLedgerEntry(ctx, sqlc.InsertLedgerEntryParams{
						ID:               entryID,
						VirtualAccountID: acc.ID,
						IdentityID:       acc.IdentityID,
						AmountKobo:       amountKobo,
						Direction:        "inbound",
						Status:           "matched",
						NombaReference:   txn.ID,
						Source:           "sweep",
						Narration:        pgNarration,
						SenderName:       pgSenderName,
						OccurredAt:       occurredAt,
					})
					return entryID, err
				})
			}
		}
	}
	
	return nil
}
