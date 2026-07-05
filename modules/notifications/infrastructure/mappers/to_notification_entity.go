package notificationsinfrastructuremappers

import (
	"time"

	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
)

func ToNotificationEntity(model models.Notification, email string) *notificationsentities.NotificationEntity {
	return &notificationsentities.NotificationEntity{
		ID:          model.ID.String(),
		Type:        model.Type,
		Title:       model.Title,
		Message:     model.Message,
		Status:      model.Status,
		ReferenceID: model.ReferenceID.String(),
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		UserID: model.UserID.String(),
		Email:email,
	}
}