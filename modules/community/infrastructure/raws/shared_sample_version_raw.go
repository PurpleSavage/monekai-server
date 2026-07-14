package communityraws

import (
	"time"

	communityvalueobjects "github.com/PurpleSavage/monekai-server/modules/community/domain/valueobjects"
	"github.com/google/uuid"
)

type SharedSampleVersionRaw struct {
	// Info base de SharedSample
	ID              uuid.UUID  `gorm:"column:id" json:"id"`
	Likes           int        `gorm:"column:likes" json:"likes"`
	Downloads       int        `gorm:"column:downloads" json:"downloads"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"createdAt"`

	// Info unida de SampleVersion (Se mapea automáticamente con sql.Scanner)
	SampleVersionID uuid.UUID `gorm:"column:sample_version_id" json:"sampleVersionId"`
	Effects         communityvalueobjects.EffectsVO `gorm:"column:effects" json:"effects"`
	FinalAudioURL   string     `gorm:"column:final_audio_url" json:"finalAudioUrl"`
	
	// Info unida de Sample (Vía alias o Join directo)
	SampleName      string     `gorm:"column:sample_name" json:"sampleName"`
	Prompt          string     `gorm:"column:prompt" json:"prompt"`

	// Info unida de User (Haciendo match con los alias AS del SELECT)
	UserEmail       string     `gorm:"column:user_email" json:"userEmail"`
	UserID          uuid.UUID  `gorm:"column:user_id" json:"userId"`
	UserName        string     `gorm:"column:user_name" json:"name"`
}