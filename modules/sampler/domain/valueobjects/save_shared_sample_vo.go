package samplervalueobjects

import (
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	"github.com/google/uuid"
)

type SaveSharedSampleVO struct {
	SampleID  uuid.UUID
	UserID    uuid.UUID
	Likes     int
	Downloads int
}

func CreateSaveSharedSampleVO(
	sampleID string,
	userID string,
	likes *int,
	downloads *int,
) (*SaveSharedSampleVO, error) {
	userUUID, err := authvalueobjects.NewUUIDVO(userID)
	if err != nil {
		return nil, err
	}
	sampleUUID,err:= authvalueobjects.NewUUIDVO(sampleID)
	if err != nil {
		return nil, err
	}

	finalLikes := 0
	if likes != nil {
		finalLikes = *likes
	}

	finalDownloads := 0
	if downloads != nil {
		finalDownloads = *downloads
	}

	return &SaveSharedSampleVO{
		SampleID:  userUUID.Value(),
		UserID:    sampleUUID.Value(),
		Likes:     finalLikes,
		Downloads: finalDownloads,
	}, nil
}