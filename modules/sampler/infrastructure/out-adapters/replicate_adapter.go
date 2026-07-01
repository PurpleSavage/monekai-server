package sampleroutadapters

import (
	"context"
	"log"

	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"github.com/replicate/replicate-go"
)


const SONG_MODEL string="meta/musicgen:671ac645ce5e552cc63a54a2bbff63fcf798043055d2dac5fc9e36a837eedcfb"

type ReplicateAdapterService struct{
	client *replicate.Client
}

func NewReplicateAdapterService() (samplerports.SongGeneratorPort, error) {
	client, err := replicate.NewClient(
		replicate.WithToken(config.Envs.ReplicateKey),
	)

	if err != nil {
		return nil, err
	}

	return &ReplicateAdapterService{
		client: client,
	}, nil
}
func (re *ReplicateAdapterService) GenerateSong(vo *samplervalueobjects.PayloadGenerateSongVO)(*samplerentities.SongGeneratorResponse, error){
	ctx := context.Background()
	input := replicate.PredictionInput{
		"prompt": vo.Prompt,
		"model_version":vo.ModelVersion,
		"duration":vo.Duration,
		"normalization_strategy":vo.NormalizationStrategy,
		"classifier_free_guidance":vo.ClassifierFreeGuidance,
		"output_format":vo.OutputFormat,
		"top_k":vo.TopK,
		"top_p":vo.TopP,
		"temperature":vo.Temperature,
	}
	log.Printf("input: %v", vo.WebhookUrl)
	webhook := replicate.Webhook{
		URL:    vo.WebhookUrl,
		Events: []replicate.WebhookEventType{"start", "completed"},
	}
	prediction, err := re.client.CreatePrediction(
		ctx,
		SONG_MODEL,
		input,
		&webhook,
		false,
	)
	if err != nil {
		return nil, globalerrors.NewAppError(
			502, // Bad Gateway (porque falló un servicio externo)
			"Replicate API Error",
			"Failed to trigger song generation via Replicate",
			err,
		)
	}
	response := &samplerentities.SongGeneratorResponse{
		GenerationId:     prediction.ID,
		StatusGeneration: samplerenums.Status(prediction.Status),
		UserId: vo.UserId,
	}

	return response, nil

}