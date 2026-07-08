package notificationssreponsesdtos

import notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"

type NotificationDetailsResponseDTO struct {
	Data  			*any 	`json:"data"`
	NotificationID          string `json:"id"`
	Type        	notificationsenums.TypeNotification `json:"type"`
	Title       	string `json:"title"`
	Message       	string `json:"message"`
	Status       	string `json:"status"`
	ReferenceID 	string `json:"referenceId"`
	CreatedAt   	string `json:"createdAt"`
	UserID      	string `json:"userId"`
	Email       	string `json:"email"`
}
