package tech.triumphsystems.kobo.model;

/** Response from the /healthz endpoint. */
public final class HealthResponse {
    private final String status;
    private final String db;

    public HealthResponse(String status, String db) {
        this.status = status;
        this.db = db;
    }

    public String getStatus() { return status; }
    public String getDb() { return db; }
    public boolean isHealthy() { return "ok".equals(status); }
}
