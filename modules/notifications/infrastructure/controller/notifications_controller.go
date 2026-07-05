package notificationscontroller

import (
	"net/http"

	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/go-chi/chi/v5"
)

type NotificationsController struct {
	authMiddleware *authmiddlewares.AuthMiddleware
}

func NewNotificationsController(
	am *authmiddlewares.AuthMiddleware,
) *NotificationsController {
	return &NotificationsController{
		authMiddleware: am,
	}
}

func (c *NotificationsController) ListNotifications(w http.ResponseWriter, r *http.Request) {
	
}

func NotificationsMaproutes(nc *NotificationsController) chi.Router {
	r := chi.NewRouter()
	r.Use(nc.authMiddleware.AccessToken)
	r.Get("/", nc.ListNotifications)
	return r
}