package service

import (
	repo "backend/intenal/repo"
	"context"
)

type DataService struct {
	repo *repo.DataRepository
}

func NewDataService(repo *repo.DataRepository) *DataService {
	if repo == nil {
		panic("AuthService: repo is nil")
	}
	return &DataService{repo: repo}
}

func (s *DataService) GetUserData(ctx context.Context, userID string) (string, string, error) {
	data, err := s.repo.GetUserData(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Assuming data is a struct with Name and Email fields
	return data.Username, data.Email, nil
}
