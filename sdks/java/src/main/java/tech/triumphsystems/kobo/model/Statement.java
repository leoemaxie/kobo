package tech.triumphsystems.kobo.model;

import java.util.Collections;
import java.util.List;

/** Structured monthly account statement. */
public final class Statement {

    private final String accountId;
    private final String period;
    private final long openingBalanceKobo;
    private final long closingBalanceKobo;
    private final long totalInflowKobo;
    private final List<Transaction> transactions;

    public Statement(
            String accountId,
            String period,
            long openingBalanceKobo,
            long closingBalanceKobo,
            long totalInflowKobo,
            List<Transaction> transactions) {
        this.accountId = accountId;
        this.period = period;
        this.openingBalanceKobo = openingBalanceKobo;
        this.closingBalanceKobo = closingBalanceKobo;
        this.totalInflowKobo = totalInflowKobo;
        this.transactions = transactions != null ? transactions : Collections.emptyList();
    }

    public String getAccountId() { return accountId; }
    /** Statement period in {@code YYYY-MM} format. */
    public String getPeriod() { return period; }
    public long getOpeningBalanceKobo() { return openingBalanceKobo; }
    public long getClosingBalanceKobo() { return closingBalanceKobo; }
    public long getTotalInflowKobo() { return totalInflowKobo; }
    public List<Transaction> getTransactions() { return transactions; }

    @Override
    public String toString() {
        return "Statement{accountId='" + accountId + "', period='" + period + "', txns=" + transactions.size() + "}";
    }
}
