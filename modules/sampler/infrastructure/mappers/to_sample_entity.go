package samplerinfrastructuremappers

import (
	"time"
	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
)
func ToSampleEntity(dbSample models.Sample) *samplerentities.SampleEntity {
	return &samplerentities.SampleEntity{
		Id:           dbSample.ID,
		SampleName:   dbSample.SampleName,
		Prompt:       dbSample.Prompt,
		AudioUrl:     dbSample.InitialAudioURL,
		Duration:     dbSample.Duration,
		OutputFormat: dbSample.OutputFormat,
		ModelVersion: dbSample.ModelVersion,
		Status:       dbSample.Status,
		CreatedAt:    dbSample.CreatedAt.Format(time.RFC3339),
	}
}