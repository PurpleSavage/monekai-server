package communityentities

import "github.com/google/uuid"
type SampleInfo struct {
	ID        uuid.UUID
	SampleName string
	InitialAudioURL string
	Prompt string
	Duration int
}
