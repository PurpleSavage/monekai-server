package accountresponsesdtos

type PaymentHistoryItemDTO struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Credits     int    `json:"credits"`
	AmountCents int    `json:"amountCents"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
}