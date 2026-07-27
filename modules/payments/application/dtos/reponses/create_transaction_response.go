package paymentsresponsesdtos

import paymentsenums "github.com/PurpleSavage/monekai-server/modules/payments/domain/enums"


type CreateTransactionResponseDTO struct {
	TransactionID string
	Status paymentsenums.TransactionStatus
}