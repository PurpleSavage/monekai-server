package paymentsoutadapters

import (
	"context"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"
	paymentsresponsesdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/reponses"
	paymentsports "github.com/PurpleSavage/monekai-server/modules/payments/application/ports"
	paymentsentites "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
	paymentsentities "github.com/PurpleSavage/monekai-server/modules/payments/domain/entities"
	paymentsenums "github.com/PurpleSavage/monekai-server/modules/payments/domain/enums"
	paymentsvalueobjects "github.com/PurpleSavage/monekai-server/modules/payments/domain/valueobjects"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
)

type PaymentServiceAdapter struct {
	paddleClient *paddle.SDK
}

func NewPaymentServiceAdapter() (paymentsports.PaymentServicePort, error) {

	client, err := paddle.NewSandbox(config.Envs.PaddleApiKey)
	// En producción:
	// client, err := paddle.New(config.Envs.PaddleApiKey)

	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Payment Service Error",
			"Unable to initialize Paddle client",
			err,
		)
	}

	return &PaymentServiceAdapter{
		paddleClient: client,
	}, nil
}
func (p *PaymentServiceAdapter) CreateCustomer(
	ctx context.Context,
	vo *paddle.CreateCustomerRequest,
) (res *paymentsvalueobjects.CustomerIDVO, err error) {
	customer, err := p.paddleClient.CreateCustomer(ctx, vo)
	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Payment Error",
			"Error to create customer",
			err,
		)
	}
	customerVO, err := paymentsvalueobjects.CreateCustomerIDVO(customer.ID)
	if err != nil {
		return nil, err
	}
	return customerVO, nil

}

func (p *PaymentServiceAdapter) CreateTransaction(
	ctx context.Context,
	vo *paymentsvalueobjects.TransactionVO,
) (*paymentsresponsesdtos.CreateTransactionResponseDTO, error) {

	transaction, err := p.paddleClient.TransactionsClient.CreateTransaction(
		ctx,
		&paddle.CreateTransactionRequest{
			CustomerID: &vo.CustomerID,

			Items: []paddle.CreateTransactionItems{
				*paddle.NewCreateTransactionItemsTransactionItemFromCatalog(
					&paddle.TransactionItemFromCatalog{
						PriceID:  vo.PriceID,
						Quantity: 1,
					},
				),
			},

			CustomData: paddle.CustomData{
				"user_id": vo.UserID.String(),
			},
		},
	)

	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Payment Error",
			"Error processing the transaction",
			err,
		)
	}
	return &paymentsresponsesdtos.CreateTransactionResponseDTO{
		TransactionID: transaction.ID,
		Status:        paymentsenums.TransactionStatus(transaction.Status),
	}, nil
}

func (p *PaymentServiceAdapter) GetTransaction(
	ctx context.Context,
	transactionID string,
) (*paymentsentites.Transaction, error) {
	transaction, err := p.paddleClient.TransactionsClient.GetTransaction(
		ctx,
		&paddle.GetTransactionRequest{
			TransactionID: transactionID,
		},
	)

	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Payment Error",
			"Unable to retrieve transaction",
			err,
		)
	}

	return &paymentsentities.Transaction{
		ID:       transaction.ID,
		Status:   paymentsenums.TransactionStatus(transaction.Status),
		PriceID:  transaction.Items[0].Price.ID,
		Currency: string(transaction.CurrencyCode),
	}, nil
}
