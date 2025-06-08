package main

import (
	config "backend/config"
	"backend/internal/app"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env файла (опционально, для локальной разработки)
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// Конфигурация
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Запуск приложения
	app.Run(cfg)
}
