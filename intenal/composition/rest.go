package composition

import (
	repository "backend/intenal/repo"
	"backend/intenal/service"
	"backend/intenal/transport/rest/v1/handlers"
	"backend/pkg/postgres"
)

type RESTModules struct {
	PingHandler *handlers.PingHandler
	AuthHandler *handlers.AuthHandler
	DataHandler *handlers.DataHandler
}

func InitRESTModules(pg *postgres.Postgres) *RESTModules {
	return &RESTModules{
		PingHandler: handlers.NewPingHandler(service.NewPingService(repository.NewPingRepository(pg))),
		AuthHandler: handlers.NewAuthHandler(service.NewAuthService(repository.NewAuthRepository(pg))),
		DataHandler: handlers.NewDataHandler(service.NewDataService(repository.NewDataRepository(pg))),
	}
}
