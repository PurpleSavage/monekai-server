package notificationscontroller

import (
	"net/http"
	"strconv"
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/common/application/dtos/requests"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
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
//controlador para listar las notificaciones con paginado
func (c *NotificationsController) ListNotifications(w http.ResponseWriter, r *http.Request) {
	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
	if rawData == nil {
		commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil))
		return
	}
	dtoSession, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(500, "Internal Error", "Could not parse session data", err))
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(http.StatusBadRequest, "Bad Request", "'page' must be a valid integer", err),
		)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(http.StatusBadRequest, "Bad Request", "'limit' must be a valid integer", err),
		)
		return
	}

	queryDto := commonrequestsdtos.ListNotificationsQueryDTO{
		Page:  page,
		Limit: limit,
	}

	notifications, err:= c.listNotificationsUC.Execute(
		dtoSession.Id,
		queryDto.Limit,
		queryDto.Page,
	)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, notifications)
}


//controlador para marcr todas las notificaciones no leídas como leídas
func (c *NotificationsController) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	
}


//controlaador para marcar una notificación como leída
func (c *NotificationsController) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	
}

func NotificationsMaproutes(nc *NotificationsController) chi.Router {
	r := chi.NewRouter()
	r.Use(nc.authMiddleware.AccessToken)
	r.Get("/", nc.ListNotifications)
	return r
}