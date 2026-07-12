package notificationsusecases

import (
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
	"github.com/google/uuid"
)

type MarkMultiplesNotificationsAsReadUC struct {
	notificationPersistencePort notificationsports.NotificationsPersistencePort
}

func NewMarkMultiplesNotificationsAsReadUC(notificationPersistencePort notificationsports.NotificationsPersistencePort) *MarkMultiplesNotificationsAsReadUC {
	return &MarkMultiplesNotificationsAsReadUC{
		notificationPersistencePort: notificationPersistencePort,
	}
}

func (uc *MarkMultiplesNotificationsAsReadUC) Execute(
	notificationIDs []uuid.UUID,
) ([]notificationssreponsesdtos.NotificationMarkResponseDTO, error) {
	response, err := uc.notificationPersistencePort.MarkAllNotificationsAsRead(notificationIDs)
	if err != nil {
		return nil, err
	}
	return response, nil
}

