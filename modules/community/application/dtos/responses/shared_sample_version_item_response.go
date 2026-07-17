package communityresponsesdtos

import "github.com/google/uuid"
type SharedSampleVersionItemDTO struct {
	ID            uuid.UUID             `json:"id"`
	Likes         int                   `json:"likes"`
	Downloads     int                   `json:"downloads"`
	CreatedAt     string                `json:"createdAt"`
	SampleVersion SampleVersionInfoDTO  `json:"sampleVersion"`
	SharedBy      SharedByDTO           `json:"sharedBy"` // Reutiliza el SharedByDTO del sample común
}