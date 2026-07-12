package notificationsrequestsdtos

type MarkAllNotificationsAsReadRequestDTO struct {
	NotificationIds []string `json:"notificationIds" validate:"required,min=1,dive,uuid4"`
}