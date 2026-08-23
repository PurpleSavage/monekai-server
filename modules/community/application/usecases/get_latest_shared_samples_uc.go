package communityusecases

import (
	"context"

	communityports "github.com/PurpleSavage/monekai-server/modules/community/application/ports"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)

type GetLatestSharedSamplesUC struct{
	repo communityports.CommunityPersistencePort
}

func (g *GetLatestSharedSamplesUC)Execute(
	ctx context.Context,
	limit int,
)(*commonentities.PaginatedResult[communityentities.SharedSample], error){
	response,err := g.repo.GetLatestSharedSamples(ctx,limit)
	if err != nil {
		return  nil, err
	}
	
}
