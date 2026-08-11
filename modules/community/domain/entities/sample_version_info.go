package communityentities

import (
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	"github.com/google/uuid"
)


type SampleVersionInfo struct {
	ID        uuid.UUID
	Effects   commonvalueobjects.EffectsVO
	FinalAudioURL string
	SampleName string
	Prompt string
}