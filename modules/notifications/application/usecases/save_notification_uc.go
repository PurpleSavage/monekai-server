package notificationsusecases

import (
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
	notificationsvalueobjects "github.com/PurpleSavage/monekai-server/modules/notifications/domain/valueobjects"
)

type SaveNotificationUseCase struct {
	notificationRepo notificationsports.NotificationsPersistencePort
}

func NewSaveNotificationUseCase(
	repo notificationsports.NotificationsPersistencePort,
) *SaveNotificationUseCase {
	return &SaveNotificationUseCase{
		notificationRepo: repo,
	}
}

func (uc *SaveNotificationUseCase) Execute(
	userID string,
	notificationType string,
	title string,
	message string,
	referenceID string,
) (*notificationsentities.NotificationEntity, error){
	vo, err := notificationsvalueobjects.CreateSaveNotificationVO(
		userID,
		notificationType,
		title,
		message,
		referenceID,
	)
	if err != nil {
		return nil, err
	}
	response, error := uc.notificationRepo.SaveNotification(*vo)
	if error != nil {
		return nil, error
	}
	return response, nil
}
