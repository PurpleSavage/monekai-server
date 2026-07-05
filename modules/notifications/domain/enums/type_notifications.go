package notificationsenums
type TypeNotification string

const (
	ReplicateError   TypeNotification = "replicate_error"
	ReplicateSuccess TypeNotification = "replicate_success"
	Payment          TypeNotification = "payment"
	Info            TypeNotification = "info"
)