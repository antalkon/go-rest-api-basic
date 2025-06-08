package app

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"

	// источники и драйверы
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	retryTimeout    = time.Second
)

func init() {
	_ = godotenv.Load() // вот это добавь

	dsn, ok := os.LookupEnv("PG_URL")
	if !ok || dsn == "" {
		log.Fatalf("migrate: PG_URL not set")
	}

	dsn += "?sslmode=disable"

	var (
		m   *migrate.Migrate
		err error
	)

	for attempts := defaultAttempts; attempts > 0; attempts-- {
		m, err = migrate.New("file://migrations", dsn)
		if err == nil {
			break
		}
		log.Printf("Migrate: waiting for DB (%d attempts left)", attempts)
		time.Sleep(retryTimeout)
	}

	if err != nil {
		log.Fatalf("Migrate: failed to connect: %v", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate: no changes")
		return
	}

	log.Println("Migrate: success")
}
