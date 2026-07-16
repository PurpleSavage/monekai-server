package communityresponsesdtos

import "github.com/google/uuid"
type SampleInfoDTO struct {
	ID              uuid.UUID `json:"id"`
	SampleName      string    `json:"sampleName"`
	InitialAudioURL string    `json:"initialAudioUrl"`
	Prompt          string    `json:"prompt"`
	Duration        int       `json:"duration"`
}