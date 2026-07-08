package notificationssreponsesdtos

import (
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	"github.com/google/uuid"
)

type ItemNotificationDTO struct {
	ID          uuid.UUID
	Type        notificationsenums.TypeNotification    
	Title       string    
	Message     string    
	Status      notificationsenums.NotificationStatus   
	ReferenceID string   
	CreatedAt   string
	UserID      uuid.UUID
}
