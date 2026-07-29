package paymentsresponsesdtos

type PaymentNotificationDTO struct {
	PaymentID string `json:"payment_id"`
	Credits   int    `json:"credits"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
}
