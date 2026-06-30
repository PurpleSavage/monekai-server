package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PurpleSavage/monekai-server/cmd/scripts"
	connection "github.com/PurpleSavage/monekai-server/configurations/persistenceconnections"
	"github.com/PurpleSavage/monekai-server/modules/shared/auth"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
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
    //dtoValidator := sharedvalidators.NewDTOValidator()
	r.Mount(
		"/auth",
		auth.AuthBootstrap(
			db,
		),
	)
	
	
	http.ListenAndServe(":8080", r)
}
