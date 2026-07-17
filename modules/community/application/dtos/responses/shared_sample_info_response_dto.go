package communityresponsesdtos

import (
	communityvalueobjects "github.com/PurpleSavage/monekai-server/modules/community/domain/valueobjects"
	"github.com/google/uuid"
)
type SampleVersionInfoDTO struct {
	ID            uuid.UUID                       `json:"id"`
	Effects       communityvalueobjects.EffectsVO `json:"effects"` // Se serializa a JSON de forma automática
	FinalAudioURL string                          `json:"finalAudioUrl"`
	SampleName    string                          `json:"sampleName"`
	Prompt        string                          `json:"prompt"`
}