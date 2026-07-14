package communityentities

import "github.com/google/uuid"

type SharedSample struct {
	ID        uuid.UUID
	Sample  SampleInfo
	Likes     int
	Downloads int
	CreatedAt string
	SharedBy SharedBy
}
