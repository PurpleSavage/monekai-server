package notificationsports

import (
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
	notificationsvalueobjects "github.com/PurpleSavage/monekai-server/modules/notifications/domain/valueobjects"
)

type NotificationsPersistencePort interface {
	SaveNotification(vo notificationsvalueobjects.SaveNotificationVO) (*notificationsentities.NotificationEntity, error )
	ListNotifications(userID string,limit int,page int,) ([]notificationssreponsesdtos.ItemNotificationDTO, error)
}

