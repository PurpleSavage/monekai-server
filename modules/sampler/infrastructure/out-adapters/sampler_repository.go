package sampleroutadapters

import (
	"context"
	"errors"

	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
	samplerinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/mappers"
	samplerraws "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/raws"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/gorm"
)
 type SamplerRepository struct {
	db *gorm.DB
}

func NewSamplerrepository(db *gorm.DB) samplerports.SamplerPersistencePort{
	return  &SamplerRepository{
		db: db,
	}
}

func (a *SamplerRepository) SaveSampleWithoutAudioUrl(vo samplervalueobjects.SaveSampleVO) error {
	newSample := models.Sample{
		UserID:          vo.UserID,
		SampleName:      vo.SampleName,
		Prompt:          vo.Prompt,
		InitialAudioURL: nil,
		Duration:        vo.Duration,
		ModelVersion:    vo.ModelVersion,
		PredictionID:    vo.PredictionID,
		Status:          vo.Status,
		OutputFormat:    vo.OutputFormat,
	}

	if err := a.db.Create(&newSample).Error; err != nil {
		return globalerrors.NewAppError(
			500,
			"Database Error",
			"The sample record could not be saved in the database",
			err,
		)
	}

	return nil
}

func (a *SamplerRepository) UpdateAudioUrl(url string, predictionId string) (*samplervalueobjects.ResponseUpdateurlVo, error) {

	err := a.db.Model(&models.Sample{}).
			Where("prediction_id = ?", predictionId).
			Updates(map[string]any{
				"initial_audio_url": url,
				"status":            samplerenums.Succeeded,
			}).Error

		if err != nil {
			return nil, globalerrors.NewAppError(
				500,
				"Database Error",
				"The sample audio URL could not be updated",
				err,
			)
		}

		// 2. Recuperar el sample actualizado + datos del usuario
		var result samplerraws.JoinResult

		err = a.db.Table("samples").
			Select(`
				samples.id as id,
				samples.sample_name,
				samples.prompt,
				samples.initial_audio_url,
				samples.duration,
				samples.model_version,
				samples.output_format,
				samples.prediction_id,
				samples.status,
				samples.created_at,
				users.id as user_id,
				users.email
			`).
			Joins("INNER JOIN users ON users.id = samples.user_id").
			Where("samples.prediction_id = ?", predictionId).
			First(&result).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, globalerrors.NewAppError(
					404,
					"Not Found",
					"No sample record found with the provided prediction ID",
					err,
				)
			}

			return nil, globalerrors.NewAppError(
				500,
				"Database Error",
				"Error loading updated sample data",
				err,
			)
		}

	

		return samplerinfrastructuremappers.ToResponseUpdateUrlVo(result), nil
}
func(a *SamplerRepository) ListSamples(ctx context.Context,userID string, page int, limit int) ([]*samplerentities.SampleEntity,error){
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 
	}

	offset := (page - 1) * limit
	var dbSamples []models.Sample

	err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("initial_audio_url IS NOT NULL").
		Limit(limit).  
		Offset(offset).
		Find(&dbSamples).
		Error
	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error listing your samples",
			err,
		)
	}
	domainSaples := make([]*samplerentities.SampleEntity, len(dbSamples))
	for i,dbSample:= range (dbSamples){
		domainSaples[i]=samplerinfrastructuremappers.ToSampleEntity(dbSample)
	}
	return  domainSaples, nil
}

func(a *SamplerRepository) CountTotalSamples(ctx context.Context,userID string)(int, error){
	var count int64
	err := a.db.WithContext(ctx).
		Model(&models.Sample{}). 
		Where("user_id = ?", userID).
		Count(&count).Error

	if err !=nil{
		return 0, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error retrieving total sample count from database",
			err,
		)
	}
	return int(count), nil
}

func (a*SamplerRepository) SaveSharedSample(vo samplervalueobjects.SaveSharedSampleVO)(*samplerentities.ShareSampleEntity,error){
	shareSampleModel:= models.SharedSample{
		SampleID:&vo.SampleID,
		UserID:vo.UserID,
		Likes:vo.Likes,
		Downloads:vo.Downloads,
	}
	if err := a.db.Create(&shareSampleModel).Error; err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"The sample record could not be saved in the database",
			err,
		)
	}

	sharedSampleResponse := samplerinfrastructuremappers.ToSharedSampleEntity(
		shareSampleModel,
	)

	return sharedSampleResponse, nil
}

func (a *SamplerRepository) FindSampleByPredictionId(id string) (*samplerentities.SampleEntity,string, error) {
	var dbSample models.Sample
	err := a.db.First(&dbSample, "prediction_id = ?", id).Error
	if err != nil {
		return nil,"", globalerrors.NewAppError(
			500,
			"Database Error",
			"Error retrieving sample by prediction id",
			err,
		)
	}
	sample:= samplerinfrastructuremappers.ToSampleEntity(dbSample)

	return sample,dbSample.UserID.String(), nil
}
