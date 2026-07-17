package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/monnify"
	"github.com/leoemaxie/kobo/internal/payout"
)

type PayoutHandler struct {
	svc         *payout.Service
	monnifyClient *monnify.Client
}

func NewPayoutHandler(svc *payout.Service, monnifyClient *monnify.Client) *PayoutHandler {
	return &PayoutHandler{
		svc:         svc,
		monnifyClient: monnifyClient,
	}
}

func (h *PayoutHandler) Svc() *payout.Service {
	return h.svc
}

func (h *PayoutHandler) ListBanks(w http.ResponseWriter, r *http.Request) {
	banks, err := h.monnifyClient.ListBanks(r.Context())
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusInternalServerError, "internal_error", "Failed to list banks", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banks)
}

func (h *PayoutHandler) LookupBankAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.AccountNumber == "" || req.BankCode == "" {
		apierrors.WriteError(w, http.StatusBadRequest, "missing_fields", "account_number and bank_code are required")
		return
	}

	resp, err := h.monnifyClient.LookupBankAccount(r.Context(), req.AccountNumber, req.BankCode)
	if err != nil {
		apierrors.WriteError(w, http.StatusUnprocessableEntity, "invalid_account", "Bank account verification failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PayoutHandler) SaveBankAccount(w http.ResponseWriter, r *http.Request) {
	consoleCtx := middleware.GetConsoleSessionContext(r.Context())

	var req struct {
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
		BankName      string `json:"bank_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.AccountNumber == "" || req.BankCode == "" || req.BankName == "" {
		apierrors.WriteError(w, http.StatusBadRequest, "missing_fields", "account_number, bank_code, and bank_name are required")
		return
	}

	account, err := h.svc.SaveBankAccount(r.Context(), consoleCtx.IntegratorID, consoleCtx.UserID, req.AccountNumber, req.BankCode, req.BankName)
	if err != nil {
		apierrors.WriteError(w, http.StatusUnprocessableEntity, "invalid_account", "Failed to save bank account: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *PayoutHandler) GetBankAccount(w http.ResponseWriter, r *http.Request) {
	consoleCtx := middleware.GetConsoleSessionContext(r.Context())

	account, err := h.svc.GetBankAccount(r.Context(), consoleCtx.IntegratorID)
	if err != nil {
		if err == payout.ErrNoBankAccount {
			apierrors.WriteError(w, http.StatusNotFound, "no_bank_account", "No active bank account configured")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to get bank account")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *PayoutHandler) RequestPayout(w http.ResponseWriter, r *http.Request) {
	consoleCtx := middleware.GetConsoleSessionContext(r.Context())

	var req struct {
		AmountKobo int64 `json:"amount_kobo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.AmountKobo <= 0 {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_amount", "amount_kobo must be positive")
		return
	}

	payoutRecord, err := h.svc.RequestPayout(r.Context(), consoleCtx.IntegratorID, consoleCtx.UserID, req.AmountKobo)
	if err != nil {
		switch err {
		case payout.ErrNoBankAccount:
			apierrors.WriteError(w, http.StatusUnprocessableEntity, "no_bank_account", "No active bank account configured")
		case payout.ErrBelowMinimum:
			apierrors.WriteError(w, http.StatusUnprocessableEntity, "below_minimum", "Amount is below the minimum threshold")
		case payout.ErrInsufficientBalance:
			apierrors.WriteError(w, http.StatusUnprocessableEntity, "insufficient_balance", "Insufficient balance for payout")
		case payout.ErrPayoutInProgress:
			apierrors.WriteError(w, http.StatusConflict, "payout_in_progress", "A payout is already in progress")
		default:
			apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to process payout request")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payoutRecord)
}

func (h *PayoutHandler) ListPayouts(w http.ResponseWriter, r *http.Request) {
	consoleCtx := middleware.GetConsoleSessionContext(r.Context())

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(20)
	if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil && l > 0 {
		if l > 100 {
			limit = 100
		} else {
			limit = int32(l)
		}
	}

	offset := int32(0)
	if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && o > 0 {
		offset = int32(o)
	}

	payouts, err := h.svc.ListPayouts(r.Context(), consoleCtx.IntegratorID, limit, offset)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to list payouts")
		return
	}

	response := map[string]interface{}{
		"data":   payouts,
		"limit":  limit,
		"offset": offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
