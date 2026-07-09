package tech.triumphsystems.kobo.model;

/** Bank account details for a provisioned virtual account. */
public final class VirtualAccountSummary {

    private final String accountNumber;
    private final String bankName;
    private final String accountName;
    private final Long expectedAmountKobo;
    private final boolean isExpired;

    public VirtualAccountSummary(String accountNumber, String bankName, String accountName, Long expectedAmountKobo, boolean isExpired) {
        this.accountNumber = accountNumber;
        this.bankName = bankName;
        this.accountName = accountName;
        this.expectedAmountKobo = expectedAmountKobo;
        this.isExpired = isExpired;
    }

    public String getAccountNumber() { return accountNumber; }
    public String getBankName() { return bankName; }
    public String getAccountName() { return accountName; }
    public Long getExpectedAmountKobo() { return expectedAmountKobo; }
    public boolean isExpired() { return isExpired; }

    @Override
    public String toString() {
        return "VirtualAccountSummary{accountNumber='" + accountNumber + "', bankName='" + bankName + "'}";
    }
}
