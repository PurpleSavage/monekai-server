package notificationsentities

import notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
type NotificationEntity struct {
	ID          string
	Type        notificationsenums.TypeNotification    
	Title       string    
	Message     string    
	Status      notificationsenums.NotificationStatus   
	ReferenceID string   
	CreatedAt   string
	UserID      string
	Email       string
}