package paymentsentites

import (
	paymentsvalueobjects "github.com/PurpleSavage/monekai-server/modules/payments/domain/valueobjects"
	"github.com/google/uuid"
)

type CreditPackageEntity struct{
	Id          uuid.UUID
	Provider 	paymentsvalueobjects.Provider
	PriceId		string
	Name 		string
	Credits 	int
	PriceCents 	int
	Currency	string
	Active 		bool
}