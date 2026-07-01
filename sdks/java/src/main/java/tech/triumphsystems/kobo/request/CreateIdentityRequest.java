package tech.triumphsystems.kobo.request;

import tech.triumphsystems.kobo.JsonSerializable;
import tech.triumphsystems.kobo.model.KycTier;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Request body for {@code POST /v1/identities}.
 *
 * <p>Use the builder:
 * <pre>{@code
 * var req = CreateIdentityRequest.builder()
 *     .externalReference("school:student:4471")
 *     .displayName("John Ade")
 *     .kycTierHint(KycTier.TIER_1)
 *     .metadata(Map.of("class", "JSS 2"))
 *     .build();
 * }</pre>
 */
public final class CreateIdentityRequest implements JsonSerializable {

    private final String externalReference;
    private final String displayName;
    private final KycTier kycTierHint;
    private final Map<String, Object> metadata;

    private CreateIdentityRequest(Builder b) {
        this.externalReference = Objects.requireNonNull(b.externalReference, "externalReference is required");
        this.displayName = Objects.requireNonNull(b.displayName, "displayName is required");
        this.kycTierHint = b.kycTierHint;
        this.metadata = b.metadata;
    }

    public static Builder builder() { return new Builder(); }

    @Override
    public Map<String, Object> toJsonMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("external_reference", externalReference);
        m.put("display_name", displayName);
        if (kycTierHint != null) m.put("kyc_tier_hint", kycTierHint.getValue());
        if (metadata != null && !metadata.isEmpty()) m.put("metadata", metadata);
        return m;
    }

    public static final class Builder {
        private String externalReference;
        private String displayName;
        private KycTier kycTierHint;
        private Map<String, Object> metadata;

        public Builder externalReference(String v) { externalReference = v; return this; }
        public Builder displayName(String v) { displayName = v; return this; }
        public Builder kycTierHint(KycTier v) { kycTierHint = v; return this; }
        public Builder metadata(Map<String, Object> v) { metadata = v; return this; }
        public CreateIdentityRequest build() { return new CreateIdentityRequest(this); }
    }
}
