package tech.triumphsystems.kobo.model;

import java.time.Instant;
import java.util.Map;

/**
 * The root Kobo resource. Represents an individual or entity with a
 * provisioned virtual bank account.
 */
public final class Identity {

    private final String id;
    private final String externalReference;
    private final String displayName;
    private final IdentityState state;
    private final VirtualAccountSummary virtualAccount;
    private final Map<String, Object> metadata;
    private final String failureReason;
    private final Instant createdAt;
    private final Instant updatedAt;

    public Identity(
            String id,
            String externalReference,
            String displayName,
            IdentityState state,
            VirtualAccountSummary virtualAccount,
            Map<String, Object> metadata,
            String failureReason,
            Instant createdAt,
            Instant updatedAt) {
        this.id = id;
        this.externalReference = externalReference;
        this.displayName = displayName;
        this.state = state;
        this.virtualAccount = virtualAccount;
        this.metadata = metadata;
        this.failureReason = failureReason;
        this.createdAt = createdAt;
        this.updatedAt = updatedAt;
    }

    /** UUIDv4 identity ID. */
    public String getId() { return id; }

    /** Integrator-supplied stable reference (e.g. their own customer ID). */
    public String getExternalReference() { return externalReference; }

    public String getDisplayName() { return displayName; }

    public IdentityState getState() { return state; }

    /** Null while in {@code pending} or {@code failed} state. */
    public VirtualAccountSummary getVirtualAccount() { return virtualAccount; }

    public Map<String, Object> getMetadata() { return metadata; }

    /** Non-null only when state is {@code failed}. */
    public String getFailureReason() { return failureReason; }

    public Instant getCreatedAt() { return createdAt; }

    public Instant getUpdatedAt() { return updatedAt; }

    @Override
    public String toString() {
        return "Identity{id='" + id + "', state=" + state + ", displayName='" + displayName + "'}";
    }
}
