package notificationssreponsesdtos

import (
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
)

type NotificationResponseDTO[T any] struct {
	Data *T `json:"data"`
	NotificationID          string `json:"id"`
	Type        notificationsenums.TypeNotification `json:"type"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	ReferenceID string `json:"reference_id"`
	CreatedAt   string `json:"created_at"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
}