package notificationscontroller

import (
	"net/http"

	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
)

type NotificationsController struct {
	authMiddleware *authmiddlewares.AuthMiddleware
	validator           *validators.DTOValidator
	listNotificationsUC *notificationsusecases.ListNotificationsUseCase
	sseManager          *notificationsevents.SSEManager
}

func NewNotificationsController(
	am *authmiddlewares.AuthMiddleware,
	v *validators.DTOValidator,
	listNotificationsUC *notificationsusecases.ListNotificationsUseCase,
	sseManager *notificationsevents.SSEManager,
) *NotificationsController {
	return &NotificationsController{
		authMiddleware:      am,
		validator:           v,
		listNotificationsUC: listNotificationsUC,
		sseManager:          sseManager,
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