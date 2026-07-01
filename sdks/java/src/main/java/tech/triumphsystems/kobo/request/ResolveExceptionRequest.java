package tech.triumphsystems.kobo.request;

import tech.triumphsystems.kobo.JsonSerializable;
import tech.triumphsystems.kobo.model.ResolutionAction;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Request body for {@code POST /v1/exceptions/{id}/resolve}.
 */
public final class ResolveExceptionRequest implements JsonSerializable {

    private final ResolutionAction action;
    private final String successorIdentityId;
    private final String notes;

    private ResolveExceptionRequest(Builder b) {
        this.action = Objects.requireNonNull(b.action, "action is required");
        this.successorIdentityId = b.successorIdentityId;
        this.notes = b.notes;
    }

    public static Builder builder() { return new Builder(); }

    @Override
    public Map<String, Object> toJsonMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("action", action.getValue());
        if (successorIdentityId != null) m.put("successor_identity_id", successorIdentityId);
        if (notes != null) m.put("notes", notes);
        return m;
    }

    public static final class Builder {
        private ResolutionAction action;
        private String successorIdentityId;
        private String notes;

        public Builder action(ResolutionAction v) { action = v; return this; }
        public Builder successorIdentityId(String v) { successorIdentityId = v; return this; }
        public Builder notes(String v) { notes = v; return this; }
        public ResolveExceptionRequest build() { return new ResolveExceptionRequest(this); }
    }
}
