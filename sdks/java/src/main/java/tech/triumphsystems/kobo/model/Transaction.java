package tech.triumphsystems.kobo.model;

import java.time.Instant;

/**
 * A single reconciled ledger transaction on a virtual account.
 *
 * <p>{@code amountKobo} is always positive. v1 only records inbound transactions.
 */
public final class Transaction {

    private final String id;
    private final String accountId;
    /** Amount in kobo (1/100 Naira). Always positive. */
    private final long amountKobo;
    private final String direction;
    private final TransactionStatus status;
    private final String monnifyReference;
    private final Instant occurredAt;

    public Transaction(
            String id,
            String accountId,
            long amountKobo,
            String direction,
            TransactionStatus status,
            String monnifyReference,
            Instant occurredAt) {
        this.id = id;
        this.accountId = accountId;
        this.amountKobo = amountKobo;
        this.direction = direction;
        this.status = status;
        this.monnifyReference = monnifyReference;
        this.occurredAt = occurredAt;
    }

    public String getId() { return id; }
    public String getAccountId() { return accountId; }
    /** Amount in kobo (1/100 Naira). Always positive. */
    public long getAmountKobo() { return amountKobo; }
    /** Returns the amount in full Naira as a double. */
    public double getAmountNaira() { return amountKobo / 100.0; }
    public String getDirection() { return direction; }
    public TransactionStatus getStatus() { return status; }
    public String getMonnifyReference() { return monnifyReference; }
    public Instant getOccurredAt() { return occurredAt; }

    @Override
    public String toString() {
        return "Transaction{id='" + id + "', amountKobo=" + amountKobo + ", status=" + status + "}";
    }
}
