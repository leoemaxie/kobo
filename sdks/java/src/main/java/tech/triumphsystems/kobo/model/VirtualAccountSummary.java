package tech.triumphsystems.kobo.model;

/** Bank account details for a provisioned virtual account. */
public final class VirtualAccountSummary {

    private final String accountNumber;
    private final String bankName;
    private final String accountName;

    public VirtualAccountSummary(String accountNumber, String bankName, String accountName) {
        this.accountNumber = accountNumber;
        this.bankName = bankName;
        this.accountName = accountName;
    }

    public String getAccountNumber() { return accountNumber; }
    public String getBankName() { return bankName; }
    public String getAccountName() { return accountName; }

    @Override
    public String toString() {
        return "VirtualAccountSummary{accountNumber='" + accountNumber + "', bankName='" + bankName + "'}";
    }
}
