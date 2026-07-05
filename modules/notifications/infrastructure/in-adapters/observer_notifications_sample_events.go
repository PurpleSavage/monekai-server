package notificationsinadapters

import (
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)
type ObserverNotificationSampleEvent struct {
	
}

func NewObserverNotificationSampleEvent(
	
) commonports.Observer{
	return &ObserverNotificationSampleEvent{}
}

func (o *ObserverNotificationSampleEvent) ReceiveEvent(event commonentities.Event[any]) {
	
}