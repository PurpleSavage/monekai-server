package paymentscontroller

import (
	"net/http"

	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)


type PaymentsController struct {
	validator  *validators.DTOValidator
	authMiddleware *authmiddlewares.AuthMiddleware
}
func NewPaymentsController(
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
) *PaymentsController {
	return &PaymentsController{
		validator:v,
		authMiddleware: am,
	}
}
func (sc *PaymentsController) HandlePayments(w http.ResponseWriter, r *http.Request) {

}
func (sc *PaymentsController) ReceivePaymentWebhook(w http.ResponseWriter, r *http.Request){
	
}

func PaymentsMapRoutes(sc *PaymentsController) chi.Router{
	r := chi.NewRouter()
	r.Use(sc.authMiddleware.AccessToken)
	r.Post("/create", sc.HandlePayments)
	r.Post("/webhook", sc.ReceivePaymentWebhook)
	return r
}