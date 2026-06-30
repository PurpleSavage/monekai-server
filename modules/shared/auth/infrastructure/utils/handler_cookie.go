package authutils

import (
	"net/http"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
)
func HandleCookie(w http.ResponseWriter, token string) {
	sameSite := http.SameSiteLaxMode
    secure := false
	enviroment:= config.Envs.Enviroment
	if enviroment == "prod" {
		sameSite = http.SameSiteNoneMode
		secure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		MaxAge:   86400,
		SameSite: sameSite,
	})
}