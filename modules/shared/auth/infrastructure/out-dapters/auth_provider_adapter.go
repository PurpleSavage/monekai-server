package authoutadapters

import (
	"context"
	"fmt"

	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
	"google.golang.org/api/idtoken"
)

type AuthProviderAdapter struct {
	clientID string
}

// NewAuthProviderAdapter inicializa el adaptador consumiendo tu variable global Envs
func NewAuthProviderAdapter() authports.AuthProviderPort {
	// Consumimos el Client ID de Google que ya cargas en tu LoadEnvs()
	googleClientID := config.Envs.GoogleClientId
	
	return &AuthProviderAdapter{
		clientID: googleClientID,
	}
}

func (a *AuthProviderAdapter) VerifyAndExtract(token string) (*authports.ResponseProvider, error) {
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, token, a.clientID)
	if err != nil {
		return nil, fmt.Errorf("token de Google inválido o expirado: %w", err)
	}
	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	photo, _ := payload.Claims["picture"].(string)
	
	phone, _ := payload.Claims["phone_number"].(string)
	response := &authports.ResponseProvider{
		AuthID: googleID, 
		Email:   email,
	}

	if photo != "" {
		response.PhotoUrl = &photo
	}
	if phone != "" {
		response.PhoneNumber = &phone
	}

	return response, nil
}
