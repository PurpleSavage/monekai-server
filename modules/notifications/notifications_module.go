package notifications

import (
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationscontroller "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/controller"
	notificationsinadapters "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/in-adapters"
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
) chi.Router{
	
	//adapters 
	notificationsRepo:= notificationsinadapters.NewNotificationsRepository(db)

	//usecases
	saveNotificationUC := notificationsusecases.NewSaveNotificationUseCase(notificationsRepo)
	listNotificationsUC := notificationsusecases.NewListNotificationsUseCase(notificationsRepo)
	
	//adapters -depent - uc
	notificationsObserver:= notificationsinadapters.NewObserverNotificationSampleEvent(saveNotificationUC)
	ob.AddObserver(notificationsObserver, "sample_event")


	
	controller:= notificationscontroller.NewNotificationsController(
		am,
		v,
		listNotificationsUC,
	)
	return  notificationscontroller.NotificationsMaproutes(controller)
}