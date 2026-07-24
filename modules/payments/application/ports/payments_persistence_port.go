package paymentsports

import (
	"context"

	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
)

type PaymentsPersistencePort interface {
	GetCreditPackage(packageID string) (*paymentsentites.CreditPackageEntity, error)
	ListCreditPackages(ctx context.Context)([]*paymentsentites.CreditPackageEntity,error)
}