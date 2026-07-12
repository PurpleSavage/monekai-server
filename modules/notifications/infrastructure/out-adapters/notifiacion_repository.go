package notificationsoutadapters

import (
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
	notificationsentities "github.com/PurpleSavage/monekai-server/modules/notifications/domain/entities"
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	notificationsvalueobjects "github.com/PurpleSavage/monekai-server/modules/notifications/domain/valueobjects"
	notificationsinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/mappers"
	notificationsraws "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/raws"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NotificationsRepository struct {
	db *gorm.DB
}

func NewNotificationsRepository(db *gorm.DB) notificationsports.NotificationsPersistencePort {
	return &NotificationsRepository{db: db}
}

func (r *NotificationsRepository) SaveNotification(
	vo notificationsvalueobjects.SaveNotificationVO,
) (*notificationsentities.NotificationEntity, error) {
	newNotification := models.Notification{
		UserID:      vo.UserID,
		Type:        vo.Type,
		Title:       vo.Title,
		Message:     vo.Message,
		ReferenceID: vo.ReferenceID,
	}
	
	err := r.db.Create(&newNotification).Error
	if err != nil {
		return nil, globalerrors.NewAppError(500, "Failed to save notification", "", nil)
	}

	var queryResult notificationsraws.NotificationQueryResultRaw
	err = r.db.Table("notifications").
		Select("notifications.*, users.email").
		Joins("INNER JOIN users ON users.id = notifications.user_id").
		Where("notifications.id = ?", newNotification.ID).
		First(&queryResult).Error
	if err != nil {
		return nil, globalerrors.NewAppError(500, "Failed to retrieve notification details with user email", "", nil)
	}

	savedNotificationResponse := notificationsinfrastructuremappers.ToNotificationEntity(newNotification, queryResult.Email)
	return savedNotificationResponse, nil
}

func (r *NotificationsRepository) ListNotifications(
	userID string,
	limit int,
	page int,
) ([]notificationssreponsesdtos.ItemNotificationDTO, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&notifications).Error
	if err != nil {
		return nil, globalerrors.NewAppError(500, "Failed to list notifications", "", nil)
	}
	
	return notificationsinfrastructuremappers.ToListOfItmesNotification(notifications), nil
}

func (r *NotificationsRepository) MarkAllNotificationsAsRead(
	notificationIDs []uuid.UUID,
) ([]notificationssreponsesdtos.NotificationMarkResponseDTO, error) {
	if len(notificationIDs) == 0 {
		return []notificationssreponsesdtos.NotificationMarkResponseDTO{}, nil
	}
	
	var updatedNotifications []models.Notification
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Returning{}).
			Model(&models.Notification{}).
			Where("id IN ? AND status = ?", notificationIDs, notificationsenums.NotificationUnread).
			Update("status", notificationsenums.NotificationRead).
			Find(&updatedNotifications).Error	
		return err
	})

	if err != nil {
		return nil, globalerrors.NewAppError(500, "Failed to update notifications status", "", nil)
	} 

	updatedMap := make(map[string]bool)
	for _, n := range updatedNotifications {
		updatedMap[n.ID.String()] = true
	}

	response := make([]notificationssreponsesdtos.NotificationMarkResponseDTO, len(notificationIDs))
	for i, id := range notificationIDs {
		idStr := id.String()

		if updatedMap[idStr] {
			response[i] = notificationssreponsesdtos.NotificationMarkResponseDTO{
				NotificationID:     idStr,
				StatusNotification: notificationsenums.NotificationRead,
			}
		} else {
			response[i] = notificationssreponsesdtos.NotificationMarkResponseDTO{
				NotificationID:     idStr,
				StatusNotification: notificationsenums.NotificationUnread,
			}
		}
	}

	return response, nil
} 

func (r *NotificationsRepository) MarkNotificationAsRead(
	notificationID uuid.UUID,
) (*notificationssreponsesdtos.NotificationMarkResponseDTO, error) {
	// 1. El constructor de tu VO requiere un string, así que pasamos notificationID.String()
	uuidVO, err := authvalueobjects.NewUUIDVO(notificationID.String())
	if err != nil {
		return nil, err
	}

	var updatedNotification models.Notification

	// 2. Extraemos el tipo nativo uuid.UUID usando .Value() para la consulta de GORM
	err = r.db.Clauses(clause.Returning{}).
		Model(&models.Notification{}).
		Where("id = ? AND status = ?", uuidVO.Value(), notificationsenums.NotificationUnread).
		Update("status", notificationsenums.NotificationRead).
		Find(&updatedNotification).Error

	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Failed to update notification status",
			"",
			nil,
		)
	}

	if updatedNotification.ID == uuid.Nil {
		return &notificationssreponsesdtos.NotificationMarkResponseDTO{
			NotificationID:     uuidVO.String(),
			StatusNotification: notificationsenums.NotificationUnread,
		}, nil
	}

	return &notificationssreponsesdtos.NotificationMarkResponseDTO{
		NotificationID:     updatedNotification.ID.String(),
		StatusNotification: notificationsenums.NotificationRead,
	}, nil
}