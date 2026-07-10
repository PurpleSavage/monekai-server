package notificationsinadapters

import (
	"encoding/json"
	"log"

	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)
type ObserverNotificationSampleEvent struct {
	saveNotificationUC *notificationsusecases.SaveNotificationUseCase
	sse *notificationsevents.SSEManager
}

func NewObserverNotificationSampleEvent(
	sn *notificationsusecases.SaveNotificationUseCase,
	sse *notificationsevents.SSEManager,
) commonports.Observer{
	return &ObserverNotificationSampleEvent{
		saveNotificationUC: sn,
		sse: sse,
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
	switch event.Name {

		case commonentities.EventSampleReady:

			sampleData, ok := (*data.Data).(samplerresponsessdtos.SampleResponseDTO)
			if !ok {
				log.Println("invalid sample data")
				return
			}

			dto := notificationssreponsesdtos.NotificationResponseDTO[samplerresponsessdtos.SampleResponseDTO]{
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

			payload, err := json.Marshal(dto)
			if err != nil {
				log.Println(err)
				return
			}

			o.sse.BroadcastToUser(
				data.Email,
				string(event.Name),
				string(payload),
			)

		case commonentities.EventSampleError:

			dto := notificationssreponsesdtos.NotificationResponseDTO[struct{}]{
				Data:           nil,
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

			payload, err := json.Marshal(dto)
			if err != nil {
				log.Println(err)
				return
			}

			o.sse.BroadcastToUser(
				data.Email,
				string(event.Name),
				string(payload),
			)
	}
} 