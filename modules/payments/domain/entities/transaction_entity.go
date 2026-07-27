package paymentsentities

import paymentsenums "github.com/PurpleSavage/monekai-server/modules/payments/domain/enums"

type Transaction struct {
	ID        string
	Status    paymentsenums.TransactionStatus
	PriceID   string
	Currency  string
	Completed bool
}