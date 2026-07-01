package tech.triumphsystems.kobo.model;

/** Category of a flagged exception. */
public enum ExceptionType {
    PAYMENT_TO_CLOSED_ACCOUNT("payment_to_closed_account"),
    PAYMENT_TO_UNKNOWN_ACCOUNT("payment_to_unknown_account"),
    PAYMENT_DURING_CLOSING("payment_during_closing");

    private final String value;
    ExceptionType(String value) { this.value = value; }

    public String getValue() { return value; }

    public static ExceptionType of(String s) {
        for (ExceptionType v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown ExceptionType: " + s);
    }
}
