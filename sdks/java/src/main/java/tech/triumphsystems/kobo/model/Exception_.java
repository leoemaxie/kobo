package tech.triumphsystems.kobo.model;

import java.time.Instant;

/**
 * A misdirected-payment or unmatched-transfer exception.
 *
 * <p>Named {@code Exception_} with a trailing underscore to avoid a collision
 * with {@link java.lang.Exception}.
 */
public final class Exception_ {

    private final String id;
    private final ExceptionType type;
    private final long amountKobo;
    private final String monnifyReference;
    private final String relatedAccountId;
    private final ExceptionStatus status;
    private final ExceptionResolution resolution;
    private final Instant detectedAt;
    private final Instant resolvedAt;

    public Exception_(
            String id,
            ExceptionType type,
            long amountKobo,
            String monnifyReference,
            String relatedAccountId,
            ExceptionStatus status,
            ExceptionResolution resolution,
            Instant detectedAt,
            Instant resolvedAt) {
        this.id = id;
        this.type = type;
        this.amountKobo = amountKobo;
        this.monnifyReference = monnifyReference;
        this.relatedAccountId = relatedAccountId;
        this.status = status;
        this.resolution = resolution;
        this.detectedAt = detectedAt;
        this.resolvedAt = resolvedAt;
    }

    public String getId() { return id; }
    public ExceptionType getType() { return type; }
    public long getAmountKobo() { return amountKobo; }
    public double getAmountNaira() { return amountKobo / 100.0; }
    public String getMonnifyReference() { return monnifyReference; }
    public String getRelatedAccountId() { return relatedAccountId; }
    public ExceptionStatus getStatus() { return status; }
    public ExceptionResolution getResolution() { return resolution; }
    public Instant getDetectedAt() { return detectedAt; }
    /** Null if not yet resolved. */
    public Instant getResolvedAt() { return resolvedAt; }

    @Override
    public String toString() {
        return "Exception_{id='" + id + "', type=" + type + ", status=" + status + "}";
    }
}
