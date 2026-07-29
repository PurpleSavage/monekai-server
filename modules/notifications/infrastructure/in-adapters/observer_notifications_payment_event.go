package notificationsinadapters

import (
	"encoding/json"
	"log"

	paymentsrequestsdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/requests"
	paymentsresponsesdtos "github.com/PurpleSavage/monekai-server/modules/payments/application/dtos/reponses"
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsusecases "github.com/PurpleSavage/monekai-server/modules/notifications/application/usecases"
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)

type ObserverNotificationPaymentEvent struct {
	saveNotificationUC *notificationsusecases.SaveNotificationUseCase
	sse                *notificationsevents.SSEManager
}

func NewObserverNotificationPaymentEvent(
	sn *notificationsusecases.SaveNotificationUseCase,
	sse *notificationsevents.SSEManager,
) commonports.Observer {
	return &ObserverNotificationPaymentEvent{
		saveNotificationUC: sn,
		sse:                sse,
	}
}

func (o *ObserverNotificationPaymentEvent) ReceiveEvent(event commonentities.Event) {
	data, ok := event.Data.(paymentsrequestsdtos.DataPaymentNotify)
	if !ok {
		log.Println("invalid payment event data")
		return
	}

	var notifType, title, message string
	switch event.Name {
	case commonentities.EventPaymentSuccess:
		notifType = string(notificationsenums.Payment)
		title = "Payment successful"
		message = "Credits have been added to your account"
	case commonentities.EventPaymentFailed:
		notifType = string(notificationsenums.Payment)
		title = "Payment failed"
		message = "The payment could not be processed"
	default:
		log.Println("unhandled payment event name:", event.Name)
		return
	}

	notificationData, err := o.saveNotificationUC.Execute(
		data.UserID,
		notifType,
		title,
		message,
		data.PaymentID,
	)
	if err != nil {
		log.Println("error saving payment notification: ", err)
		return
	}

	dto := notificationssreponsesdtos.NotificationResponseDTO[paymentsresponsesdtos.PaymentNotificationDTO]{
		Data: &paymentsresponsesdtos.PaymentNotificationDTO{
			PaymentID: data.PaymentID,
			Credits:   data.Credits,
			Amount:    data.Amount,
			Currency:  data.Currency,
			Status:    data.Status,
		},
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
