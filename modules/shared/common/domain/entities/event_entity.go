package commonentities

type EventName string

const (
	EventSampleReady EventName = "sample_ready"
	EventSampleError EventName = "sample_error"

	EventPaymentSuccess EventName = "payment_success"
	EventPaymentFailed  EventName = "payment_failed"
)

type Event struct {
	Name EventName
	Data any
}
