package httpserver

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	app     *echo.Echo
	address string
	notify  chan error
}

type Config struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func New(cfg Config) *Server {
	app := echo.New()

	// базовые middleware
	app.Use(middleware.Recover())

	s := &Server{
		app:     app,
		address: cfg.Address,
		notify:  make(chan error, 1),
	}

	return s
}

func (s *Server) App() *echo.Echo {
	return s.app
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.app.Start(s.address)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.Shutdown(ctx)
}
