package billing

import (
	"context"
	"log"
	"time"

	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type InvoiceJob struct {
	q           sqlc.Querier
	nombaClient *nomba.Client
}

func NewInvoiceJob(q sqlc.Querier, nombaClient *nomba.Client) *InvoiceJob {
	return &InvoiceJob{q: q, nombaClient: nombaClient}
}

// Run should be called periodically (e.g. daily) by the worker
func (j *InvoiceJob) Run(ctx context.Context) error {
	log.Println("Starting automated billing job...")

	// 1. Generate new invoices for the current period
	currentPeriod := time.Now().Format("2006-01") // e.g. "2026-07"
	err := j.q.GenerateMonthlyInvoices(ctx, currentPeriod)
	if err != nil {
		log.Printf("failed to generate monthly invoices: %v", err)
	}

	// 2. Process pending/failed invoices that are due for retry
	invoices, err := j.q.GetPendingInvoices(ctx)
	if err != nil {
		log.Printf("failed to get pending invoices: %v", err)
		return err
	}

	for _, inv := range invoices {
		// Try wallet deduction first
		walletBalance, err := j.q.GetIntegratorWalletBalance(ctx, inv.IntegratorID)
		if err == nil && walletBalance > 0 {
			deduction := walletBalance
			if deduction >= inv.AmountKobo {
				deduction = inv.AmountKobo
			}

			// Deduct from wallet
			j.q.UpdateIntegratorWalletBalance(ctx, sqlc.UpdateIntegratorWalletBalanceParams{
				ID:                inv.IntegratorID,
				WalletBalanceKobo: -deduction, // Negative to subtract
			})

			inv.AmountKobo -= deduction

			if inv.AmountKobo <= 0 {
				log.Printf("invoice %s fully paid via wallet balance", inv.ID)
				j.q.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
					ID:         inv.ID,
					Status:     "paid",
					RetryCount: inv.RetryCount,
				})
				continue
			}
		}

		pm, err := j.q.GetDefaultPaymentMethod(ctx, inv.IntegratorID)
		if err != nil || pm.NombaTokenKey == "" {
			log.Printf("no default payment method for integrator %s", inv.IntegratorID)
			j.failInvoice(ctx, inv)
			continue
		}

		resp, err := j.nombaClient.ChargeToken(ctx, nomba.ChargeTokenRequest{
			TokenKey: pm.NombaTokenKey,
			Order: nomba.OrderInfo{
				Amount:         float64(inv.AmountKobo) / 100.0,
				Currency:       "NGN",
				OrderReference: "inv_" + inv.ID.String(),
				CustomerEmail:  "billing@yourdomain.com", // Required by Nomba API
				CallbackUrl:    "https://api.yourdomain.com/webhooks/nomba", // Required by Nomba API
			},
		})

		if err != nil || !resp.Status {
			log.Printf("failed to charge invoice %s: %v", inv.ID, err)
			j.failInvoice(ctx, inv)
			continue
		}

		log.Printf("successfully charged invoice %s", inv.ID)
		j.q.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
			ID:         inv.ID,
			Status:     "paid",
			RetryCount: inv.RetryCount,
		})
	}

	log.Println("Completed automated billing job.")
	return nil
}

func (j *InvoiceJob) failInvoice(ctx context.Context, inv sqlc.Invoice) {
	newRetry := inv.RetryCount + 1
	status := "failed"
	if newRetry >= 3 {
		j.q.SuspendIntegrator(ctx, inv.IntegratorID)
	}

	j.q.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
		ID:         inv.ID,
		Status:     status,
		RetryCount: newRetry,
	})
}
