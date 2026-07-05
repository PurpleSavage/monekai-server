package samplerinfrastructuremappers

import (
	"time"
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
)
func ToSharedSampleEntity(dbSharedSample models.SharedSample) *samplerentities.ShareSampleEntity{
	return &samplerentities.ShareSampleEntity{
		Id:              dbSharedSample.ID,
		SampleId:        dbSharedSample.SampleID,
		SampleVersionId: dbSharedSample.SampleVersionID,
		UserId:          dbSharedSample.UserID.String(),
		Likes:           dbSharedSample.Likes,
		Downloads:       dbSharedSample.Downloads,
		CreatedAt:       dbSharedSample.CreatedAt.Format(time.RFC3339),
	}
}
