package sqlc

const (
	ConsoleEnvironmentSandbox    = "sandbox"
	ConsoleEnvironmentProduction = "production"
)

type ConsoleEnvironment = string
type Invoice = ConsoleInvoice
type PaymentMethod = ConsolePaymentMethod
