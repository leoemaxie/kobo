import { KoboError } from "./errors.js";
import type {
  CloseIdentityRequest,
  CreateIdentityRequest,
  Exception,
  ExceptionListResponse,
  GetStatementOptions,
  HealthResponse,
  Identity,
  ListExceptionsParams,
  ListIdentitiesOptions,
  ListTransactionsOptions,
  PaginationParams,
  ResolveExceptionRequest,
  Statement,
  TransactionListResponse,
  UpdateIdentityRequest,
} from "./types.js";

// ─── Constants ────────────────────────────────────────────────────────────────

const PRODUCTION_BASE_URL = "https://api.kobo.triumphsystems.tech/v1";
const SANDBOX_BASE_URL = "https://sandbox.api.kobo.triumphsystems.tech/v1";
const SDK_VERSION = "0.1.0";

// ─── Client Configuration ─────────────────────────────────────────────────────

export interface KoboClientOptions {
  /**
   * Override the base URL. Useful for pointing at a local server during testing.
   * @default "https://api.kobo.triumphsystems.tech/v1"
   */
  baseUrl?: string;
  /**
   * AbortSignal timeout in milliseconds applied to every request.
   * @default 30_000
   */
  timeoutMs?: number;
}

// ─── Sub-clients ──────────────────────────────────────────────────────────────

/**
 * IdentitiesClient groups all /identities operations.
 */
export class IdentitiesClient {
  constructor(private readonly http: KoboHttpLayer) {}

  /**
   * Register a new identity and provision its virtual account.
   *
   * The identity is returned immediately in `pending` state. Poll `get()` or
   * listen for the `identity.activated` webhook to confirm provisioning.
   */
  async create(body: CreateIdentityRequest): Promise<Identity> {
    return this.http.post<Identity>("/identities", body);
  }

  /**
   * Fetch an identity record, its current state, and linked account.
   */
  async get(identityId: string): Promise<Identity> {
    return this.http.get<Identity>(`/identities/${identityId}`);
  }

  /**
   * List identities for the integrator.
   */
  async list(opts?: ListIdentitiesOptions): Promise<Identity[]> {
    const q: Record<string, string> = {};
    if (opts?.state) q["state"] = opts.state;
    if (opts?.limit !== undefined) q["limit"] = String(opts.limit);
    if (opts?.offset !== undefined) q["offset"] = String(opts.offset);
    return this.http.get<Identity[]>("/identities", q);
  }

  /**
   * Update display profile fields (`display_name`, `metadata`).
   * At least one field must be provided.
   */
  async update(
    identityId: string,
    body: UpdateIdentityRequest
  ): Promise<Identity> {
    return this.http.patch<Identity>(`/identities/${identityId}`, body);
  }

  /**
   * Initiate closure of an identity's virtual account.
   *
   * Closure is asynchronous — the response (202) is returned immediately.
   * Listen for the `identity.closed` webhook or poll `get()`.
   */
  async close(
    identityId: string,
    body: CloseIdentityRequest
  ): Promise<Identity> {
    return this.http.post<Identity>(`/identities/${identityId}/close`, body);
  }

  /**
   * Reopen a CLOSED identity's account without re-running full provisioning.
   * The identity must be in `closed` state.
   */
  async reopen(identityId: string): Promise<Identity> {
    return this.http.post<Identity>(`/identities/${identityId}/reopen`);
  }
}

/**
 * AccountsClient groups all /accounts operations.
 */
export class AccountsClient {
  constructor(private readonly http: KoboHttpLayer) {}

  /**
   * Retrieve a paginated, reconciled transaction history for an account.
   */
  async listTransactions(
    accountId: string,
    params?: ListTransactionsOptions
  ): Promise<TransactionListResponse> {
    const q = buildPaginationQuery(params);
    return this.http.get<TransactionListResponse>(
      `/accounts/${accountId}/transactions`,
      q
    );
  }

  /**
   * Retrieve a structured statement for a given period.
   *
   * @param period - A `YYYY-MM` string. Defaults to the current month when omitted.
   */
  async getStatement(
    accountId: string,
    opts?: GetStatementOptions
  ): Promise<Statement> {
    const q: Record<string, string> = {};
    if (opts?.period) q["period"] = opts.period;
    return this.http.get<Statement>(`/accounts/${accountId}/statement`, q);
  }
}

/**
 * ExceptionsClient groups all /exceptions operations.
 */
export class ExceptionsClient {
  constructor(private readonly http: KoboHttpLayer) {}

  /**
   * List unresolved (or resolved) misdirected-payment or unmatched-transfer cases.
   */
  async list(params?: ListExceptionsParams): Promise<ExceptionListResponse> {
    const q = buildPaginationQuery(params);
    if (params?.status) q["status"] = params.status;
    return this.http.get<ExceptionListResponse>("/exceptions", q);
  }

  /**
   * Apply a resolution to a flagged exception.
   */
  async resolve(
    exceptionId: string,
    body: ResolveExceptionRequest
  ): Promise<Exception> {
    return this.http.post<Exception>(
      `/exceptions/${exceptionId}/resolve`,
      body
    );
  }
}

// ─── Main Client ─────────────────────────────────────────────────────────────

/**
 * KoboClient is the main entry point for the Kobo SDK.
 *
 * ```ts
 * import { KoboClient } from "@kobo/sdk";
 *
 * const kobo = new KoboClient("kobo_live_pk_...", "kobo_live_sk_...");
 * const identity = await kobo.identities.create({ ... });
 * ```
 *
 * Zero external dependencies — uses the native `fetch` API.
 */
export class KoboClient {
  /** Operations on Identity resources */
  readonly identities: IdentitiesClient;
  /** Operations on Account resources (transactions, statements) */
  readonly accounts: AccountsClient;
  /** Operations on Exception resources */
  readonly exceptions: ExceptionsClient;

  private readonly http: KoboHttpLayer;

  constructor(
    apiKey: string,
    apiSecret: string,
    opts: KoboClientOptions = {}
  ) {
    const baseUrl = opts.baseUrl ?? PRODUCTION_BASE_URL;
    const timeoutMs = opts.timeoutMs ?? 30_000;

    this.http = new KoboHttpLayer(apiKey, apiSecret, baseUrl, timeoutMs);
    this.identities = new IdentitiesClient(this.http);
    this.accounts = new AccountsClient(this.http);
    this.exceptions = new ExceptionsClient(this.http);
  }

  /**
   * Create a KoboClient pointing at the sandbox environment.
   *
   * ```ts
   * const kobo = KoboClient.sandbox("kobo_test_pk_...", "kobo_test_sk_...");
   * ```
   */
  static sandbox(
    apiKey: string,
    apiSecret: string,
    opts: Omit<KoboClientOptions, "baseUrl"> = {}
  ): KoboClient {
    return new KoboClient(apiKey, apiSecret, {
      ...opts,
      baseUrl: SANDBOX_BASE_URL,
    });
  }

  /**
   * Call the unauthenticated health endpoint.
   * Useful as a connectivity check before sending real traffic.
   */
  async health(): Promise<HealthResponse> {
    return this.http.get<HealthResponse>("/healthz");
  }
}

// ─── Internal HTTP Layer ──────────────────────────────────────────────────────

/** @internal */
class KoboHttpLayer {
  private readonly authHeader: string;

  constructor(
    apiKey: string,
    apiSecret: string,
    private readonly baseUrl: string,
    private readonly timeoutMs: number
  ) {
    // HTTP Basic Auth: base64(apiKey:apiSecret)
    this.authHeader =
      "Basic " + btoa(`${apiKey}:${apiSecret}`);
  }

  async get<T>(path: string, query?: Record<string, string>): Promise<T> {
    return this.request<T>("GET", path, undefined, query);
  }

  async post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>("POST", path, body);
  }

  async patch<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>("PATCH", path, body);
  }

  private async request<T>(
    method: string,
    path: string,
    body?: unknown,
    query?: Record<string, string>
  ): Promise<T> {
    let url = this.baseUrl + path;
    if (query && Object.keys(query).length > 0) {
      url += "?" + new URLSearchParams(query).toString();
    }

    const headers: Record<string, string> = {
      Authorization: this.authHeader,
      Accept: "application/json",
      "User-Agent": `kobo-ts/${SDK_VERSION}`,
    };
    if (body !== undefined) {
      headers["Content-Type"] = "application/json";
    }

    const signal = AbortSignal.timeout(this.timeoutMs);

    const init: RequestInit = {
      method,
      headers,
      signal,
    };
    if (body !== undefined) {
      init.body = JSON.stringify(body);
    }

    const response = await fetch(url, init);

    // Parse body regardless of status so errors can be decoded
    const text = await response.text();
    let parsed: unknown;
    try {
      parsed = text.length > 0 ? JSON.parse(text) : undefined;
    } catch {
      parsed = undefined;
    }

    if (!response.ok) {
      if (
        parsed &&
        typeof parsed === "object" &&
        "code" in parsed &&
        "message" in parsed
      ) {
        throw new KoboError(response.status, parsed as {
          code: string;
          message: string;
          details?: Record<string, unknown>;
        });
      }
      throw new KoboError(response.status, {
        code: "http_error",
        message: `HTTP ${response.status}: ${text}`,
      });
    }

    return parsed as T;
  }
}

// ─── Pagination Helpers ───────────────────────────────────────────────────────

function buildPaginationQuery(
  params?: PaginationParams
): Record<string, string> {
  const q: Record<string, string> = {};
  if (params?.cursor) q["page[cursor]"] = params.cursor;
  if (params?.limit !== undefined) q["page[limit]"] = String(params.limit);
  return q;
}

export type {
  GetStatementOptions,
  ListExceptionsParams,
  ListIdentitiesOptions,
  ListTransactionsOptions,
};
