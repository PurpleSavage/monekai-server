package authcontrollers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	authrequestsdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/requests"
	authresponsesdtos "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/dtos/responses"
	authusecases "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/usecases"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	authutils "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/utils"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
	"github.com/go-chi/chi/v5"
)

type AuthController struct {
	loginUseCase *authusecases.LoginUseCase
    authMiddleware *authmiddlewares.AuthMiddleware
    findUserUseCase *authusecases.FindUserByEmailUseCase
    refreshTokenUseCase *authusecases.RefreshTokenUseCase
}

func NewAuthController(
	lu *authusecases.LoginUseCase,
    am *authmiddlewares.AuthMiddleware,
    fu *authusecases.FindUserByEmailUseCase,
    ru *authusecases.RefreshTokenUseCase,
) *AuthController {
	return &AuthController{
		loginUseCase: lu,
        authMiddleware: am,
        findUserUseCase: fu,
        refreshTokenUseCase:ru,
	}
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req authrequestsdtos.AuthRequestDto
    userAgent:=r.UserAgent()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        commoninfrastructuremappers.RespondWithError(
	        w, 
	        globalerrors.NewAppError(400, "Bad Request", "Invalid JSON body", err),
        )
        return
    }
	if req.Token == "" {
        commoninfrastructuremappers.RespondWithError(
	        w, 
	        globalerrors.NewAppError(400, "Validation Error", "Token is required", nil),
        )
        return
    }
    if userAgent==""{
        userAgent = "unknown"
    }
    session, err := ac.loginUseCase.Execute(req.Token,userAgent) 
    if err != nil {
    	log.Println("ERROR en handler Login:", err)
        commoninfrastructuremappers.RespondWithError(w, err)
        return
    }
    refreshToken:= session.RefreshToken
    responseSession:= authresponsesdtos.ResponseSessionDto{
        UserData:authresponsesdtos.UserResponseDto{
            ID:        session.UserData.Id,
            Email:     session.UserData.Email,
            PhotoURL:  session.UserData.PhotoUrl,
            CreatedAt: session.UserData.CreatedAt.Format(time.RFC3339),
            Credits:   session.UserData.Credits,
        },
        AccessToken: session.AccessToken,
    }
    authutils.HandleCookie(w,refreshToken)
    commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, responseSession)
}

func (ac *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
 	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
    if rawData == nil {
		commoninfrastructuremappers.RespondWithError(
			w, 
			globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil),
		)
		return
	}
	dto, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(
			w, 
			globalerrors.NewAppError(500, "Internal Error", "Could not parse session data", err),
		)
		return
	}
	    newtoken, err := ac.refreshTokenUseCase.Execute(dto)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	response := authresponsesdtos.RefreshTokenResponseDto{
		AccessToken: newtoken.Value(),
	}
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK, response)
}


func (ac *AuthController) GetProfile(w http.ResponseWriter, r *http.Request) {
	rawData := r.Context().Value(authmiddlewares.SessionContextKey)
	if rawData == nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(401, "Unauthorized", "Session data not found in context", nil),
		)
		return
	}
	
	// 2. Lo mapeamos usando TU función MapClaimsToStruct
	dto, err := authinadapters.MapClaimsToStruct[authrequestsdtos.SessionRequestDto](rawData)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(
			w,
			globalerrors.NewAppError(500, "Internal Error", "Could not parse session data", err),
		)
		return
	}
	newtoken, err := ac.refreshTokenUseCase.Execute(dto)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return 
	}
	
	// 3. Extraemos el Email directamente del DTO mapeado
	session, err := ac.findUserUseCase.Execute(dto.Email)
	if err != nil {
		commoninfrastructuremappers.RespondWithError(w, err)
		return
	}
	//TODO: usar el dto que ya existe en login 
	userDto := authresponsesdtos.UserResponseDto{
		ID:        session.Id,
		Email:     session.Email,
		PhotoURL:  session.PhotoUrl,
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		Credits:   session.Credits,
	}
	response := authresponsesdtos.RenovateSessionResponseDto{
		UserData:        userDto,
		AccessToken: newtoken.Value(),
	} 
	
	commoninfrastructuremappers.RespondWithJSON(w, http.StatusOK,response)
}



func AuthMapRoutes(ac *AuthController) chi.Router {
	r := chi.NewRouter()
	//pruta pública 
	r.Post("/login", ac.Login)

	//ruta con middleware
	r.Group(func(r chi.Router){
		r.Use(ac.authMiddleware.RefreshToken)
		r.Get("/refresh-token", ac.RefreshToken)
	})

	//ruta con middleware
	r.Group(func(r chi.Router){
		r.Use(ac.authMiddleware.RefreshToken)
		r.Get("/profile", ac.GetProfile)
	})
	return r
} 