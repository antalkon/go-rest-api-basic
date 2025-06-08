APP_NAME=backend # Название бинарника (опционально)
MAIN=./cmd/backend # Путь до main-файла
MIGRATE_TAGS=migrate # Миграции

# Запуск приложения
run:
	go run $(MAIN) 




# ==== MIGRATIONS ============
migrate-up: # Прогнать миграции вверх
	go run -tags $(MIGRATE_TAGS) $(MAIN)

migrate-down: # Прогнать миграции вниз (откат на 1 шаг)
	go run -tags $(MIGRATE_TAGS) $(MAIN) down

migrate-reset: # Откат всех миграций (необязательно)
	go run -tags $(MIGRATE_TAGS) $(MAIN) reset

migrate-create: # Создать новую миграцию (пример: make migrate-create name=add_users_table)
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name



# ==== УТИЛИТЫ ================
fmt: # Форматирование кода
	go fmt ./...
test: # Тесты
	go test ./...
clean: # Очистка
	go clean
# Переменные окружения (для отладки)
env:
	@echo "Using .env:"
	@cat .env