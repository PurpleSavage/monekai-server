package samplerinfrastructuremappers

import (
	"time"
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	"github.com/google/uuid"
)
func ToSharedSampleEntity(dbSharedSample models.SharedSample) *samplerentities.ShareSampleEntity{
	var sampleID *uuid.UUID
	if dbSharedSample.SampleID != nil {
		value := dbSharedSample.SampleID
		sampleID = value
	}

	var sampleVersionID *string
	if dbSharedSample.SampleVersionID != nil {
		value := dbSharedSample.SampleVersionID.String()
		sampleVersionID = &value
	}

	return &samplerentities.ShareSampleEntity{
		Id:              dbSharedSample.ID,
		SampleId:        sampleID,
		SampleVersionId: sampleVersionID,
		UserId:          dbSharedSample.UserID.String(),
		Likes:           dbSharedSample.Likes,
		Downloads:       dbSharedSample.Downloads,
		CreatedAt:       dbSharedSample.CreatedAt.Format(time.RFC3339),
	}
}
