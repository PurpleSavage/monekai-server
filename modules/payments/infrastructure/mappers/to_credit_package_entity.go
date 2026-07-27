package paymentsinfrastructuremappers

import (
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	paymentsentities "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
)

func ToCreditPackagesEntity(model models.CreditPackage) *paymentsentities.CreditPackageEntity{
	 return &paymentsentities.CreditPackageEntity{
			Id: model.ID,
			Provider:model.Provider,
			PriceId:model.PriceID,
			Name:model.Name,
			Credits:model.Credits,
			PriceCents:model.PriceCents,
			Currency:model.Currency,
			Active:model.Active,
	}
}