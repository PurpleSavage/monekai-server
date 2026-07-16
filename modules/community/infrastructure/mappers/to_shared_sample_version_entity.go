package communityinfrastructuremappers

import (
	"time"
	communityentities "github.com/PurpleSavage/monekai-server/modules/community/domain/entities"
	communityraws "github.com/PurpleSavage/monekai-server/modules/community/infrastructure/raws"
)

func MapToSharedSampleVersionsDomain(rawVersions []communityraws.SharedSampleVersionRaw) []communityentities.SharedSampleVersion {
	domainVersions := make([]communityentities.SharedSampleVersion, len(rawVersions))

	for i, raw := range rawVersions {
		

		domainVersions[i] = communityentities.SharedSampleVersion{
			ID:        raw.ID,
			Likes:     raw.Likes,
			Downloads: raw.Downloads,
			CreatedAt: raw.CreatedAt.Format(time.RFC3339),
			
			// Mapeo estructurado para el creador
			SharedBy: communityentities.SharedBy{
				ID:    raw.UserID,
				Name:  raw.UserName,
				Email: raw.UserEmail,
			},
			
			// Mapeo estructurado para la versión editada con su JSON
			SampleVersion: communityentities.SampleVersionInfo{
				ID:            raw.SampleVersionID,
				Effects:       raw.Effects, // Pasa limpio gracias a sql.Scanner
				FinalAudioURL: raw.FinalAudioURL,
				SampleName:    raw.SampleName,
				Prompt:        raw.Prompt,
			},
		}
	}

	return domainVersions
}