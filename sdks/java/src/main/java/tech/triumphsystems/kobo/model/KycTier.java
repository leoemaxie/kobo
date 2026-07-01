package tech.triumphsystems.kobo.model;

/** KYC verification tier for an identity. */
public enum KycTier {
    TIER_1("tier_1"),
    TIER_2("tier_2"),
    TIER_3("tier_3");

    private final String value;
    KycTier(String value) { this.value = value; }

    public String getValue() { return value; }

    public static KycTier of(String s) {
        for (KycTier v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown KycTier: " + s);
    }
}
