package service

import (
	"backend/intenal/models"
	repository "backend/intenal/repo"
	"context"
)

type PingService struct {
	repo *repository.PingRepository
}

func NewPingService(repo *repository.PingRepository) *PingService {
	if repo == nil {
		panic("PingService: repo is nil")
	}
	return &PingService{repo: repo}
}

func (s *PingService) SavePing(ctx context.Context, ip string) error {
	return s.repo.SavePing(ctx, ip)
}

func (s *PingService) GetPings(ctx context.Context) ([]models.Ping, error) {
	return s.repo.GetAll(ctx)
}
