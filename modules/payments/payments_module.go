package payments

import (
	paymentsusecases "github.com/PurpleSavage/monekai-server/modules/payments/application/usecases"
	paymentscontroller "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/controllers"
	paymentsoutadapters "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/out-adapters"
	authusecases "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/usecases"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	authoutadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/out-dapters"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)
func SamplerBootstrap(
	db *gorm.DB,
	ob commonports.ObserverBucketPort,
	v *validators.DTOValidator,
 	authmiddleware *authmiddlewares.AuthMiddleware,
) chi.Router{
	userRepo := authoutadapters.NewUserRepository(db)
	paymentRepo:= paymentsoutadapters.NewPaymentRepository(db)
	paymentService,_:= paymentsoutadapters.NewPaymentServiceAdapter()

	findUserByEmailUC:= authusecases.NewFindUserByEmailUseCase(userRepo)
	createPaymentUC:= paymentsusecases.NewCreatePaymentUseCase(
		paymentRepo,
		paymentService,
		findUserByEmailUC,
	)
	listCreditsUC:= paymentsusecases.NewListCreditPackageUseCase(paymentRepo)
	
	controller := paymentscontroller.NewPaymentsController(
		v,
		authmiddleware,
		createPaymentUC,
		listCreditsUC,
	)
	return paymentscontroller.PaymentsMapRoutes(controller)
}