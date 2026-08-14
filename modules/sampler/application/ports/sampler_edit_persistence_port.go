package samplerports

import (
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
)

type SamplerEditPort interface{
	SaveEditedSample(dto samplerequestsdto.SaveEditSampleWithURLDTO)(*samplerentities.EditedSampleEntity,error)
	GetEditedSampleByID(id string)(*samplerentities.EditedSampleEntity,error)
	UpdateURLEditedSample(url string)(bool,error)
	UpdateEffectsEditedSample(vo commonvalueobjects.EffectsVO)(bool,error)
}