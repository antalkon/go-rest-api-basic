package main

import (
	_ "backend/api"
	config "backend/config"
	"backend/intenal/app"
	"log"

	"github.com/joho/godotenv"
)

// @title           Backend go basic API
// @version         1.0
// @description     Basic API template for Go applications.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
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
