package paymentsusecases

import (
	"context"

	"github.com/PaddleHQ/paddle-go-sdk/v5"
	paymentsresponsesdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/reponses"
	paymentsports "github.com/PurpleSavage/monekai-server/modules/payments/application/ports"
	paymentsvalueobjects "github.com/PurpleSavage/monekai-server/modules/payments/domain/valueobjects"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
)

type CreatePaymentUseCase struct{
	paymentRepo paymentsports.PaymentsPersistencePort
	userRepo 	authports.UserPersistencePort
	paymentService paymentsports.PaymentServicePort
}
func NewCreatePaymentUseCase(
	paymentRepo paymentsports.PaymentsPersistencePort,
	userRepo 	authports.UserPersistencePort,
) *CreatePaymentUseCase{
	return  &CreatePaymentUseCase{
		paymentRepo: paymentRepo,
		userRepo:userRepo,
	}
}
func (c *CreatePaymentUseCase) Execute(
	packageId string,
	userID string,
	email string,
	ctx context.Context,
)(*paymentsresponsesdtos.CreateTransactionResponseDTO, error){
	user,err:=c.userRepo.FindUserByEmail(email)
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