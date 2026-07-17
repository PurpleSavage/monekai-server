package communityresponsesdtos

import "github.com/google/uuid"

type SharedSampleItemDTO struct {
	ID        uuid.UUID      `json:"id"`
	Likes     int            `json:"likes"`
	Downloads int            `json:"downloads"`
	CreatedAt string         `json:"createdAt"`
	Sample    SampleInfoDTO  `json:"sample"`
	SharedBy  SharedByDTO    `json:"sharedBy"`
}

