package notificationsinadapters

import (
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)
type ObserverNotificationSampleEvent struct {
	saveNotificationUC *notificationsusecases.SaveNotificationUseCase
}

func NewObserverNotificationSampleEvent(
	sn *notificationsusecases.SaveNotificationUseCase,
) commonports.Observer{
	return &ObserverNotificationSampleEvent{
		saveNotificationUC: sn,
	}
}

func (o *ObserverNotificationSampleEvent) ReceiveEvent(event commonentities.Event[any]) {
	
}