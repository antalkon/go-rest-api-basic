package app

import (
	"backend/config"
	"backend/intenal/composition"
	"backend/intenal/transport/rest"
	"backend/pkg/httpserver"
	"backend/pkg/logger"
	"backend/pkg/postgres"
	"backend/pkg/redis"
	"backend/pkg/s3"
	"context"
	"time"
)

func Run(cfg *config.Config) {
	logger.Init(cfg.Log.Level)
	l := logger.L()

	pg, err := postgres.New(postgres.Config{
		URL:      cfg.PG.URL,
		PoolSize: cfg.PG.PoolMax,
	})
	if err != nil {
		l.Fatal("Postgres init failed: %v", err)
	}
	defer pg.Close()
	l.Info("Postgres connected")

	rds, err := redis.New(redis.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		l.Fatal("Redis error: %v", err)
	}
	defer rds.Close()
	l.Info("Redis connected")

	s3Client, err := s3.New(s3.Config{
		Endpoint:  cfg.S3.Endpoint,
		AccessKey: cfg.S3.AccessKey,
		SecretKey: cfg.S3.SecretKey,
		Bucket:    cfg.S3.Bucket,
		UseSSL:    cfg.S3.UseSSL,
	})
	if err != nil {
		l.Fatal("S3 init error: %v", err)
	}
	_ = s3Client // используем s3Client в дальнейшем коде
	l.Info("S3 connected and bucket ready")

	modules := composition.InitRESTModules(pg)
	srv := httpserver.New(httpserver.Config{
		Address:         ":" + cfg.HTTP.Port,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		ShutdownTimeout: 3 * time.Second,
	})
	rest.InitRoutes(srv, modules)
	srv.Start()
	l.Info("HTTP server started on %s", cfg.HTTP.Port)
	if err := <-srv.Notify(); err != nil {
		l.Error("Server error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	l.Info("App %s started", cfg.App.Name)
}
