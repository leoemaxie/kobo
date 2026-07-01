package tech.triumphsystems.kobo.model;

/** Reconciliation status of a transaction. */
public enum TransactionStatus {
    MATCHED("matched"),
    PARTIAL("partial"),
    OVERPAYMENT("overpayment");

    private final String value;
    TransactionStatus(String value) { this.value = value; }

    public String getValue() { return value; }

    public static TransactionStatus of(String s) {
        for (TransactionStatus v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown TransactionStatus: " + s);
    }
}
