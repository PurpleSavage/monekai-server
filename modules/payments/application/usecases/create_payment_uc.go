package paymentsusecases

import (
	"context"

	"github.com/PaddleHQ/paddle-go-sdk/v5"
	paymentsresponsesdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/reponses"
	paymentsports "github.com/PurpleSavage/monekai-server/modules/payments/application/ports"
	paymentsvalueobjects "github.com/PurpleSavage/monekai-server/modules/payments/domain/valueobjects"
	authusecases "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/usecases"
)

type CreatePaymentUseCase struct{
	paymentRepo paymentsports.PaymentsPersistencePort

	paymentService paymentsports.PaymentServicePort
	findUserByEmailUC *authusecases.FindUserByEmailUseCase
}
func NewCreatePaymentUseCase(
	paymentRepo paymentsports.PaymentsPersistencePort,
	paymentService paymentsports.PaymentServicePort,
	findUserByEmailUC *authusecases.FindUserByEmailUseCase,
) *CreatePaymentUseCase{
	return  &CreatePaymentUseCase{
		paymentRepo: paymentRepo,
		paymentService: paymentService,
		findUserByEmailUC: findUserByEmailUC,
	}
}
func (c *CreatePaymentUseCase) Execute(
	packageId string,
	userID string,
	email string,
	ctx context.Context,
)(*paymentsresponsesdtos.CreateTransactionResponseDTO, error){
	user,err:=c.findUserByEmailUC.Execute(email)
	if err!=nil{
		return nil,err
	}
	var customerID string

	if user.CustomerID == nil{
	 	vo:=&paddle.CreateCustomerRequest{
	    	Email: email,
		}
		responseVo,err:=c.paymentService.CreateCustomer(ctx,vo)
		if err != nil {
    		return nil, err
		}
		customerID=responseVo.Value()
	}
	packageCredits,err:= c.paymentRepo.GetCreditPackage(packageId)
	if err!=nil{
		return nil,err
	}
	vo,err:= paymentsvalueobjects.CreateTransactionVO(
		packageCredits.PriceId,
		userID,
		customerID,	
	)
	response,err:=c.paymentService.CreateTransaction(ctx,vo)
	if err != nil {
		return nil,err
	}
	return response,nil
}