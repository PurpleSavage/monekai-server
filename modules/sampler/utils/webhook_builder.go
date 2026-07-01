package samplerutils

import (
	"fmt"
	"net/url"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
)

func BuildWebhook(path string, queryParams ...map[string]string) string {
	baseUrl := config.Envs.BackendServerBaseUrl
	fullUrl := fmt.Sprintf("%s/audio/webhook/%s", baseUrl, path)

	// Si no pasaron ningún mapa, o pasaron uno vacío, retornamos la URL base
	if len(queryParams) == 0 || len(queryParams[0]) == 0 {
		return fullUrl
	}

	u, err := url.Parse(fullUrl)
	if err != nil {
		return fullUrl
	}

	// queryParams[0] contiene el primer mapa enviado
	q := u.Query()
	for key, value := range queryParams[0] {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}