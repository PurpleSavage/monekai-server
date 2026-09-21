package accountinfrastructuremappers

import (
	accountresponsesdtos "github.com/PurpleSavage/monekai-server/modules/account/application/dtos/responses"
	accountentities "github.com/PurpleSavage/monekai-server/modules/account/domain/entities"
)

func ToPaymentHistoryItemDTO(entity accountentities.PaymentHistoryEntity) accountresponsesdtos.PaymentHistoryItemDTO {
	return accountresponsesdtos.PaymentHistoryItemDTO{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		Credits:     entity.Credits,
		AmountCents: entity.AmountCents,
		Currency:    entity.Currency,
		Status:      entity.Status,
	}
}

func BuildListPaymentHistoryResponseDTO(entities []accountentities.PaymentHistoryEntity) []accountresponsesdtos.PaymentHistoryItemDTO {
	items := make([]accountresponsesdtos.PaymentHistoryItemDTO, len(entities))
	for i, entity := range entities {
		items[i] = ToPaymentHistoryItemDTO(entity)
	}
	return items
}