package notificationsvalueobjects

import (
	"strings"

	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	"github.com/google/uuid"
)

type SaveNotificationVO struct {
	UserID      uuid.UUID 
	Type        notificationsenums.TypeNotification
	Title       string
	Message     string
	ReferenceID uuid.UUID
}

func CreateSaveNotificationVO(
	userID string,
	notificationType string,
	title string,
	message string,
	referenceID string,
) (*SaveNotificationVO, error) {

	// 1. Validar requerimiento y formato de UserID con Regex estándar

	userUUID, err := authvalueobjects.NewUUIDVO(userID)
	if err != nil {
		return nil, err
	}

	// 2. Validar y transformar el Enum TypeNotification
	var finalType notificationsenums.TypeNotification
	switch notificationsenums.TypeNotification(strings.ToLower(notificationType)) {
		case notificationsenums.ReplicateError:
			finalType = notificationsenums.ReplicateError
			case notificationsenums.ReplicateSuccess:
			finalType = notificationsenums.ReplicateSuccess
			case notificationsenums.Payment:
			finalType = notificationsenums.Payment
			case notificationsenums.Info:
			finalType = notificationsenums.Info
		default:
			return nil, commondomainerrors.NewValidationError(
				"type",
				"unsupported notification type",
		)
	}

	// -------------------------
	// TITLE VALIDATION
	// -------------------------

	title = strings.TrimSpace(title)

	if title == "" {
		return nil, commondomainerrors.NewValidationError(
			"title",
			"title is required",
		)
	}

	// -------------------------
	// MESSAGE VALIDATION
	// -------------------------

	message = strings.TrimSpace(message)

	if message == "" {
		return nil, commondomainerrors.NewValidationError(
			"message",
			"message is required",
		)
	}

	referenceUUID, err := authvalueobjects.NewUUIDVO(referenceID)
	if err != nil {
		return nil, err
	}

	return &SaveNotificationVO{
		UserID:      userUUID.Value(),
		Type:        finalType,
		Title:       title,
		Message:     message,
		ReferenceID: referenceUUID.Value(),
	}, nil
}