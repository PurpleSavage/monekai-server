package communityentities

import "github.com/google/uuid"

type SharedSampleVersion struct {
	ID        uuid.UUID
	SharedBy  SharedBy
	SampleVersion SampleVersionInfo
	Likes     int
	Downloads int
	CreatedAt string
}