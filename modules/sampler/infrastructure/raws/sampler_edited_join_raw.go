package samplerraws

import (
	"time"

	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SamplerEditedJoinRaw struct {
	ID              uuid.UUID                 `gorm:"column:id"`
	Effects         datatypes.JSON            `gorm:"column:effects"`
	FinalAudioURL   string                    `gorm:"column:final_audio_url"`
	CreatedAt       time.Time                 `gorm:"column:created_at"`
	SampleID        uuid.UUID                 `gorm:"column:sampleId"`
	SampleName      string                    `gorm:"column:sample_name"`
	Prompt          string                    `gorm:"column:prompt"`
	InitialAudioURL string                    `gorm:"column:initial_audio_url"`
	Duration        int                       `gorm:"column:duration"`
	OutputFormat    samplerenums.OutputFormat `gorm:"column:output_format"`
	ModelVersion    samplerenums.ModelVersion `gorm:"column:model_version"`
	Status          samplerenums.Status       `gorm:"column:status"`
}
