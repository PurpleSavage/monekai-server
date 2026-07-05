package models

import (
	"time"
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ---SESSIONS---
type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"` // Genera el UUID automáticamente
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	RefreshToken string    `gorm:"uniqueIndex;not null"`
	UserAgent    string    `gorm:"size:255"`
	CreatedAt    time.Time `gorm:"autoCreateTime"` // GORM o la BD lo maneja al crear
	UpdatedAt    time.Time `gorm:"autoUpdateTime"` // GORM o la BD lo maneja al actualizar
}

// --- USUARIO ---
type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ExternalID    string    `gorm:"uniqueIndex;not null"`
	Email         string    `gorm:"unique;not null"`
	Credits       int       `gorm:"default:0"`
	CreatedAt     time.Time
	PhotoUrl      *string        `gorm:"type:text"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
}

// --- GENERACIÓN IA (REPLICATE) ---
type Sample struct {
	ID              uuid.UUID               `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          uuid.UUID               `gorm:"type:uuid;not null;index"`
	SampleName      string                  `gorm:"size:100;not null"`
	Prompt          string                  `gorm:"type:text"`
	InitialAudioURL *string                 `gorm:"type:text"`
	Duration        int                     `gorm:"type:text"`
	ModelVersion    samplerenums.ModelVersion `gorm:"type:varchar(30);not null;check:model_version IN ('stereo-melody-large', 'stereo-large', 'melody-large', 'large')"`
	OutputFormat    samplerenums.OutputFormat `gorm:"type:varchar(10);not null;check:output_format IN ('mp3', 'wav')"`
	PredictionID string            `gorm:"uniqueIndex"`                                                                                                         // PredictionID guarda el ID de Replicate para cuando llegue el webhook
	Status       samplerenums.Status `gorm:"type:varchar(20);default:'processing';check:status IN ('starting', 'processing', 'succeeded', 'failed', 'canceled')"` // processing, succeeded, failed
	CreatedAt    time.Time         `gorm:"autoUpdateTime"`
}

// --- EDICIÓN Y EFECTOS ---
type SampleVersion struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SampleID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Effects       datatypes.JSON `gorm:"type:jsonb"`
	FinalAudioURL string         `gorm:"type:text"`
	CreatedAt     time.Time      `gorm:"autoUpdateTime"`
}

// --- SOCIAL ---
type SharedSample struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SampleID        *uuid.UUID `gorm:"type:uuid;index"`
	SampleVersionID *uuid.UUID `gorm:"type:uuid;index"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	Likes           int        `gorm:"default:0"`
	Downloads       int        `gorm:"default:0"`
	CreatedAt       time.Time  `gorm:"autoUpdateTime"`
}

// --- NOTIFICACIONES Y WEBHOOKS ---
type Notification struct {
	ID          uuid.UUID                   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID                   `gorm:"type:uuid;not null;index"`
	Type        notificationsenums.TypeNotification `gorm:"size:50;not null;check:type_notification IN ('replicate_error', 'replicate_success', 'payment', 'info')"` 
	Title       string                      `gorm:"size:255"`
	Message     string                      `gorm:"type:text"`
	Status      notificationsenums.NotificationStatus           `gorm:"type:varchar(20);default:'unread';check:status IN ('unread', 'read')"` 
	ReferenceID uuid.UUID                   `gorm:"type:uuid"` // ID del sample relacionado
	CreatedAt   time.Time                   `gorm:"autoUpdateTime"`
}

// --- PAGOS (PADDLE) ---
type PaymentLog struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ExternalRefID string         `gorm:"uniqueIndex"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null"`
	RawPayload    datatypes.JSON `gorm:"type:jsonb"`
	Status        string
	CreatedAt     time.Time `gorm:"autoUpdateTime"`
}
