package authusecases

import (
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	authentities "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/entities"
	authvalueobjects "github.com/PurpleSavage/monekai-server/modules/shared/auth/domain/valueobjects"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
)
type LoginUseCase struct{
	jwtService  authports.JwtPort
	repo authports.UserPersistencePort
	authService authports.AuthProviderPort
    registerUseCase *RegisterUseCase
    findUserByEmailUC  *FindUserByEmailUseCase
    updateSessionUC *UpdateSessionUseCase
}

func NewLoginUseCase(jwt authports.JwtPort,
	repo authports.UserPersistencePort,
	authService authports.AuthProviderPort,
    registerUC *RegisterUseCase, 
    findUserByEmailUC *FindUserByEmailUseCase,
    updateSessionUC *UpdateSessionUseCase) *LoginUseCase{

	return &LoginUseCase{
		jwtService:jwt,
		repo: repo,
		authService: authService,
        registerUseCase: registerUC,
        findUserByEmailUC:findUserByEmailUC,
        updateSessionUC:updateSessionUC,
	}
}
func (l *LoginUseCase) Execute(token string, userAgent string) (*authentities.SessionEntity, error) {
	userProvider, err := l.authService.VerifyAndExtract(token)
	if err != nil {
		return nil, err
	}
	
	// 2. Intentamos buscar al usuario en nuestra BD local
	user, err := l.findUserByEmailUC.Execute(userProvider.Email)

	// 3. MANEJO DE USUARIO NUEVO (Si no existe en nuestra BD)
	if globalerrors.IsNotFound(err) {
		return l.handleNewUser(userProvider, userAgent)
	}

	// 4. Si hubo otro tipo de error de base de datos (conexión, etc.), salimos
	if err != nil {
		return nil, err
	}

	// 5. USUARIO EXISTENTE: Ahora que estamos 100% seguros de que 'user' no es nil, creamos el DTO
	dto := authrequestsdtos.SessionRequestDto{
		Id:    user.Id,
		Email: user.Email,
	}

	// Generamos el Refresh Token usando el STRUCT completo
	refreshToken, err := l.jwtService.GenerateToken(dto, "72h")
	if err != nil {
		return nil, err
	}

	// Actualizamos el hash de la sesión persistente en la BD
	if err := l.updateSessionUC.Execute(refreshToken, user.Id); err != nil {
		return nil, err
	}

	return l.buildSessionResponse(user, refreshToken)
}

func (l *LoginUseCase) buildSessionResponse(user *authentities.UserEntity, refreshToken string) (*authentities.SessionEntity, error) {
	// 💡 CORRECCIÓN CRÍTICA: El Access Token también debe llevar el STRUCT, no el email plano
	dto := authrequestsdtos.SessionRequestDto{
		Id:    user.Id,
		Email: user.Email,
	}

	accessToken, err := l.jwtService.GenerateToken(dto, "1m")
	if err != nil {
		return nil, err 
	}

	return &authentities.SessionEntity{
		UserData:     *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (l *LoginUseCase) handleNewUser(providerUser *authports.ResponseProvider, userAgent string) (*authentities.SessionEntity, error) {
	// 1. Mandamos el flujo temporal al registro para persistir al usuario en la BD local
	newUserVO := authvalueobjects.CreateUserVO{
		ExternalID:   providerUser.AuthID,
		Email:        providerUser.Email,
		PhotoUrl:     providerUser.PhotoUrl,
		RefreshToken: "", // Lo dejamos vacío temporalmente porque requerimos primero el ID real de la base de datos
		UserAgent:    userAgent,
	}

	userCreated, err := l.registerUseCase.Execute(newUserVO)
	if err != nil {
		return nil, err
	}

	// 2. Ahora que la BD ya guardó al usuario y nos retornó su ID autogenerado, armamos su DTO de sesión
	dto := authrequestsdtos.SessionRequestDto{
		Id:    userCreated.Id,
		Email: userCreated.Email,
	}

	// 3. Generamos su Refresh Token estructurado real
	refreshToken, err := l.jwtService.GenerateToken(dto, "72h")
	if err != nil {
		return nil, err
	}

	// 4. Guardamos el token en la BD del nuevo usuario para dejar su sesión activa
	if err := l.updateSessionUC.Execute(refreshToken, userCreated.Id); err != nil {
		return nil, err
	}

	// 5. Construimos la respuesta final de sesión (el método internamente le generará su Access Token con el struct)
	return l.buildSessionResponse(userCreated, refreshToken)
}