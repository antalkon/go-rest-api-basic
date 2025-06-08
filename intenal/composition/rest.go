package composition

import (
	repository "backend/intenal/repo"
	"backend/intenal/service"
	"backend/intenal/transport/rest/v1/handlers"
	"backend/pkg/postgres"
)

type RESTModules struct {
	PingHandler *handlers.PingHandler
}

func InitRESTModules(pg *postgres.Postgres) *RESTModules {
	repo := repository.NewPingRepository(pg)
	svc := service.NewPingService(repo)
	handler := handlers.NewPingHandler(svc)

	return &RESTModules{
		PingHandler: handler,
	}
}
