package notificationscontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	notificationsrequestsdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/requests"
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/common/application/dtos/requests"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NotificationsController struct {
	authMiddleware *authmiddlewares.AuthMiddleware
	validator           *validators.DTOValidator
	listNotificationsUC *notificationsusecases.ListNotificationsUseCase
	sseManager          *notificationsevents.SSEManager
	markAllNotificationsAsRead *notificationsusecases.MarkMultiplesNotificationsAsReadUC
	markNotificationAsRead *notificationsusecases.MarkNotificationAsReadUseCase
}

func NewNotificationsController(
	am *authmiddlewares.AuthMiddleware,
	v *validators.DTOValidator,
	listNotificationsUC *notificationsusecases.ListNotificationsUseCase,
	sseManager *notificationsevents.SSEManager,
	markAllNotificationsAsRead *notificationsusecases.MarkMultiplesNotificationsAsReadUC,
	markNotificationAsRead *notificationsusecases.MarkNotificationAsReadUseCase,
) *NotificationsController {
	return &NotificationsController{
		authMiddleware:      am,
		validator:           v,
		listNotificationsUC: listNotificationsUC,
		sseManager:          sseManager,
		markAllNotificationsAsRead: markAllNotificationsAsRead,
		markNotificationAsRead:     markNotificationAsRead,
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
		r.Context(),
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

	var body notificationsrequestsdtos.MarkAllNotificationsAsReadRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(http.StatusBadRequest, "Bad Request", "Invalid JSON body", err),
		)
		return
	}
	defer r.Body.Close()

	if err := c.validator.ValidateStruct(body); err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	validUUIDs := make([]uuid.UUID, 0, len(body.NotificationIds))
	for _, id := range body.NotificationIds {
		parsedID, err := authvalueobjects.NewUUIDVO(id)
		if err != nil {
			commoninfrastructuremappers.RespondWithError(w, err)
			return
		}
		validUUIDs = append(validUUIDs, parsedID.Value())
	}

	// Pasamos también el userID para que el use case/repo filtre por dueño
	response, err := c.markAllNotificationsAsRead.Execute(dtoSession.Id, validUUIDs)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}


//controlaador para marcar una notificación como leída
func (c *NotificationsController) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
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

	var body notificationsrequestsdtos.MarkNotificationAsReadDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(http.StatusBadRequest, "Bad Request", "Invalid JSON body", err),
		)
		return
	}
	defer r.Body.Close()

	uuidValid, err := authvalueobjects.NewUUIDVO(body.NotificationId)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	response, err := c.markNotificationAsRead.Execute(dtoSession.Id, uuidValid.String())
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}

	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}

func NotificationsMaproutes(nc *NotificationsController) chi.Router {
	r := chi.NewRouter()
	r.Use(nc.authMiddleware.AccessToken)
	r.Get("/all", nc.ListNotifications)
	r.Patch("/read-all", nc.MarkAllNotificationsAsRead)
	r.Patch("/{id}/read", nc.MarkNotificationAsRead)
	return r
}