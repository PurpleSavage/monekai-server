package sampleroutadapters

import (
	"context"
	"encoding/json"
	"errors"

	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
	samplerinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/mappers"
	samplerraws "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/raws"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commonvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SamplerEditedRepository struct {
	db gorm.DB
}

func NewSamplerEditedRepository(
	db gorm.DB,
) samplerports.SamplerEditedPersistencePort {
	return &SamplerEditedRepository{
		db: db,
	}
}

func (s *SamplerEditedRepository) SaveEditedSample(
	dto samplerequestsdto.SaveEditSampleWithURLDTO,
) (*authvalueobjects.UUIDVO, error) {
	vo, err := samplervalueobjects.CreateSampleEditedVO(
		dto.SampleID,
		dto.Effects,
		dto.FinalAudioURL,
	)
	if err != nil {
		return nil, err
	}
	effects, err := vo.Effects.GetEffectsToJSON()
	if err != nil {
		return nil, err
	}
	samplelVersion := models.SampleVersion{
		SampleID:      vo.SampleID,
		Effects:       effects,
		FinalAudioURL: vo.FinalAudioURL,
	}
	if err := s.db.Create(&samplelVersion).Error; err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"The sample version record could not be saved in the database",
			err,
		)
	}
	createdID, err := authvalueobjects.NewUUIDVO(samplelVersion.ID.String())
	if err != nil {
		return nil, err
	}

	return createdID, nil
}
func (s *SamplerEditedRepository) GetEditedSampleByID(id string) (*samplerentities.EditedSampleEntity, error) {
	var result samplerraws.SamplerEditedJoinRaw
	err := s.db.
		Table("sample_versions").
		Select(`
			sample_versions.id,
		    sample_versions.effects,
		    sample_versions.final_audio_url,
		    sample_versions.created_at,
		    s.id AS sampleId,
		    s.sample_name,
		    s.prompt,
		    s.initial_audio_url,
		    s.duration,
		    s.output_format,
		    s.model_version,
		    s.status
			`).
		Joins("INNER JOIN samples s ON s.id = sample_versions.sample_id").
		Where("sample_versions.id = ?", id).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, globalerrors.NewAppError(
				404,
				"Not Found",
				"No edited sample record found with the provided ID",
				err,
			)
		}
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error retrieving edited sample",
			err,
		)
	}

	return samplerinfrastructuremappers.ToEditedSampleEntity(result), nil
}
func (s *SamplerEditedRepository) UpdateURLEditedSample(id string, url string) (bool, error) {
	result := s.db.Model(&models.SampleVersion{}).
		Where("id = ?", id).
		Update("final_audio_url", url)
	if result.Error != nil {
		return false, globalerrors.NewAppError(
			500,
			"Database Error",
			"The final audio URL could not be updated",
			result.Error,
		)
	}
	if result.RowsAffected == 0 {
		return false, globalerrors.NewAppError(
			404,
			"Not Found",
			"No sample version record found with the provided ID",
			nil,
		)
	}
	return true, nil
}

func (s *SamplerEditedRepository) UpdateEffectsEditedSample(id string, vo commonvalueobjects.EffectsVO) (bool, error) {
	var existing models.SampleVersion
	err := s.db.First(&existing, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, globalerrors.NewAppError(
				404,
				"Not Found",
				"No sample version record found with the provided ID",
				err,
			)
		}
		return false, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error retrieving sample version to update effects",
			err,
		)
	}

	merged, err := mergeEffects(existing.Effects, vo)
	if err != nil {
		return false, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error merging the stored effects with the updated ones",
			err,
		)
	}
	effectsJSON, err := merged.GetEffectsToJSON()
	if err != nil {
		return false, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error serializing the effects",
			err,
		)
	}

	result := s.db.Model(&models.SampleVersion{}).
		Where("id = ?", id).
		Update("effects", effectsJSON)
	if result.Error != nil {
		return false, globalerrors.NewAppError(
			500,
			"Database Error",
			"The effects could not be updated",
			result.Error,
		)
	}

	return true, nil
}

func mergeEffects(existingJSON datatypes.JSON, incoming commonvalueobjects.EffectsVO) (*commonvalueobjects.EffectsVO, error) {
	merged := commonvalueobjects.EffectsVO{}
	if len(existingJSON) > 0 && string(existingJSON) != "null" {
		if err := json.Unmarshal(existingJSON, &merged); err != nil {
			return nil, err
		}
	}

	if incoming.Reverb != nil {
		merged.Reverb = incoming.Reverb
	}
	if incoming.SlowPitch != nil {
		merged.SlowPitch = incoming.SlowPitch
	}
	if incoming.Saturation != nil {
		merged.Saturation = incoming.Saturation
	}
	if incoming.Delay != nil {
		merged.Delay = incoming.Delay
	}
	if incoming.LowPass != nil {
		merged.LowPass = incoming.LowPass
	}
	if incoming.HighPass != nil {
		merged.HighPass = incoming.HighPass
	}
	if incoming.Gain != nil {
		merged.Gain = incoming.Gain
	}
	if incoming.Reverse != nil {
		merged.Reverse = incoming.Reverse
	}

	return &merged, nil
}

func (s *SamplerEditedRepository) ListSamplesEditedByUsertID(ctx context.Context, userID string, page int, limit int) ([]*samplerentities.EditedSampleEntity, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	var result []samplerraws.SamplerEditedJoinRaw
	err := s.db.WithContext(ctx).
		Table("sample_versions").
		Select(`
			sample_versions.id,
		    sample_versions.effects,
		    sample_versions.final_audio_url,
		    sample_versions.created_at,
		    s.id AS sampleId,
		    s.sample_name,
		    s.prompt,
		    s.initial_audio_url,
		    s.duration,
		    s.output_format,
		    s.model_version,
		    s.status
			`).
		Joins("INNER JOIN samples s ON s.id = sample_versions.sample_id").
		Where("s.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	if err != nil {
		return nil, globalerrors.NewAppError(
			500,
			"Database Error",
			"Error listing your samples",
			err,
		)
	}

	domainSamples := make([]*samplerentities.EditedSampleEntity, len(result))
	for i, sample := range result {
		domainSamples[i] = samplerinfrastructuremappers.ToEditedSampleEntity(sample)
	}

	return domainSamples, nil
}
