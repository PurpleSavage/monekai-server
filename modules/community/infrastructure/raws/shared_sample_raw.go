package communityraws

import (
	"time"

	"github.com/google/uuid"
)
type SharedSampleDetailRaw struct {
	// Info base de SharedSample
	ID              uuid.UUID  `gorm:"column:id" json:"id"`
	Likes           int        `gorm:"column:likes" json:"likes"`
	Downloads       int        `gorm:"column:downloads" json:"downloads"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"createdAt"`

	// Info unida de Sample (Haciendo match exacto con las columnas de la tabla samples)
	SampleID        uuid.UUID  `gorm:"column:sample_id" json:"sampleId"`
	SampleName      string     `gorm:"column:sample_name" json:"sampleName"`
	InitialAudioURL string     `gorm:"column:initial_audio_url" json:"initialAudioUrl"`
	Prompt          string     `gorm:"column:prompt" json:"prompt"`
	Duration        int        `gorm:"column:duration" json:"duration"`

	// Info unida de User (Haciendo match con los alias AS del SELECT)
	UserEmail       string     `gorm:"column:user_email" json:"userEmail"`
	UserID          uuid.UUID  `gorm:"column:user_id" json:"userId"`
	UserName        string     `gorm:"column:user_name" json:"name"`
}