// =============================================================================
// Kobo SDK – Type Definitions
// Zero external dependencies. Native fetch API.
// =============================================================================

// ─── Enums ───────────────────────────────────────────────────────────────────

export type IdentityState =
  | "pending"
  | "active"
  | "limited"
  | "closing"
  | "closed"
  | "failed";



export type TransactionDirection = "inbound";

export type TransactionStatus = "matched" | "partial" | "overpayment";

export type ExceptionType =
  | "payment_to_closed_account"
  | "payment_to_unknown_account"
  | "payment_during_closing";

export type ExceptionStatus = "open" | "resolved";

export type ExceptionResolutionAction =
  | "return_to_sender"
  | "redirect_to_successor"
  | "manual_override";

export type SweepDestinationType =
  | "refund_to_source"
  | "integrator_account"
  | "successor_identity";

// ─── Core Models ─────────────────────────────────────────────────────────────

export interface VirtualAccountSummary {
  account_number: string;
  bank_name: string;
  account_name: string;
  expected_amount_kobo?: number | null;
  is_expired: boolean;
}

export interface Identity {
  id: string;
  external_reference: string;
  display_name: string;
  state: IdentityState;
  virtual_account?: VirtualAccountSummary | null;
  metadata?: Record<string, unknown>;
  failure_reason?: string | null;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: string;
  account_id: string;
  /** Amount in kobo (1/100 of a Naira). Always positive; use `direction` for flow. */
  amount_kobo: number;
  direction: TransactionDirection;
  status: TransactionStatus;
  nomba_reference: string;
  occurred_at: string;
}

export interface Statement {
  account_id: string;
  period: string;
  opening_balance_kobo: number;
  closing_balance_kobo: number;
  total_inflow_kobo: number;
  transactions: Transaction[];
}

export interface ExceptionResolution {
  action: ExceptionResolutionAction;
  notes?: string;
}

export interface Exception {
  id: string;
  type: ExceptionType;
  amount_kobo: number;
  nomba_reference: string;
  related_account_id?: string | null;
  status: ExceptionStatus;
  resolution?: ExceptionResolution | null;
  detected_at: string;
  resolved_at?: string | null;
}

// ─── Request Bodies ───────────────────────────────────────────────────────────

export interface CreateIdentityRequest {
  external_reference: string;
  display_name: string;
  metadata?: Record<string, unknown>;
}

export interface UpdateIdentityRequest {
  display_name?: string;
  metadata?: Record<string, unknown>;
}

export interface SweepDestination {
  type: SweepDestinationType;
  successor_identity_id?: string;
  integrator_account_reference?: string;
}

export interface CloseIdentityRequest {
  sweep_destination: SweepDestination;
  reason?: string;
}

export interface ResolveExceptionRequest {
  action: ExceptionResolutionAction;
  successor_identity_id?: string;
  notes?: string;
}

// ─── Paginated Responses ──────────────────────────────────────────────────────

export interface PagedResponse<T> {
  data: T[];
  next_cursor: string | null;
}

export type TransactionListResponse = PagedResponse<Transaction>;
export type ExceptionListResponse = PagedResponse<Exception>;

// ─── Query Params ─────────────────────────────────────────────────────────────

export interface PaginationParams {
  cursor?: string;
  limit?: number;
}

export interface ListExceptionsParams extends PaginationParams {
  status?: ExceptionStatus;
}

export interface GetStatementOptions {
  period?: string;
}

export interface ListTransactionsOptions extends PaginationParams {}
export interface ListIdentitiesOptions {
  state?: IdentityState;
  limit?: number;
  offset?: number;
}

// ─── Health ───────────────────────────────────────────────────────────────────

export interface HealthResponse {
  status: "ok";
  db: "ok" | "degraded";
}

// ─── Errors ───────────────────────────────────────────────────────────────────

export interface KoboErrorBody {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
