package payments

import (
	paymentsusecases "github.com/PurpleSavage/monekai-server/modules/payments/application/usecases"
	paymentscontroller "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/controllers"
	paymentsmiddlewares "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/middlewares"
	paymentsoutadapters "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/out-adapters"
	authusecases "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/usecases"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	authoutadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/out-dapters"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)
func PaymentsBootstrap(
	db *gorm.DB,
	ob commonports.ObserverBucketPort,
	v *validators.DTOValidator,
 	authmiddleware *authmiddlewares.AuthMiddleware,
) chi.Router{
	userRepo := authoutadapters.NewUserRepository(db)
	paymentRepo:= paymentsoutadapters.NewPaymentRepository(db)
	paymentService,_:= paymentsoutadapters.NewPaymentServiceAdapter()

	//middlewarespayments
	paymentVirifyMiddleware:=paymentsmiddlewares.NewPaymentWebhookVerifier()

	
	findUserByEmailUC:= authusecases.NewFindUserByEmailUseCase(userRepo)
	createPaymentUC:= paymentsusecases.NewCreatePaymentUseCase(
		paymentRepo,
		paymentService,
		findUserByEmailUC,
	)
	listCreditsUC:= paymentsusecases.NewListCreditPackageUseCase(paymentRepo)
	processPaymentWebhookUC:= paymentsusecases.NewProcessPaymentWebhookUseCase(
		paymentRepo,
		userRepo,
	)
	
	controller := paymentscontroller.NewPaymentsController(
		v,
		authmiddleware,
		createPaymentUC,
		listCreditsUC,
		paymentVirifyMiddleware,
		processPaymentWebhookUC,
		ob,
	)
	return paymentscontroller.PaymentsMapRoutes(controller)
}