package communityvalueobjects

import "github.com/google/uuid"

type DownloadSharedSampleVO struct {
	SampleID uuid.UUID
	Downloads int
}

func NewDownloadSharedSampleVO(sampleID uuid.UUID, downloads int) *DownloadSharedSampleVO {
	return &DownloadSharedSampleVO{
		SampleID: sampleID,
		Downloads: downloads,
	}
}