package authusecases

import (
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
)

type RefreshTokenUseCase struct{
	jwtService  authports.JwtPort
}

func NewRefreshTokenUseCase(
	jwtService authports.JwtPort, 
) *RefreshTokenUseCase{
	return &RefreshTokenUseCase{
		jwtService: jwtService,
	}
}


func (r *RefreshTokenUseCase) Execute(dto authrequestsdtos.SessionRequestDto) (authvalueobjects.TokenVO, error) {
	
	// Generamos el nuevo Access Token de corta duración (1 minuto) pasándole el struct completo
	token, err := r.jwtService.GenerateToken(dto, "1m")
	if err != nil {
		return authvalueobjects.TokenVO{}, err
	}
	
	validatedtoken, err := authvalueobjects.NewTokenVO(token)
	if err != nil {
		return authvalueobjects.TokenVO{}, err
	}
	
	return validatedtoken, nil
}
