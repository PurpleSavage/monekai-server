package notificationsinadapters

import (
	"log"
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
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

func (o *ObserverNotificationSampleEvent) ReceiveEvent(
	event commonentities.Event,
) {
	data, ok := event.Data.(notificationssreponsesdtos.NotificationDetailsResponseDTO )
	if !ok {
		log.Println("invalid event data")
		return
	}
	notificationData,err:=o.saveNotificationUC.Execute(
		data.UserID,
		string(data.Type),
		data.Title,
		data.Message,
		data.ReferenceID,
	)
	if err!=nil {
		//manejar error
		log.Println("error saving notification: ",err)
		return
	}
	sampleData, ok := (*data.Data).(samplerresponsessdtos.SampleResponseDTO)
	if !ok {
			log.Println("invalid notification data")
			return
	}
	notificationDto := notificationssreponsesdtos.NotificationResponseDTO[samplerresponsessdtos.SampleResponseDTO]{
		Data:           &sampleData, 
		NotificationID: notificationData.ID,
		Type:           notificationData.Type,
		Title:          notificationData.Title,
		Message:        notificationData.Message,
		Status:         string(notificationData.Status),
		ReferenceID:    notificationData.ReferenceID,
		CreatedAt:      notificationData.CreatedAt,
		UserID:         notificationData.UserID,
		Email:          notificationData.Email,
	}
}