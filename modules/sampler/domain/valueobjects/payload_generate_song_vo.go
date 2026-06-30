package samplervalueobjects

import (
	"strings"

	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)
type PayloadGenerateSongVO struct {
	UserId	string
	// Internal webhook used by Replicate to notify
	// the backend when the generation finishes.
	WebhookUrl string

	// Music generation model version.
	ModelVersion samplerenums.ModelVersion

	// User prompt used to generate the music.
	Prompt string

	// Duration in seconds.
	// Current MVP target: max 30 seconds.
	Duration int

	// Audio output normalization strategy.
	// Example values:
	// - peak
	// - loudness
	//
	// For now this will usually be hardcoded internally,
	// but it remains exposed in case advanced settings
	// are added later.
	NormalizationStrategy samplerenums.NormalizationStrategy

	// Controls how strongly the prompt influences
	// the generation process.
	//
	// Higher values:
	// - more prompt adherence
	// - less randomness
	//
	// Lower values:
	// - more creativity
	// - more variation
	ClassifierFreeGuidance int

	// Final audio format returned by the model.
	// Example:
	// - mp3
	// - wav
	OutputFormat samplerenums.OutputFormat



	// Internal AI parameters


	TopK        int
	TopP        int
	Temperature float32
}
func CreatePayloadGenerateSongVO(
	webhookUrl string,
	modelVersion samplerenums.ModelVersion,
	prompt string,
	duration int,
	outputFormat samplerenums.OutputFormat,
	userId string,
) (*PayloadGenerateSongVO, error) {



	// -------------------------
	// PROMPT VALIDATION
	// -------------------------

	prompt = strings.TrimSpace(prompt)

	if prompt == "" {

		return nil, commondomainerrors.NewValidationError(
			"prompt",
			"prompt is required",
		)
	}

	if len(prompt) > 500 {

		return nil, commondomainerrors.NewValidationError(
			"prompt",
			"prompt exceeds maximum length",
		)
	}

	// -------------------------
	// DURATION VALIDATION
	// -------------------------

	if duration <= 0 {

		return nil, commondomainerrors.NewValidationError(
			"duration",
			"duration must be greater than 0",
		)
	}

	if duration > 30 {

		return nil, commondomainerrors.NewValidationError(
			"duration",
			"maximum duration is 30 seconds",
		)
	}

	// -------------------------
	// MODEL VALIDATION
	// -------------------------

	switch modelVersion {

	case samplerenums.StereoMelodyLarge,
		samplerenums.StereoLarge,
		samplerenums.MelodyLarge,
		samplerenums.Large:

	default:

		return nil, commondomainerrors.NewValidationError(
			"modelVersion",
			"unsupported model version",
		)
	}



	

	// -------------------------
	// OUTPUT FORMAT VALIDATION
	// -------------------------

	switch outputFormat {

	case samplerenums.Mp3,
		samplerenums.Wav:

	default:

		return nil, commondomainerrors.NewValidationError(
			"outputFormat",
			"unsupported output format",
		)
	}

	// -------------------------
	// CREATE VO
	// -------------------------

	vo := &PayloadGenerateSongVO{
		UserId: userId,


		WebhookUrl: webhookUrl,

		ModelVersion: modelVersion,

		Prompt: prompt,

		Duration: duration,

		NormalizationStrategy:samplerenums.Peak,

		ClassifierFreeGuidance: 3,

		OutputFormat: outputFormat,

		TopK: 250,
		TopP: 0,
		Temperature:1 ,
	}

	return vo, nil
}
/*
	FUTURE FEATURES (NOT INCLUDED IN MVP)

	The following parameters are intentionally excluded
	from the current implementation to keep the MVP simple,
	fast, and focused on prompt-only generation.

	Planned future capabilities:

	- continuation
		Allows continuing an existing audio generation.

	- input_audio
		Allows using uploaded audio or a URL as reference.

	- continuation_start
		Defines where continuation starts.

	- continuation_end
		Defines where continuation ends.

	- multi_band_diffusion
		Experimental quality enhancement option.


	ADVANCED MODEL PARAMETERS (CURRENTLY HARDCODED)

	These parameters are AI sampling controls and are
	not exposed to end users for now because they add
	unnecessary complexity to the MVP UX.

	Default values used internally:

	- TopK = 250
		Controls how many candidate tokens are considered.

	- TopP = 0
		Enables nucleus sampling when > 0.
		Value 0 means TopK sampling is used instead.

	- Temperature = 1
		Controls randomness/creativity of generation.
*/