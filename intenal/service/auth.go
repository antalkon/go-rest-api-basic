package service

import (
	"backend/intenal/models"
	repo "backend/intenal/repo"
	"backend/intenal/transport/rest/v1/req"
	"backend/pkg/auth"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	repo *repo.AuthRepository
}

func NewAuthService(repo *repo.AuthRepository) *AuthService {
	if repo == nil {
		panic("AuthService: repo is nil")
	}
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, r req.RegisterRequest) (string, string, error) {
	existing, _ := s.repo.FindByEmail(ctx, r.Email)
	if existing != nil {
		return "", "", errors.New("email already registered")
	}

	hash, err := auth.HashPassword(r.Password)
	if err != nil {
		return "", "", err
	}
	userUuid := uuid.New()

	if err := s.repo.CreateUser(ctx, userUuid, r.Username, r.Email, hash); err != nil {
		return "", "", err
	}

	access, err := auth.GenerateAccessToken(userUuid, "user")
	if err != nil {
		return "", "", err
	}
	refresh, err := auth.GenerateRefreshToken(userUuid)
	if err != nil {
		return "", "", err
	}

	refreshToken := models.RefreshToken{
		ID:        uuid.New(), // ID самой записи
		UserID:    userUuid,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(auth.RefreshTokenTTL),
		Revoked:   false,
	}

	if err := s.repo.SaveRefresh(ctx, refreshToken); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) Login(ctx context.Context, r req.LoginRequest) (string, string, error) {
	user, err := s.repo.FindByEmail(ctx, r.Email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := auth.ComparePassword(r.Password, user.PasswordHash); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, err := auth.GenerateAccessToken(user.ID, "user")
	if err != nil {
		return "", "", err
	}
	refresh, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken := models.RefreshToken{
		ID:        uuid.New(), // ID самой записи
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(auth.RefreshTokenTTL),
		Revoked:   false,
	}

	if err := s.repo.SaveRefresh(ctx, refreshToken); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
