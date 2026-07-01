package tech.triumphsystems.kobo;

import tech.triumphsystems.kobo.model.*;
import tech.triumphsystems.kobo.request.*;

import java.net.http.HttpClient;
import java.time.Duration;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Entry point for the Kobo API SDK.
 *
 * <p>Create one instance and reuse it. The underlying {@link HttpClient} is
 * thread-safe. Zero external dependencies — uses only the JDK (Java 11+).
 *
 * <h2>Quick start</h2>
 * <pre>{@code
 * // Production
 * KoboClient kobo = KoboClient.builder()
 *     .apiKey("kobo_live_pk_...")
 *     .apiSecret("kobo_live_sk_...")
 *     .build();
 *
 * // Sandbox
 * KoboClient kobo = KoboClient.sandbox("kobo_test_pk_...", "kobo_test_sk_...");
 *
 * Identity identity = kobo.identities().create(
 *     CreateIdentityRequest.builder()
 *         .externalReference("school:student:4471")
 *         .displayName("John Ade")
 *         .build()
 * );
 * }</pre>
 */
public final class KoboClient {

    private static final String PRODUCTION_BASE_URL = "https://api.kobo.triumphsystems.tech/v1";
    private static final String SANDBOX_BASE_URL    = "https://sandbox.api.kobo.triumphsystems.tech/v1";

    private final KoboHttpClient http;
    private final IdentitiesClient identities;
    private final AccountsClient accounts;
    private final ExceptionsClient exceptions;

    private KoboClient(Builder b) {
        String url = b.baseUrl != null ? b.baseUrl : PRODUCTION_BASE_URL;
        HttpClient hc = b.httpClient != null
                ? b.httpClient
                : HttpClient.newBuilder()
                        .connectTimeout(Duration.ofSeconds(10))
                        .build();
        this.http = new KoboHttpClient(hc, url, b.apiKey, b.apiSecret);
        this.identities = new IdentitiesClient(this.http);
        this.accounts   = new AccountsClient(this.http);
        this.exceptions = new ExceptionsClient(this.http);
    }

    /**
     * Create a production {@link KoboClient} with the given credentials.
     */
    public static KoboClient of(String apiKey, String apiSecret) {
        return builder().apiKey(apiKey).apiSecret(apiSecret).build();
    }

    /**
     * Create a sandbox {@link KoboClient} with the given credentials.
     */
    public static KoboClient sandbox(String apiKey, String apiSecret) {
        return builder().apiKey(apiKey).apiSecret(apiSecret).baseUrl(SANDBOX_BASE_URL).build();
    }

    /** @return Operations on Identity resources. */
    public IdentitiesClient identities() { return identities; }

    /** @return Operations on Account resources (transactions, statements). */
    public AccountsClient accounts() { return accounts; }

    /** @return Operations on Exception resources. */
    public ExceptionsClient exceptions() { return exceptions; }

    /**
     * Call the unauthenticated {@code /healthz} endpoint.
     */
    public HealthResponse health() {
        return http.get("/healthz", null, HealthResponse.class);
    }

    public static Builder builder() { return new Builder(); }

    // ─── Builder ──────────────────────────────────────────────────────────────

    public static final class Builder {
        private String apiKey;
        private String apiSecret;
        private String baseUrl;
        private HttpClient httpClient;

        public Builder apiKey(String v)         { apiKey = v;      return this; }
        public Builder apiSecret(String v)      { apiSecret = v;   return this; }
        public Builder baseUrl(String v)        { baseUrl = v;     return this; }
        public Builder httpClient(HttpClient v) { httpClient = v;  return this; }

        public KoboClient build() {
            Objects.requireNonNull(apiKey,    "apiKey is required");
            Objects.requireNonNull(apiSecret, "apiSecret is required");
            return new KoboClient(this);
        }
    }

    // ─── Sub-clients ──────────────────────────────────────────────────────────

    /**
     * Operations on {@code /identities}.
     */
    public static final class IdentitiesClient {
        private final KoboHttpClient http;
        IdentitiesClient(KoboHttpClient http) { this.http = http; }

        /**
         * Register a new identity and provision its virtual account.
         * Returns immediately with the identity in {@code pending} state.
         */
        public Identity create(CreateIdentityRequest req) {
            return http.post("/identities", req, Identity.class);
        }

        /** Fetch an identity by its Kobo ID. */
        public Identity get(String identityId) {
            return http.get("/identities/" + identityId, null, Identity.class);
        }

        /** Update display profile fields. At least one field must be set. */
        public Identity update(String identityId, UpdateIdentityRequest req) {
            return http.patch("/identities/" + identityId, req, Identity.class);
        }

        /**
         * Initiate closure of an identity's virtual account.
         * Closure is asynchronous; the returned identity will be in {@code closing} state.
         */
        public Identity close(String identityId, CloseIdentityRequest req) {
            return http.post("/identities/" + identityId + "/close", req, Identity.class);
        }

        /**
         * Reopen a {@code closed} identity without re-running provisioning.
         */
        public Identity reopen(String identityId) {
            return http.post("/identities/" + identityId + "/reopen", null, Identity.class);
        }
    }

    /**
     * Operations on {@code /accounts}.
     */
    public static final class AccountsClient {
        private final KoboHttpClient http;
        AccountsClient(KoboHttpClient http) { this.http = http; }

        /**
         * List reconciled transactions for an account with cursor-based pagination.
         *
         * @param cursor Pass {@code null} for the first page.
         * @param limit  Number of records per page (1-200). Pass {@code null} for server default (50).
         */
        public TransactionPage listTransactions(String accountId, String cursor, Integer limit) {
            Map<String, String> q = new LinkedHashMap<>();
            if (cursor != null) q.put("page[cursor]", cursor);
            if (limit  != null) q.put("page[limit]",  limit.toString());
            return http.get("/accounts/" + accountId + "/transactions", q, TransactionPage.class);
        }

        /**
         * Get the structured statement for a given period.
         *
         * @param period {@code YYYY-MM} string. Pass {@code null} to default to the current month.
         */
        public Statement getStatement(String accountId, String period) {
            Map<String, String> q = new LinkedHashMap<>();
            if (period != null) q.put("period", period);
            return http.get("/accounts/" + accountId + "/statement", q, Statement.class);
        }
    }

    /**
     * Operations on {@code /exceptions}.
     */
    public static final class ExceptionsClient {
        private final KoboHttpClient http;
        ExceptionsClient(KoboHttpClient http) { this.http = http; }

        /**
         * List exceptions with optional status filter and cursor-based pagination.
         *
         * @param status Pass {@code null} to default to {@code open}.
         * @param cursor Pass {@code null} for the first page.
         * @param limit  Pass {@code null} for server default (50).
         */
        public ExceptionPage list(ExceptionStatus status, String cursor, Integer limit) {
            Map<String, String> q = new LinkedHashMap<>();
            if (status != null) q.put("status", status.getValue());
            if (cursor != null) q.put("page[cursor]", cursor);
            if (limit  != null) q.put("page[limit]",  limit.toString());
            return http.get("/exceptions", q, ExceptionPage.class);
        }

        /**
         * Apply a resolution to a flagged exception.
         */
        public Exception_ resolve(String exceptionId, ResolveExceptionRequest req) {
            return http.post("/exceptions/" + exceptionId + "/resolve", req, Exception_.class);
        }
    }
}
