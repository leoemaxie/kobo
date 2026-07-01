package tech.triumphsystems.kobo.model;

/** Where account funds should be swept on identity closure. */
public enum SweepDestinationType {
    REFUND_TO_SOURCE("refund_to_source"),
    INTEGRATOR_ACCOUNT("integrator_account"),
    SUCCESSOR_IDENTITY("successor_identity");

    private final String value;
    SweepDestinationType(String value) { this.value = value; }

    public String getValue() { return value; }

    public static SweepDestinationType of(String s) {
        for (SweepDestinationType v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown SweepDestinationType: " + s);
    }
}
