package communityinfrastructuremappers

import (
	"time"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	communityraws "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/raws"
)


func MapToSharedSamplesDomain(rawDetails []communityraws.SharedSampleDetailRaw) []communityentities.SharedSample {
	domainSamples := make([]communityentities.SharedSample, len(rawDetails))

	for i, raw := range rawDetails {
	
		domainSamples[i] = communityentities.SharedSample{
			ID:        raw.ID,
			Likes:     raw.Likes,
			Downloads: raw.Downloads,
			// Convertimos time.Time a string usando un formato estándar (RFC3339)
			CreatedAt: raw.CreatedAt.Format(time.RFC3339),
			Sample: communityentities.SampleInfo{
				ID:              raw.SampleID,
				SampleName:      raw.SampleName,
				InitialAudioURL: raw.InitialAudioURL,
				Prompt:          raw.Prompt,
				Duration:        raw.Duration,
			},
			SharedBy: communityentities.SharedBy{
				ID:    raw.UserID,
				Name:  raw.UserName,
				Email: raw.UserEmail,
			},
		}
	}

	return domainSamples
}