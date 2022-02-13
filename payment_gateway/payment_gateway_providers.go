package payment_gateway_providers

type PaymentGatewayProvider string

const (
	Xendit   PaymentGatewayProvider = "Xendit"
	Midtrans PaymentGatewayProvider = "Midtrans"
)
