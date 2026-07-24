package paymentsinfrastructuremappers

import (
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
)

func ToCreditPackagesEntity(model models.CreditPackage) *paymentsentites.CreditPackageEntity{
	 return &paymentsentites.CreditPackageEntity{
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