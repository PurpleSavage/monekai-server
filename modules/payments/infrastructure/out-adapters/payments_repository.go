package paymentsoutadapters

import (
	"context"
	"errors"
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	paymentsports "github.com/PurpleSavage/monekai-server/modules/payments/application/ports"
	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
	paymentsinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/mappers"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) paymentsports.PaymentsPersistencePort{
	return  &PaymentRepository{
		db: db,
	}
}
func  (r *PaymentRepository)ListCreditPackages(ctx context.Context)([]*paymentsentites.CreditPackageEntity,error){
	var creditPackages []models.CreditPackage
	err:= r.db.WithContext(ctx).Find(&creditPackages).Error
	if err!= nil{
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error getting credit package",
			err,
		)
	}
	packages:=make([]*paymentsentites.CreditPackageEntity,len(creditPackages))
	for i,pkg:= range(creditPackages){
		packages[i]=paymentsinfrastructuremappers.ToCreditPackagesEntity(pkg)
	}
	return packages,nil
}
func (r *PaymentRepository) GetCreditPackage(packageID string) (*paymentsentites.CreditPackageEntity,error) {
	var creditPackage models.CreditPackage
	err:= r.db.Where("id = ?",packageID).Find(&creditPackage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, globalerrors.NewAppError(
				404,
				"Not Found",
				"Credit package not found",
				err,
			)
		}
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error getting credit package",
			err,
		)
	}
	return paymentsinfrastructuremappers.ToCreditPackagesEntity(creditPackage) , nil
}

func (r *PaymentRepository) GetCreditPackageByPriceID(priceID string) (*paymentsentites.CreditPackageEntity, error) {
	var creditPackage models.CreditPackage
	err := r.db.Where("price_id = ?", priceID).First(&creditPackage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, globalerrors.NewAppError(
				404,
				"Not Found",
				"Credit package not found for the given price ID",
				err,
			)
		}
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error getting credit package by price ID",
			err,
		)
	}
	return paymentsinfrastructuremappers.ToCreditPackagesEntity(creditPackage), nil
}

func (r *PaymentRepository) SavePayment(ctx context.Context, payment *paymentsentites.PaymentEntity) error {
	model := paymentsinfrastructuremappers.ToPaymentModel(payment)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return globalerrors.NewAppError(
			500,
			"Database Error",
			"Failed to save payment",
			err,
		)
	}
	payment.ID = model.ID.String()
	return nil
}

func (r *PaymentRepository) FindPaymentByProviderTransactionID(transactionID string) (*paymentsentites.PaymentEntity, error) {
	var payment models.Payment
	err := r.db.Where("provider_transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error finding payment by transaction ID",
			err,
		)
	}
	return paymentsinfrastructuremappers.ToPaymentEntityFromModel(payment), nil
}
