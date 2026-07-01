package samplervalueobjects

import (
	"strings"
	audioenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	"github.com/google/uuid"
)

type SaveNotificationVO struct {
	UserID      uuid.UUID 
	Type        audioenums.TypeNotification
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
	var finalType audioenums.TypeNotification
	switch audioenums.TypeNotification(strings.ToLower(notificationType)) {
		case audioenums.ReplicateError:
			finalType = audioenums.ReplicateError
		case audioenums.ReplicateSuccess:
			finalType = audioenums.ReplicateSuccess
		case audioenums.Payment:
			finalType = audioenums.Payment
		case audioenums.Info:
			finalType = audioenums.Info
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