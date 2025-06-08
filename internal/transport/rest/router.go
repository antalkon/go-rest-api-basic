package rest

import (
	"backend/internal/composition"
	"backend/internal/transport/rest/middleware"
	v1 "backend/internal/transport/rest/v1"
	"backend/pkg/httpserver"
)

func InitRoutes(server *httpserver.Server, m *composition.RESTModules) {
	api := server.App().Group("/api")
	v1Group := api.Group("/v1", middleware.RequestLogger())

	v1.RegisterRoutes(v1Group, m.PingHandler)
}
