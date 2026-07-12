package notifications

import (
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationscontroller "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/controller"
	notificationsinadapters "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/in-adapters"
	notificationsoutadapters "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/out-adapters"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func NotificationsBootstrap(
	db *gorm.DB,
	ob commonports.ObserverBucketPort,
	v *validators.DTOValidator,
	am *authmiddlewares.AuthMiddleware,
	sseManager *notificationsevents.SSEManager,
) chi.Router{
	
	//adapters 
	notificationsRepo:= notificationsoutadapters.NewNotificationsRepository(db)

	//usecases
	saveNotificationUC := notificationsusecases.NewSaveNotificationUseCase(notificationsRepo)
	listNotificationsUC := notificationsusecases.NewListNotificationsUseCase(notificationsRepo)
	markAllNotificationsAsReadUC := notificationsusecases.NewMarkMultiplesNotificationsAsReadUseCase(notificationsRepo)
	markNotificationAsReadUC := notificationsusecases.NewMarkNotificationAsReadUseCase(notificationsRepo)

	//adapters -depent - uc
	notificationsObserver:= notificationsinadapters.NewObserverNotificationSampleEvent(
		saveNotificationUC,
		sseManager,
	)
	ob.AddObserver(notificationsObserver, "sample_event")


	
	controller:= notificationscontroller.NewNotificationsController(
		am,
		v,
		listNotificationsUC,
		sseManager,
		markAllNotificationsAsReadUC,
		markNotificationAsReadUC,
	)
	return  notificationscontroller.NotificationsMaproutes(controller)
}