package samplervalueobjects

import (
	"strings"
	samplerenums "github.com/PurpleSavage/monekai-server/modules/sampler/domain/enums"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
	"github.com/google/uuid"
)
type SaveSampleVO struct {
	UserID       uuid.UUID
	SampleName   string
	Prompt       string
	Duration     int
	ModelVersion samplerenums.ModelVersion
	OutputFormat samplerenums.OutputFormat // <-- Agregado al VO
	PredictionID string
	Status       samplerenums.Status 
}

func CreateSaveSampleVO(
	userId string,
	sampleName string,
	prompt string,
	duration int,
	modelVersion string,
	outputFormat string, // <-- Agregado como argumento string
	predictionId string,
	status string,
) (*SaveSampleVO, error) {

	// -------------------------
	// USER ID VALIDATION
	// -------------------------
	userUUID, err := authvalueobjects.NewUUIDVO(userId)
	if err != nil {
		return nil, err
	}

	// -------------------------
	// SAMPLE NAME VALIDATION
	// -------------------------
	sampleName = strings.TrimSpace(sampleName)
	if sampleName == "" {
		return nil, commondomainerrors.NewValidationError(
			"sampleName",
			"sample name is required",
		)
	}
	if len(sampleName) > 100 {
		return nil, commondomainerrors.NewValidationError(
			"sampleName",
			"sample name exceeds maximum length of 100 characters",
		)
	}

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
	// MODEL VERSION VALIDATION & CASTING
	// -------------------------
	var typedModelVersion samplerenums.ModelVersion
	switch samplerenums.ModelVersion(modelVersion) {
	case samplerenums.StereoMelodyLarge,
	samplerenums.StereoLarge,
	samplerenums.MelodyLarge,
	samplerenums.Large:
	typedModelVersion = samplerenums.ModelVersion(modelVersion)
	default:
		return nil, commondomainerrors.NewValidationError(
			"modelVersion",
			"unsupported model version",
		)
	}

	// -------------------------
	// OUTPUT FORMAT VALIDATION & CASTING
	// -------------------------
	var typedOutputFormat samplerenums.OutputFormat
	switch samplerenums.OutputFormat(outputFormat) {
	case samplerenums.Mp3,
	samplerenums.Wav:
	typedOutputFormat = samplerenums.OutputFormat(outputFormat)
	default:
		return nil, commondomainerrors.NewValidationError(
			"outputFormat",
			"unsupported output format",
		)
	}

	// -------------------------
	// PREDICTION ID VALIDATION
	// -------------------------
	predictionId = strings.TrimSpace(predictionId)
	if predictionId == "" {
		return nil, commondomainerrors.NewValidationError(
			"predictionId",
			"prediction id is required",
		)
	}

	// -------------------------
	// STATUS VALIDATION & CASTING
	// -------------------------
	var typedStatus samplerenums.Status
	switch samplerenums.Status(status) {
	case samplerenums.Starting,
	samplerenums.Processing,
	samplerenums.Succeeded,
	samplerenums.Failed,
	samplerenums.Canceled:
	typedStatus = samplerenums.Status(status)
	default:
		return nil, commondomainerrors.NewValidationError(
			"status",
			"unsupported status value",
		)
	}

	// -------------------------
	// CREATE VO
	// -------------------------
	vo := &SaveSampleVO{
		UserID:       userUUID.Value(),
		SampleName:   sampleName,
		Prompt:       prompt,
		Duration:     duration,
		ModelVersion: typedModelVersion,
		OutputFormat: typedOutputFormat, // <-- Asignado al struct final
		PredictionID: predictionId,
		Status:       typedStatus,
	}

	return vo, nil
}