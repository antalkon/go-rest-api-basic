package service

import (
	"backend/intenal/models"
	repo "backend/intenal/repo"
	"backend/intenal/transport/rest/v1/req"
	"backend/pkg/auth"
	"context"
	"errors"
	"fmt"
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
		Token:     refresh,
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
		Token:     refresh,
		Revoked:   false,
	}

	if err := s.repo.SaveRefresh(ctx, refreshToken); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	fmt.Println("=== REFRESH START ===")
	fmt.Println("Токен из куки:", refreshToken)

	// 1. Парсим токен
	userID, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("Ошибка парсинга токена:", err)
		return "", "", errors.New("invalid refresh token")
	}
	fmt.Println("User ID из токена:", userID)

	// 2. Ищем токен в БД
	tokenData, err := s.repo.GetRefreshTokoenByToken(ctx, refreshToken)
	if err != nil {
		fmt.Println("Токен не найден в БД:", err)
		return "", "", errors.New("refresh token not found")
	}

	if tokenData.Revoked {
		fmt.Println("Токен отозван")
		return "", "", errors.New("refresh token is revoked")
	}
	if tokenData.ExpiresAt.Before(time.Now()) {
		fmt.Println("Токен истёк")
		return "", "", errors.New("refresh token expired")
	}
	fmt.Println("Токен действителен")

	// 3. Находим пользователя
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		fmt.Println("Пользователь не найден:", err)
		return "", "", errors.New("user not found")
	}
	fmt.Println("Пользователь найден:", user.Email)

	// 4. Отзываем старый токен
	if err := s.repo.RevokeToken(ctx, refreshToken); err != nil {
		fmt.Println("Ошибка отзыва токена:", err)
		return "", "", errors.New("failed to revoke old refresh token")
	}
	fmt.Println("Старый токен отозван")

	// 5. Генерим новые токены
	access, err := auth.GenerateAccessToken(user.ID, "user")
	if err != nil {
		fmt.Println("Ошибка генерации access:", err)
		return "", "", err
	}
	refresh, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		fmt.Println("Ошибка генерации refresh:", err)
		return "", "", err
	}

	// 6. Сохраняем новый refresh токен
	newToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(auth.RefreshTokenTTL),
		Token:     refresh,
		Revoked:   false,
	}
	if err := s.repo.SaveRefresh(ctx, newToken); err != nil {
		fmt.Println("Ошибка сохранения нового токена:", err)
		return "", "", err
	}
	fmt.Println("Новый токен сохранён")
	fmt.Println("=== REFRESH SUCCESS ===")

	return access, refresh, nil
}
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {

	if err := s.repo.RevokeToken(ctx, refreshToken); err != nil {

	}

	return nil
}
