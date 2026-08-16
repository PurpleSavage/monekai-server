package samplerports

import (
	"context"

	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
)

type SamplerEditedPersistencePort interface {
	SaveEditedSample(dto samplerequestsdto.SaveEditSampleWithURLDTO) (*authvalueobjects.UUIDVO, error)
	GetEditedSampleByID(id string) (*samplerentities.EditedSampleEntity, error)
	UpdateURLEditedSample(id string, url string) (bool, error)
	UpdateEffectsEditedSample(id string, vo commonvalueobjects.EffectsVO) (bool, error)
	ListSamplesEditedByUsertID(ctx context.Context, userID string, page int, limit int) ([]*samplerentities.EditedSampleEntity, error)
}
