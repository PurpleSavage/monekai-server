package samplerentities

import (
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	"github.com/google/uuid"
)


type EditedSampleEntity struct{
	ID uuid.UUID
	Sample SampleEntity
	Effects commonvalueobjects.EffectsVO
	FinalAudioURL string
	CreatedAt string
}