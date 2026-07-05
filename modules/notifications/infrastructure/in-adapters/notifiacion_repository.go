package notificationsinadapters

import (
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
	notificationsvalueobjects "github.com/PurpleSavage/monekai-server/modules/notifications/domain/valueobjects"
	notificationsinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/mappers"
	notificationsraws "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/raws"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/gorm"
)

type NotificationsRepository struct {
	db *gorm.DB
}

func NewNotificationsRepository(db *gorm.DB) notificationsports.NotificationsPersistencePort {
	return &NotificationsRepository{db: db}
}

func (r *NotificationsRepository) SaveNotification(
	vo notificationsvalueobjects.SaveNotificationVO,
) (*notificationsentities.NotificationEntity, error ){
	
	newNoptification:= models.Notification{
		UserID:      vo.UserID,
		Type:        vo.Type,
		Title:       vo.Title,
		Message:     vo.Message,
		ReferenceID: vo.ReferenceID,	
	}	
	err := r.db.Create(&newNoptification).Error
	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Failed to save notification",
			"",
			nil,
		)
	}
	 var queryResult notificationsraws.NotificationQueryResultRaw
	err = r.db.Table("notifications").
		Select("notifications.*, users.email").
		Joins("INNER JOIN users ON users.id = notifications.user_id").
		Where("notifications.id = ?", newNoptification.ID).
		First(&queryResult).Error
	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Failed to retrieve notification details with user email",
			"",
			nil,
		)
	}
	savedNotificationResponse := notificationsinfrastructuremappers.ToNotificationEntity(newNoptification, queryResult.Email)
	return savedNotificationResponse, nil
}




