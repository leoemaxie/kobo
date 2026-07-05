package reconciliation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/stretchr/testify/assert"
)

func TestRunSweep_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/token/issue", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"code": "00",
			"description": "Success",
			"data": {
				"access_token": "mock-token",
				"refresh_token": "mock-refresh",
				"expiresAt": "2030-01-01T00:00:00Z"
			}
		}`))
	})

	mux.HandleFunc("/transactions/virtual", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"code": "00",
			"description": "Success",
			"data": {
				"results": [
					{
						"id": "txn-sweep-1",
						"amount": "200.50",
						"type": "vact_transfer",
						"entryType": "CREDIT",
						"recipientAccountType": "VIRTUAL",
						"timeCreated": "2026-07-01T10:00:00Z"
					}
				],
				"cursor": ""
			}
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := nomba.NewClient(server.URL, "client", "secret", "account", server.Client())

	accID := uuid.New()

	mq := &mockQuerier{
		ListActiveVirtualAccountsFunc: func(ctx context.Context) ([]sqlc.VirtualAccount, error) {
			return []sqlc.VirtualAccount{
				{
					ID:            accID,
					AccountNumber: pgtype.Text{String: "1234567890", Valid: true},
				},
			}, nil
		},
		InsertLedgerEntryFunc: func(ctx context.Context, arg sqlc.InsertLedgerEntryParams) (sqlc.LedgerEntry, error) {
			assert.Equal(t, accID, arg.VirtualAccountID)
			assert.Equal(t, int64(20050), arg.AmountKobo) // 200.50 * 100
			assert.Equal(t, "txn-sweep-1", arg.NombaReference)
			assert.Equal(t, "sweep", arg.Source)
			return sqlc.LedgerEntry{ID: arg.ID}, nil
		},
	}

	mIdem := &mockIdempotencyRepo{
		CheckOrSetIdempotencyFunc: func(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error) {
			assert.Equal(t, "txn-sweep-1", reference)
			assert.Equal(t, "sweep", source)
			_, err := insertLedgerFunc()
			return false, err
		},
	}

	sw := NewSweeper(mq, mIdem, client)

	err := sw.RunSweep(context.Background())
	assert.NoError(t, err)
}
