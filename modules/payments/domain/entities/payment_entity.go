package paymentsentities

type PaymentEntity struct {
	ID                    string
	UserID                string
	CreditPackageID       string
	Provider              string
	ProviderTransactionID string
	PriceID               string
	CreditsPurchased      int
	AmountCents           int
	Currency              string
	Status                string
	RawPayload            []byte
}
