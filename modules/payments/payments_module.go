package payments

import (
	paymentscontroller "github.com/PurpleSavage/monekai-server/modules/payments/infrastructure/controllers"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
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
	controller := paymentscontroller.NewPaymentsController(
		v,
		authmiddleware,
	)
	return paymentscontroller.PaymentsMapRoutes(controller)
}