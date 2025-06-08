package v1

import (
	repository "backend/internal/repo"
	"backend/internal/service"
	"backend/internal/transport/rest/v1/handlers"
	"backend/pkg/postgres"
)

type Module struct {
	PingHandler *handlers.PingHandler
}

func InitModules(db *postgres.Postgres) *Module {
	pingRepo := repository.NewPingRepository(db)
	pingService := service.NewPingService(pingRepo)
	pingHandler := handlers.NewPingHandler(pingService)

	return &Module{
		PingHandler: pingHandler,
	}
}
