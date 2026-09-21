package account

import (
	accountusecases "github.com/PurpleSavage/monekai-server/modules/account/application/usecases"
	accountcontrollers "github.com/PurpleSavage/monekai-server/modules/account/infrastructure/controllers"
	accountoutadapters "github.com/PurpleSavage/monekai-server/modules/account/infrastructure/out-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AccountBootstrap(
	db *gorm.DB,
	am *authmiddlewares.AuthMiddleware,
) chi.Router {
	accountRepo := accountoutadapters.NewAccountRepository(db)
	listPaymentHistoryUC := accountusecases.NewListPaymentHistoryUseCase(accountRepo)

	controller := accountcontrollers.NewAccountController(
		am,
		listPaymentHistoryUC,
	)
	return accountcontrollers.AccountMapRoutes(controller)
}