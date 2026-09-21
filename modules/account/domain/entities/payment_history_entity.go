package accountentities

// PaymentHistoryEntity representa una compra de créditos del usuario
// en el módulo de cuenta (historial de pagos).
type PaymentHistoryEntity struct {
	ID          string
	CreatedAt   string
	Credits     int
	AmountCents int
	Currency    string
	Status      string
}