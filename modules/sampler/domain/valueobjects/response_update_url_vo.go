package samplervalueobjects

import (
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	"github.com/google/uuid"
)

type PersonalData struct {
	Email string 
	ID uuid.UUID
}

type ResponseUpdateurlVo struct {
	UserData PersonalData
	Sample samplerentities.SampleEntity
}