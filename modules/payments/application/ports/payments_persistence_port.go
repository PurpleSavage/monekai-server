package paymentsports

import (
	"context"

	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
)

type PaymentsPersistencePort interface {
	GetCreditPackage(packageID string) (*paymentsentites.CreditPackageEntity, error)
	GetCreditPackageByPriceID(priceID string) (*paymentsentites.CreditPackageEntity, error)
	ListCreditPackages(ctx context.Context)([]*paymentsentites.CreditPackageEntity,error)
	SavePayment(ctx context.Context, payment *paymentsentites.PaymentEntity) error
	FindPaymentByProviderTransactionID(transactionID string) (*paymentsentites.PaymentEntity, error)
}