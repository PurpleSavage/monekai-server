package connection

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient inicializa la conexión a GORM con reintentos (útil para docker-compose)
func NewClient(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxRetries := 10
	baseDelay := time.Second

	for i := range maxRetries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		delay := baseDelay * (1 << i)
		log.Printf(
			"Intento %d/%d: no se pudo conectar a la base de datos. Reintentando en %v...",
			i+1, maxRetries, delay,
		)
		time.Sleep(delay)
	}

	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos: %w", maxRetries, err)
	}

	// Configuración del pool de conexiones
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}