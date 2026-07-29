package paymentsrequestsdtos

import paymentsresponsesdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/reponses"

type DataPaymentNotify struct {
	UserID    string
	PaymentID string
	Credits   int
	Amount    int
	Currency  string
	Status    string
	Data      *paymentsresponsesdtos.PaymentNotificationDTO
}
