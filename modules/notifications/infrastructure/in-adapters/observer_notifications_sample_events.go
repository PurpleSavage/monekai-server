package notificationsinadapters

import (
	"encoding/json"
	"log"

	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerresponsessdtos "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/responsess"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)
type ObserverNotificationSampleEvent struct {
	saveNotificationUC *notificationsusecases.SaveNotificationUseCase
	sse                *notificationsevents.SSEManager
}

func NewObserverNotificationSampleEvent(
	sn *notificationsusecases.SaveNotificationUseCase,
	sse *notificationsevents.SSEManager,
) commonports.Observer {
	return &ObserverNotificationSampleEvent{
		saveNotificationUC: sn,
		sse:                sse,
	}
}

func (o *ObserverNotificationSampleEvent) ReceiveEvent(event commonentities.Event) {
	data, ok := event.Data.(samplerequestsdto.DataSampleNotify)
	if !ok {
		log.Println("invalid event data")
		return
	}

	var notifType, title, message string
	switch event.Name {
	case commonentities.EventSampleReady:
		notifType = string(notificationsenums.ReplicateSuccess)
		title = "Sample ready"
		message = "Your audio sample has been generated successfully"
	case commonentities.EventSampleError:
		notifType = string(notificationsenums.ReplicateError)
		title = "Sample failed"
		message = "Error in sample generation"
	default:
		log.Println("unhandled event name:", event.Name)
		return
	}

	notificationData, err := o.saveNotificationUC.Execute(
		data.UserID,
		notifType,
		title,
		message,
		data.SampleID,
	)
	if err != nil {
		log.Println("error saving notification: ", err)
		return
	}

	switch event.Name {
	case commonentities.EventSampleReady:
		dto := notificationssreponsesdtos.NotificationResponseDTO[samplerresponsessdtos.SampleResponseDTO]{
			Data:           data.Data,
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
		o.sse.BroadcastToUser(notificationData.Email, string(event.Name), string(payload))

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
		o.sse.BroadcastToUser(notificationData.Email, string(event.Name), string(payload))
	}
}