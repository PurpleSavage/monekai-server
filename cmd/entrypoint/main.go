package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/PurpleSavage/monekai-server/cmd/scripts"
	connection "github.com/PurpleSavage/monekai-server/configurations/persistenceconnections"
	"github.com/PurpleSavage/monekai-server/modules/community"
	"github.com/PurpleSavage/monekai-server/modules/notifications"
	notificationsevents "github.com/PurpleSavage/monekai-server/modules/notifications/infrastructure/serverevents"
	"github.com/PurpleSavage/monekai-server/modules/sampler"
	"github.com/PurpleSavage/monekai-server/modules/shared/auth"
	authinadapters "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/in-adapters"
	authmiddlewares "github.com/PurpleSavage/monekai-server/modules/shared/auth/infrastructure/middlewares"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
	commoninadapters "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/in-adapters"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/validators"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{
			"Content-Type",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		AllowCredentials: true, 
	})
	r.Use(corsHandler.Handler)
	config.LoadEnvs()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Envs.Host,
		config.Envs.DbUser,
		config.Envs.DbPassword,
		config.Envs.DbName,
		config.Envs.DbPort,
		config.Envs.SslMode,
	)
	migrateDSN := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        config.Envs.DbUser,
        config.Envs.DbPassword,
        config.Envs.Host,
        config.Envs.DbPort,
        config.Envs.DbName,
        config.Envs.SslMode,
    )
	db, err := connection.NewClient(dsn)
	
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	scripts.RunMigrations(migrateDSN)
	// ─── 1. INSTANCIAS GLOBALES COMPARTIDAS ───
	jwt := authinadapters.NewJwtAdapterService()
    dtoValidator := validators.NewDTOValidator()
    bucketObserver:= commoninadapters.NewObserverBucket()
    authmiddleware:=authmiddlewares.NewAuthMiddleware(jwt)

    //sse
    sseManager := notificationsevents.NewSSEManager()
	audioSSEHandler := notificationsevents.NewAudioSSEHandler(sseManager,authmiddleware)

    //root sse
    notificationsevents.MapSSERoutes(r,audioSSEHandler)
    //root routes
	r.Mount(
		"/auth",
		auth.AuthBootstrap(
			db,
		),
	)
	r.Mount(
		"/notifications",
		notifications.NotificationsBootstrap(
			db,
			bucketObserver,
			dtoValidator,
			authmiddleware,
			sseManager,
		),
	)
	r.Mount(
		"/audio",
		sampler.SamplerBootstrap(
			db,
			bucketObserver,
			dtoValidator,
			authmiddleware,
		),
	)
	r.Mount(
		"/community",
		community.CommunityBootstrap(
			db,
			dtoValidator,
			authmiddleware,
		),
	)
	
	
	http.ListenAndServe(":8080", r)
}
