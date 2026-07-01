package samplerraws

import (
	"time"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	"github.com/google/uuid"
)
type JoinResult struct {
	ID              uuid.UUID               `gorm:"column:id"`
    UserID          uuid.UUID               `gorm:"column:user_id"`
    SampleName      string                  `gorm:"column:sample_name"`
    Prompt          string                  `gorm:"column:prompt"`
    InitialAudioURL *string                 `gorm:"column:initial_audio_url"`
    Duration        int                     `gorm:"column:duration"`
    ModelVersion    samplerenums.ModelVersion `gorm:"column:model_version"`
    OutputFormat    samplerenums.OutputFormat `gorm:"column:output_format"`
    PredictionID    string                  `gorm:"column:prediction_id"`
    Status          samplerenums.Status       `gorm:"column:status"`
    CreatedAt       time.Time               `gorm:"column:created_at"`

    Email string `gorm:"column:email"`
}