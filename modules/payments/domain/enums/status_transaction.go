package paymentsenums
type TransactionStatus string

const (
	TransactionStatusDraft     TransactionStatus = "draft"
	TransactionStatusReady     TransactionStatus = "ready"
	TransactionStatusBilled    TransactionStatus = "billed"
	TransactionStatusPaid      TransactionStatus = "paid"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCanceled  TransactionStatus = "canceled"
	TransactionStatusPastDue   TransactionStatus = "past_due"
)