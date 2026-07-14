package communityentities

import (
	communityvalueobjects "github.com/PurpleSavage/monekai-server/modules/community/domain/valueobjects"
	"github.com/google/uuid"
)


type SampleVersionInfo struct {
	ID        uuid.UUID
	Effects   communityvalueobjects.EffectsVO
	FinalAudioURL string
	SampleName string
	Prompt string
}