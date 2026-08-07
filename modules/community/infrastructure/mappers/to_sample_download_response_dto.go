package communityinfrastructuremappers

import (
	communityresponsesdtos "github.com/PurpleSavage/monekai-server/modules/community/application/dtos/responses"
	communityvalueobjects "github.com/PurpleSavage/monekai-server/modules/community/domain/valueobjects"
)


func ToSampleDownloadResponseDTO(
	vo *communityvalueobjects.DownloadSharedSampleVO,
)communityresponsesdtos.DownloadSampleResponseDTO{
	return communityresponsesdtos.DownloadSampleResponseDTO{
		SampleId:vo.SampleID.String(),
	}
}