package paymentsinfrastructuremappers

import (
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	paymentsentities "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func ToPaymentEntityFromModel(model models.Payment) *paymentsentities.PaymentEntity {
	return &paymentsentities.PaymentEntity{
		ID:                    model.ID.String(),
		UserID:                model.UserID.String(),
		CreditPackageID:       model.CreditPackageID.String(),
		Provider:              model.Provider,
		ProviderTransactionID: model.ProviderTransactionID,
		PriceID:               model.PriceID,
		CreditsPurchased:      model.CreditsPurchased,
		AmountCents:           model.AmountCents,
		Currency:              model.Currency,
		Status:                model.Status,
	}
}

func ToPaymentModel(entity *paymentsentities.PaymentEntity) *models.Payment {
	userID, _ := uuid.Parse(entity.UserID)
	packageID, _ := uuid.Parse(entity.CreditPackageID)

	return &models.Payment{
		UserID:                userID,
		CreditPackageID:       packageID,
		Provider:              entity.Provider,
		ProviderTransactionID: entity.ProviderTransactionID,
		PriceID:               entity.PriceID,
		CreditsPurchased:      entity.CreditsPurchased,
		AmountCents:           entity.AmountCents,
		Currency:              entity.Currency,
		Status:                entity.Status,
		RawPayload:            datatypes.JSON(entity.RawPayload),
	}
}
