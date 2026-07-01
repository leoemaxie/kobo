package tech.triumphsystems.kobo.model;

/** Lifecycle state of a Kobo virtual account. */
public enum IdentityState {
    PENDING("pending"),
    ACTIVE("active"),
    LIMITED("limited"),
    CLOSING("closing"),
    CLOSED("closed"),
    FAILED("failed");

    private final String value;
    IdentityState(String value) { this.value = value; }

    public String getValue() { return value; }

    public static IdentityState of(String s) {
        for (IdentityState v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown IdentityState: " + s);
    }
}
