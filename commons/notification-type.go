package commons

const (
	PaymentNotification = "Payment Notification"
	NewConsultantClient = "New Consultant Client"
)

const (
	EmailNotificationType = "Email"
	BellNotificationType  = "Bell"
	AllNotificationType   = "All"
)

type NotificationServices struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var JoranvestNotificationServices = []NotificationServices{
	{
		Name: PaymentNotification,
		Type: EmailNotificationType,
	},
	{
		Name: NewConsultantClient,
		Type: EmailNotificationType,
	},
}
