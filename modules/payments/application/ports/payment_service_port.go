package paymentsports

import (
	"context"

	"github.com/PaddleHQ/paddle-go-sdk/v5"
	paymentsresponsesdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/reponses"
	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
	paymentsvalueobjects "github.com/PurpleSavage/monekai-server/modules/payments/domain/valueobjects"
)

type PaymentServicePort interface {
	CreateTransaction(
		ctx context.Context,
		vo *paymentsvalueobjects.TransactionVO,
	) (*paymentsresponsesdtos.CreateTransactionResponseDTO, error)

	GetTransaction(
		ctx context.Context,
		transactionID string,
	) (*paymentsentites.Transaction, error)

	CreateCustomer(
		ctx context.Context,
		vo *paddle.CreateCustomerRequest,
	)(res *paymentsvalueobjects.CustomerIDVO, err error)
}