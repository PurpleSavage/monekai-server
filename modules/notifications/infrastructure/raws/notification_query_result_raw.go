package notificationsraws

import models "github.com/PurpleSavage/monekai-server/configurations/persistence"
type NotificationQueryResultRaw struct {
	models.Notification
	Email string `gorm:"column:email"` 
}