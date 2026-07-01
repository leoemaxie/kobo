package tech.triumphsystems.kobo.model;

/** Whether an exception has been resolved. */
public enum ExceptionStatus {
    OPEN("open"),
    RESOLVED("resolved");

    private final String value;
    ExceptionStatus(String value) { this.value = value; }

    public String getValue() { return value; }

    public static ExceptionStatus of(String s) {
        for (ExceptionStatus v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown ExceptionStatus: " + s);
    }
}
