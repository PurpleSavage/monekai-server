package sampler

import (
	samplercontroller "github.com/PurpleSavage/monekai-server/modules/sampler/infrastructure/controllers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)
func SamplerBootstrap(db *gorm.DB) chi.Router{
	controller := samplercontroller.NewSamplerController()
	return samplercontroller.SamplerMapRoutes(controller)
}