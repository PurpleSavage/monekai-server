package scripts

import (
	"log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" 
    _ "github.com/golang-migrate/migrate/v4/source/file"       
)
func RunMigrations(dsn string) {
    m, err := migrate.New(
        "file://migrations",
        dsn,
    )
    if err != nil {
        log.Fatal("Error creando migrator:", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Error corriendo migraciones:", err)
    }
    log.Println("Migraciones OK")
}