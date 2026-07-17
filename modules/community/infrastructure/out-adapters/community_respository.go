package communityoutadapters

import (
	"context"

	models "github.com/PurpleSavage/monekai-server/configurations/persistence"
	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	communityinfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/mappers"
	communityraws "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/raws"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityRepository struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) communityports.CommunityPersistencePort{
	return &CommunityRepository{db: db}
}
//respositorio para lista los samples compartidos 
func (r *CommunityRepository) ListSharedSamples(
	ctx context.Context, 
	page int, 
	limit int,
) ([]communityentities.SharedSample, error){
 if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 
	}

	offset := (page - 1) * limit
	var dbSharedSamplesDetails []communityraws.SharedSampleDetailRaw
	err := r.db.WithContext(ctx).Model(&models.SharedSample{}).
			Select(`
				shared_samples.id,
				shared_samples.sample_id,
				shared_samples.sample_version_id,
				shared_samples.likes,
				shared_samples.downloads,
				shared_samples.created_at,
				
				samples.sample_name,
				samples.initial_audio_url,
				samples.prompt,
				samples.duration,
				
				users.email AS user_email,
				shared_samples.user_id AS user_id,
				users.name AS user_name
			`).
			Joins("INNER JOIN samples ON samples.id = shared_samples.sample_id").
			Joins("INNER JOIN users ON users.id = shared_samples.user_id").
			Order("shared_samples.created_at DESC").
			Limit(limit).
			Offset(offset).
			Scan(&dbSharedSamplesDetails).Error

		if err != nil {
			return nil, globalerrors.NewAppError(500, "Failed to list shared samples", err.Error(), err)
		}
	return communityinfrastructuremappers.MapToSharedSamplesDomain(dbSharedSamplesDetails), nil
}
//función para listar los samples comaprtidos pero por verssion --> un sample compartido con verssion es un sample con efectos aplicados
func (r *CommunityRepository) ListSharedSamplesVersion(
	ctx context.Context,
	page int,
	limit int,
) ([]communityentities.SharedSampleVersion, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	var dbSharedVersions []communityraws.SharedSampleVersionRaw

	err := r.db.WithContext(ctx).Model(&models.SharedSample{}).
		Select(`
			shared_samples.id,
			shared_samples.likes,
			shared_samples.downloads,
			shared_samples.created_at,
			shared_samples.sample_version_id,
			sample_versions.effects,
			sample_versions.final_audio_url,
			samples.sample_name,
			samples.prompt,
			users.email AS user_email,
			shared_samples.user_id AS user_id,
			users.name AS user_name
		`).
		Joins("INNER JOIN sample_versions ON sample_versions.id = shared_samples.sample_version_id").
		Joins("INNER JOIN samples ON samples.id = sample_versions.sample_id").
		Joins("INNER JOIN users ON users.id = shared_samples.user_id").
		Order("shared_samples.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&dbSharedVersions).Error

	if err != nil {
		return nil, globalerrors.NewAppError(500, "Failed to list shared sample versions", err.Error(), nil)
	}

	return communityinfrastructuremappers.MapToSharedSampleVersionsDomain(dbSharedVersions), nil
}

func (r *CommunityRepository) LikeToSharedSample(
	ctx context.Context,
	sampleID uuid.UUID,
) (*authvalueobjects.UUIDVO, error) {
	var sharedSampleModel models.SharedSample
	err := r.db.WithContext(ctx).Where("id = ?", sampleID).First(&sharedSampleModel).Error
	if err != nil {
		return nil, globalerrors.NewAppError(404, "Shared sample not found", err.Error(), nil)
	}
	sharedSampleModel.Likes++
	err = r.db.WithContext(ctx).Save(&sharedSampleModel).Error
	if err != nil {
		return nil, globalerrors.NewAppError(500, "Failed to like shared sample", err.Error(), nil)
	}
	return authvalueobjects.NewUUIDVO(sharedSampleModel.ID.String())
}
