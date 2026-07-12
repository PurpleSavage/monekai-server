package notificationsusecases

import (
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
)

type MarkNotificationAsReadUseCase struct {
	notificationRepository  notificationsports.NotificationsPersistencePort
}
func NewMarkNotificationAsReadUseCase(
	notificationRepository notificationsports.NotificationsPersistencePort,
) *MarkNotificationAsReadUseCase {
	return &MarkNotificationAsReadUseCase{
		notificationRepository: notificationRepository,
	}
}

func (uc *MarkNotificationAsReadUseCase) Execute(notificationID string)(*notificationssreponsesdtos.NotificationMarkResponseDTO,error) {
	parseId , err:= authvalueobjects.NewUUIDVO(notificationID)
	if err != nil {
		return nil, err
	}
	response, err := uc.notificationRepository.MarkNotificationAsRead(parseId.Value())
	if err != nil {
		return nil, err
	}
	return response, nil
}
