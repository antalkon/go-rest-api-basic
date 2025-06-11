package v1

import (
	"backend/intenal/transport/rest/v1/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, hPing *handlers.PingHandler, hAuth *handlers.AuthHandler) {
	g.GET("/ping", hPing.Ping)
	g.GET("/ping/all", hPing.GetAll)

	authGroup := g.Group("/auth")
	{
		authGroup.POST("/login", hAuth.Login)
		authGroup.POST("/register", hAuth.Register)
		authGroup.POST("/refresh", hAuth.Refresh)
		authGroup.POST("/logout", hAuth.Logout)
	}

}
