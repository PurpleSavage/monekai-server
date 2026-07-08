package notificationsinfrastructuremappers

import (
	"time"

	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
)

func ToListOfItmesNotification(notifications []models.Notification) []notificationssreponsesdtos.ItemNotificationDTO {
	var result []notificationssreponsesdtos.ItemNotificationDTO
	for _, notification := range notifications {
		result = append(result, notificationssreponsesdtos.ItemNotificationDTO{
			ID:          notification.ID,
			Type:        notification.Type,
			Title:       notification.Title,
			Message:     notification.Message,
			Status:      notification.Status,
			ReferenceID: notification.ReferenceID.String(),
			CreatedAt:   notification.CreatedAt.Format(time.RFC3339),
			UserID:      notification.UserID,
		})
	}
	return result
}
