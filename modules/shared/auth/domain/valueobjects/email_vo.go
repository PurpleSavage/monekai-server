package authvalueobjects

import (
	"regexp"
	"strings"

	commondomainerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/errors"
)

var (
	// Misma regla para el usuario: sin puntos al inicio/fin ni consecutivos
	userPattern = regexp.MustCompile(`^[a-zA-Z0-9_+]+(\.[a-zA-Z0-9_+]+)*$`)

	// Modificado para el dominio completo:
	// Permite subdominios separados por puntos, etiquetados con guiones sin duplicar, 
	// y exige un TLD de 2 o más letras al final.
	domainPattern = regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*(\.[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)*\.[a-zA-Z]{2,}$`)
)

type EmailVO struct {
	value string
}

func CreateEmailVO(value string) (*EmailVO, error) {
	// 1. Limpieza de bordes para evaluar el contenido real
	valueValidatedEmpty := strings.TrimSpace(value)

	// Regla: No puede estar vacío después del trim
	if len(valueValidatedEmpty) == 0 {
		return nil, commondomainerrors.NewValidationError(
			"Email",
			"The email format is invalid.",
		)
	}

	// 2. Validación caracter por caracter: Ausencia de espacios y presencia única de @
	counterAts := 0
	for _, runeValue := range valueValidatedEmpty {
		// Regla: Ausencia de espacios
		if runeValue == ' ' {
			return nil, commondomainerrors.NewValidationError(
				"Email",
				"The email format is invalid.",
			)
		}

		if runeValue == '@' {
			counterAts++
		}
	}

	// Regla: Presencia ÚNICA del símbolo @ (debe ser exactamente 1)
	if counterAts != 1 {
		return nil, commondomainerrors.NewValidationError(
			"Email",
			"The email format is invalid.",
		)
	}

	// 3. División en Fragmentos (Usuario y Dominio)
	emailFragments := strings.Split(valueValidatedEmpty, "@")
	userPart := emailFragments[0]
	domainPart := emailFragments[1]

	// 4. Validación de la Parte Local (Usuario)
	// Regla: Letras, números, _, +, . (sin puntos al inicio/fin ni puntos consecutivos)
	if !userPattern.MatchString(userPart) {
		return nil, commondomainerrors.NewValidationError(
			"Email",
			"The email format is invalid.",
		)
	}

	// 5. Validación del Dominio y TLD
	// Regla: Letras, números, guiones, puntos (sin guiones al inicio/fin)
	// Regla: Punto obligatorio y TLD de 2 o más letras al final (.com, .pe, etc.)
	if !domainPattern.MatchString(domainPart) {
		return nil, commondomainerrors.NewValidationError(
			"Email",
			"The email format is invalid.",
		)
	}

	// Si todas las validaciones pasan, retornamos el Value Object válido
	return &EmailVO{
		value: valueValidatedEmpty,
	}, nil
}

// Getter para exponer el valor de forma inmutable
func (e *EmailVO) Value() string {
	return e.value
}