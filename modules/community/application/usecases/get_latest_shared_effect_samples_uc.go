package communityusecases

import (
	"context"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)

type GetLatestSharedEffectSamplesUC struct {
	repo  communityports.CommunityPersistencePort
	cache commonports.LocalCachePort[[]communityentities.SharedSampleVersion]
}

func NewGetLatestSharedEffectSamplesUC(
	repo communityports.CommunityPersistencePort,
	cache commonports.LocalCachePort[[]communityentities.SharedSampleVersion],
) *GetLatestSharedEffectSamplesUC {
	return &GetLatestSharedEffectSamplesUC{
		repo:  repo,
		cache: cache,
	}
}

func (g *GetLatestSharedEffectSamplesUC) Execute(
	ctx context.Context,
	limit int,
) (*commonentities.PaginatedResult[communityentities.SharedSampleVersion], error) {
	responseCached, _ := g.cache.Get("latest:shared-effect-samples")
	if responseCached != nil {
		return &commonentities.PaginatedResult[communityentities.SharedSampleVersion]{
			Total: len(*responseCached),
			Limit: limit,
			Page:  1,
			Data:  *responseCached,
		}, nil
	}
	response, err := g.repo.GetLatestSharedEffectSamples(ctx, limit)
	if err != nil {
		return nil, err
	}
	newResponseCached, err := g.cache.Set(response, "latest:shared-effect-samples", 3000)
	if err != nil {
		return &commonentities.PaginatedResult[communityentities.SharedSampleVersion]{
			Total: len(response),
			Limit: limit,
			Page:  1,
			Data:  response,
		}, nil
	}
	return &commonentities.PaginatedResult[communityentities.SharedSampleVersion]{
		Total: len(*newResponseCached),
		Limit: limit,
		Page:  1,
		Data:  *newResponseCached,
	}, nil
}
