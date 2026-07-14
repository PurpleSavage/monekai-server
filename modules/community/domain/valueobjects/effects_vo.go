package communityvalueobjects

import commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"

// -------------------------
// RANGOS VÁLIDOS
// -------------------------
const (
	minReverb = 0
	maxReverb = 100

	minSlowPitch = -12
	maxSlowPitch = 0

	minSaturation = 0
	maxSaturation = 100

	minDelay = 0
	maxDelay = 30

	minLowPass = 12500
	maxLowPass = 20000

	minHighPass = 0
	maxHighPass = 40

	minGain = -3
	maxGain = 0
)

type EffectsVO struct {
	Reverb     int
	SlowPitch  int
	Saturation int
	Delay      int
	LowPass    int
	HighPass   int
	Gain       int
	Reverse    bool
}

// -------------------------
// CONSTRUCTOR CON VALIDACIÓN
// -------------------------
func CreateEffectsVO(
	Reverb int,
	SlowPitch int,
	Saturation int,
	Delay int,
	LowPass int,
	HighPass int,
	Gain int,
	Reverse bool,
) (*EffectsVO, error) {
	// -------------------------
	// Reverb VALIDATION
	// -------------------------
	if Reverb < minReverb || Reverb> maxReverb {
		return nil, commondomainerrors.NewValidationError(
			"reverb",
			"reverb must be between 0 and 100",
		)
	}
	// -------------------------
	// SLOW PITCH VALIDATION
	// -------------------------
	if SlowPitch < minSlowPitch || SlowPitch > maxSlowPitch {
		return nil, commondomainerrors.NewValidationError(
			"slowPitch",
			"slowPitch must be between -12 and 0",
		)
	}
	// -------------------------
	// SATURATION VALIDATION
	// -------------------------
	if Saturation < minSaturation || Saturation > maxSaturation {
		return nil, commondomainerrors.NewValidationError(
			"saturation",
			"saturation must be between 0 and 100",
		)
	}
	// -------------------------
	// DELAY VALIDATION
	// -------------------------
	if Delay < minDelay || Delay > maxDelay {
		return nil, commondomainerrors.NewValidationError(
			"delay",
			"delay must be between 0 and 30",
		)
	}
	// -------------------------
	// LOW PASS VALIDATION
	// -------------------------
	if LowPass < minLowPass || LowPass > maxLowPass {
		return nil, commondomainerrors.NewValidationError(
			"lowPass",
			"lowPass must be between 12500 and 20000",
		)
	}
	// -------------------------
	// HIGH PASS VALIDATION
	// -------------------------
	if HighPass < minHighPass || HighPass > maxHighPass {
		return nil, commondomainerrors.NewValidationError(
			"highPass",
			"highPass must be between 0 and 40",
		)
	}
	// -------------------------
	// GAIN VALIDATION
	// -------------------------
	if Gain < minGain || Gain > maxGain {
		return nil, commondomainerrors.NewValidationError(
			"gain",
			"gain must be between -3 and 0",
		)
	}
	// -------------------------
	// CREATE VO
	// -------------------------
	vo := &EffectsVO{
		Reverb:     Reverb,
		SlowPitch:  SlowPitch,
		Saturation: Saturation,
		Delay:      Delay,
		LowPass:    LowPass,
		HighPass:   HighPass,
		Gain:       Gain,
		Reverse:    Reverse,
	}
	return vo, nil
}