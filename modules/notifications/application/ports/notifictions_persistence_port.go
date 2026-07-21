package notificationsports

import (
	"context"

	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
	notificationsvalueobjects "github.com/PurpleSavage/monekai-server/modules/notifications/domain/valueobjects"
	"github.com/google/uuid"
)

type NotificationsPersistencePort interface {
	CountTotalNotifications(ctx context.Context, userID string) (int, error)
	SaveNotification(vo notificationsvalueobjects.SaveNotificationVO) (*notificationsentities.NotificationEntity, error)
	ListNotifications(ctx context.Context,userID string, limit int, page int) ([]notificationssreponsesdtos.ItemNotificationDTO, error)
	MarkAllNotificationsAsRead(userID string,notificationIDs []uuid.UUID) ([]notificationssreponsesdtos.NotificationMarkResponseDTO, error)
	MarkNotificationAsRead(userID string, notificationID uuid.UUID) (*notificationssreponsesdtos.NotificationMarkResponseDTO, error) // <-- Corregido a singular
}
