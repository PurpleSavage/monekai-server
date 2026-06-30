package auth

import (
	authusecases "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/usecases"
	authcontrollers "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/controllers"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	authoutadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/out-dapters"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)
func AuthBootstrap(db *gorm.DB) chi.Router{
	repo := authoutadapters.NewUserRepository(db)
	jwt := authinadapters.NewJwtAdapterService()
	authProviderService := authoutadapters.NewAuthProviderAdapter() 

	// middlewares 
	authMiddleware := authmiddlewares.NewAuthMiddleware(
		jwt,
	)

	// 2. Casos de Uso de Apoyo (Capa de Aplicación)
	// Estos son dependencias de los casos de uso principales
	findUserByEmailUC := authusecases.NewFindUserByEmailUseCase(repo)
	updateSessionUC := authusecases.NewUpdateSessionUseCase(repo)

	// El RegisterUseCase suele necesitar el repo para crear el usuario
	registerUC := authusecases.NewRegisterUseCase(repo)

	// 3. Caso de Uso Principal (Orquestador)
	// Inyectamos todas las dependencias necesarias
	loginUC := authusecases.NewLoginUseCase(
		jwt, 
		repo, 
		authProviderService, 
		registerUC, 
		findUserByEmailUC, 
		updateSessionUC,
	)

	// buscar usuario caso de uso
	findUserUC:= authusecases.NewFindUserByEmailUseCase(repo)

	// refrescar token, caso de uso
	refreshTokenUC:=authusecases.NewRefreshTokenUseCase(jwt)
	
	controller := authcontrollers.NewAuthController(
		loginUC,
		authMiddleware,
		findUserUC,
		refreshTokenUC,
	)

	return authcontrollers.AuthMapRoutes(controller)
}