package models

import (
	"time"
	notificationsenums "github.com/PurpleSavage/monekai-server/modules/notifications/domain/enums"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

//
// -------------------- SESSIONS --------------------
//

type Session struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	User User `gorm:"foreignKey:UserID"`

	RefreshToken string `gorm:"type:varchar;uniqueIndex;not null"`

	UserAgent string `gorm:"size:255"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

//
// -------------------- USERS --------------------
//

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	ExternalID string `gorm:"type:varchar;uniqueIndex;not null"`

	Email string `gorm:"type:varchar;uniqueIndex;not null"`

	Credits int `gorm:"default:0"`

	PhotoURL *string `gorm:"type:text"`

	Sessions []Session `gorm:"foreignKey:UserID"`

	Notifications []Notification `gorm:"foreignKey:UserID"`

	Payments []Payment `gorm:"foreignKey:UserID"`

	Samples []Sample `gorm:"foreignKey:UserID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//
// -------------------- SAMPLES --------------------
//

type Sample struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	User User `gorm:"foreignKey:UserID"`

	SampleName string `gorm:"size:100;not null"`

	Prompt string `gorm:"type:text"`

	InitialAudioURL *string `gorm:"type:text"`

	Duration int

	ModelVersion samplerenums.ModelVersion `gorm:"type:varchar(30);not null;check:model_version IN ('stereo-melody-large','stereo-large','melody-large','large')"`

	OutputFormat samplerenums.OutputFormat `gorm:"type:varchar(10);not null;check:output_format IN ('mp3','wav')"`

	PredictionID string `gorm:"uniqueIndex"`

	Status samplerenums.Status `gorm:"type:varchar(20);default:'processing';check:status IN ('starting','processing','succeeded','failed','canceled')"`

	Versions []SampleVersion `gorm:"foreignKey:SampleID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//
// -------------------- SAMPLE VERSIONS --------------------
//

type SampleVersion struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	SampleID uuid.UUID `gorm:"type:uuid;not null;index"`

	Sample Sample `gorm:"foreignKey:SampleID"`

	Effects datatypes.JSON `gorm:"type:jsonb"`

	FinalAudioURL string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//
// -------------------- SHARED SAMPLES --------------------
//

type SharedSample struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	SampleID *uuid.UUID `gorm:"type:uuid;index"`

	Sample *Sample `gorm:"foreignKey:SampleID"`

	SampleVersionID *uuid.UUID `gorm:"type:uuid;index"`

	SampleVersion *SampleVersion `gorm:"foreignKey:SampleVersionID"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	User User `gorm:"foreignKey:UserID"`

	Likes int `gorm:"default:0"`

	Downloads int `gorm:"default:0"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//
// -------------------- NOTIFICATIONS --------------------
//

type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	User User `gorm:"foreignKey:UserID"`

	Type notificationsenums.TypeNotification `gorm:"type:varchar(50);not null;check:type_notification IN ('replicate_error','replicate_success','payment','info')"`

	Title string `gorm:"size:255"`

	Message string `gorm:"type:text"`

	Status notificationsenums.NotificationStatus `gorm:"type:varchar(20);default:'unread';check:status IN ('unread','read')"`

	ReferenceID uuid.UUID `gorm:"type:uuid"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

//
// -------------------- CREDIT PACKAGES --------------------
//

type CreditPackage struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// paddle, stripe, polar...
	Provider string `gorm:"type:varchar(30);not null;default:'paddle';index"`

	// ID del precio en el proveedor
	PriceID string `gorm:"type:varchar;uniqueIndex;not null"`

	Name string `gorm:"size:100;not null"`

	Credits int `gorm:"not null"`

	// $25.00 = 2500
	PriceCents int `gorm:"not null"`

	Currency string `gorm:"size:3;not null;default:'USD'"`

	Active bool `gorm:"not null;default:true"`

	Payments []Payment `gorm:"foreignKey:CreditPackageID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

//
// -------------------- PAYMENTS --------------------
//

type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	User User `gorm:"foreignKey:UserID"`

	CreditPackageID uuid.UUID `gorm:"type:uuid;not null;index"`

	CreditPackage CreditPackage `gorm:"foreignKey:CreditPackageID"`

	// paddle, stripe...
	Provider string `gorm:"type:varchar(30);not null;default:'paddle';index"`

	// txn_xxxxxx
	ProviderTransactionID string `gorm:"type:varchar;uniqueIndex;not null"`

	// price_xxxxxx
	PriceID string `gorm:"type:varchar;not null"`

	// Créditos entregados en ESTA compra
	CreditsPurchased int `gorm:"not null"`

	// $25 = 2500
	AmountCents int `gorm:"not null"`

	Currency string `gorm:"size:3;not null"`

	// pending, paid, failed, refunded...
	Status string `gorm:"type:varchar(30);not null"`

	// Respuesta completa del proveedor
	RawPayload datatypes.JSON `gorm:"type:jsonb"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}


/// falta una tabla de hisotrial de creditos y 