package tech.triumphsystems.kobo.request;

import tech.triumphsystems.kobo.JsonSerializable;
import tech.triumphsystems.kobo.model.SweepDestinationType;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Request body for {@code POST /v1/identities/{id}/close}.
 *
 * <pre>{@code
 * var req = CloseIdentityRequest.builder()
 *     .sweepDestinationType(SweepDestinationType.REFUND_TO_SOURCE)
 *     .reason("Customer request")
 *     .build();
 * }</pre>
 */
public final class CloseIdentityRequest implements JsonSerializable {

    private final SweepDestinationType sweepDestinationType;
    private final String successorIdentityId;
    private final String integratorAccountReference;
    private final String reason;

    private CloseIdentityRequest(Builder b) {
        this.sweepDestinationType = Objects.requireNonNull(b.sweepDestinationType, "sweepDestinationType is required");
        this.successorIdentityId = b.successorIdentityId;
        this.integratorAccountReference = b.integratorAccountReference;
        this.reason = b.reason;
    }

    public static Builder builder() { return new Builder(); }

    @Override
    public Map<String, Object> toJsonMap() {
        Map<String, Object> sweep = new LinkedHashMap<>();
        sweep.put("type", sweepDestinationType.getValue());
        if (successorIdentityId != null) sweep.put("successor_identity_id", successorIdentityId);
        if (integratorAccountReference != null) sweep.put("integrator_account_reference", integratorAccountReference);

        Map<String, Object> m = new LinkedHashMap<>();
        m.put("sweep_destination", sweep);
        if (reason != null) m.put("reason", reason);
        return m;
    }

    public static final class Builder {
        private SweepDestinationType sweepDestinationType;
        private String successorIdentityId;
        private String integratorAccountReference;
        private String reason;

        public Builder sweepDestinationType(SweepDestinationType v) { sweepDestinationType = v; return this; }
        public Builder successorIdentityId(String v) { successorIdentityId = v; return this; }
        public Builder integratorAccountReference(String v) { integratorAccountReference = v; return this; }
        public Builder reason(String v) { reason = v; return this; }
        public CloseIdentityRequest build() { return new CloseIdentityRequest(this); }
    }
}
