package notificationsports

import (
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
	notificationsvalueobjects "github.com/PurpleSavage/monekai-server/modules/notifications/domain/valueobjects"
)

type NotificationsPersistencePort interface {
	SaveNotification(vo notificationsvalueobjects.SaveNotificationVO) (*notificationsentities.NotificationEntity, error )
}

