package tech.triumphsystems.kobo.request;

import tech.triumphsystems.kobo.JsonSerializable;

import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Request body for {@code PATCH /v1/identities/{id}}.
 * At least one field must be set.
 */
public final class UpdateIdentityRequest implements JsonSerializable {

    private final String displayName;
    private final Map<String, Object> metadata;

    private UpdateIdentityRequest(Builder b) {
        this.displayName = b.displayName;
        this.metadata = b.metadata;
    }

    public static Builder builder() { return new Builder(); }

    @Override
    public Map<String, Object> toJsonMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        if (displayName != null) m.put("display_name", displayName);
        if (metadata != null) m.put("metadata", metadata);
        return m;
    }

    public static final class Builder {
        private String displayName;
        private Map<String, Object> metadata;

        public Builder displayName(String v) { displayName = v; return this; }
        public Builder metadata(Map<String, Object> v) { metadata = v; return this; }
        public UpdateIdentityRequest build() { return new UpdateIdentityRequest(this); }
    }
}
