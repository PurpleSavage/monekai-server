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

type SamplerRepository struct {
	db *gorm.DB
}

func NewSamplerrepository(db *gorm.DB) paymentsports.PaymentsPersistencePort{
	return  &SamplerRepository{
		db: db,
	}
}
func  (r *SamplerRepository)ListCreditPackages(ctx context.Context)([]*paymentsentites.CreditPackageEntity,error){
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
func (r *SamplerRepository) GetCreditPackage(packageID string) (*paymentsentites.CreditPackageEntity,error) {
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
