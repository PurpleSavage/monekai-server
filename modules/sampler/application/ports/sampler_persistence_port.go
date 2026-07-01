package samplerports

import (
	"context"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
)
type SamplerPersistencePort interface {
	SaveSampleWithoutAudioUrl(vo samplervalueobjects.SaveSampleVO) error
	UpdateAudioUrl(url string, predictionId string) (*samplervalueobjects.ResponseUpdateurlVo, error)
	ListSamples(ctx context.Context,userId string,page int,limit int) ([]*samplerentities.SampleEntity,error)
	CountTotalSamples(ctx context.Context,userId string)(int, error)
	SaveSharedSample(vo samplervalueobjects.SaveSharedSampleVO)(*samplerentities.ShareSampleEntity,error)
	FindSampleByPredictionId(id string) (*samplerentities.SampleEntity,string, error)
}
