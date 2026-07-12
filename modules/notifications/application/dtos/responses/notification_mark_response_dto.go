package notificationssreponsesdtos

import (
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
)

type NotificationMarkResponseDTO struct {
	NotificationID string `json:"notificationId"`
	StatusNotification  notificationsenums.NotificationStatus`json:"statusNotification"`
}