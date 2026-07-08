package notificationsusecases

import (
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
)

type ListNotificationsUseCase struct {
	notificationsRepo notificationsports.NotificationsPersistencePort
}

func NewListNotificationsUseCase(
	notificationsRepo notificationsports.NotificationsPersistencePort,
) *ListNotificationsUseCase {
	return &ListNotificationsUseCase{
		notificationsRepo: notificationsRepo,
	}
}

func (uc *ListNotificationsUseCase) Execute(userID string,limit int, page int) ([]notificationssreponsesdtos.ItemNotificationDTO, error) {
	notifications, err := uc.notificationsRepo.ListNotifications(userID, limit, page)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
