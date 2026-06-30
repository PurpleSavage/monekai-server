package authinadapters

import (
	"encoding/json"
	"fmt"
	"time"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"github.com/golang-jwt/jwt/v5"
)


type JwtAdapterService struct{}

func NewJwtAdapterService() authports.JwtPort{
	return  &JwtAdapterService{}
}

func (a *JwtAdapterService) GenerateToken(data any, durationStr string) (string, error) {
	secret := config.Envs.SecretJwt
	duration, err := time.ParseDuration(durationStr)
	
	if err != nil {
		// 400 Bad Request: The provided duration string is malformed
		return "", globalerrors.NewAppError(400, "Invalid Duration", "The token duration format is incorrect", err)
	}

	claims := jwt.MapClaims{
		"user_data": data,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
    
	if err != nil {
		// 500 Internal Server Error: Cryptographic signing failed
		return "", globalerrors.NewAppError(500, "Signature Error", "Could not sign the security token", err)
	}

	return signedToken, nil
}

func (a *JwtAdapterService) VerifyToken(token string) ( any, error) {
	secret := config.Envs.SecretJwt

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return  nil, globalerrors.NewAppError(401, "Invalid Token", "Session has expired or token is corrupt", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return  nil, globalerrors.NewAppError(401, "Unauthorized", "Token content is not processable", nil)
	}

	userData, ok := claims["user_data"]
	if !ok {
		return  nil, globalerrors.NewAppError(401, "Unauthorized", "Missing user data in token", nil)
	}

	return  userData, nil
}

func MapClaimsToStruct[T any](rawData any) (T, error) {
	var target T
	bytes, err := json.Marshal(rawData)
	if err != nil {
		return target, err
	}
	err = json.Unmarshal(bytes, &target)
	return target, err
}