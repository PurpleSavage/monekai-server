package authmiddlewares

import (
	"context"
	"net/http"
	"strings"
	authports "github.com/PurpleSavage/monekai-server/modules/shared/auth/application/ports"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	commoninfrastructuremappers "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/mappers"
)

type contextKey string

const SessionContextKey contextKey = "user_session"

type AuthMiddleware struct{
	JwtService authports.JwtPort
}
func NewAuthMiddleware(JwtService authports.JwtPort) *AuthMiddleware {
	return &AuthMiddleware{
		JwtService: JwtService,
	}
}
func (a *AuthMiddleware) AccessToken(next http.Handler)http.Handler{
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "Missing or invalid token", nil))
            return
        }
		token:= strings.Split(authHeader," ")[1]
		userData, err:= a.JwtService.VerifyToken(token)
		if err != nil {
			commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "token expired", nil))
			return 
		}
		ctx := context.WithValue(r.Context(), SessionContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (a *AuthMiddleware) RefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "Missing or invalid refresh token", nil))
			return
		}
		
		// extraemos el payload estructurado
		userData, err := a.JwtService.VerifyToken(cookie.Value)
		if err != nil {
			commoninfrastructuremappers.RespondWithError(w, globalerrors.NewAppError(401, "Unauthorized", "Expired refresh token", nil))
			return
		}

		ctx := context.WithValue(r.Context(), SessionContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
