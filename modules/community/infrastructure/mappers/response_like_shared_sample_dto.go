package communityinfrastructuremappers

import (
	communityresponsesdtos "github.com/PurpleSavage/monekai-server/modules/community/application/dtos/responses"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
)

func ResponseLikeSharedSampleDTOMapper(user *authvalueobjects.UUIDVO) communityresponsesdtos.LikeSharedSampleResponseDTO {
	return communityresponsesdtos.LikeSharedSampleResponseDTO{
		Message:         "Sample liked successfully",
		SampleIDModify:  user.String(),
	}
}