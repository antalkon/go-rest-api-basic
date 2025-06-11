package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

var JwtSecretKey = []byte("jwtSecret") // секретный ключ для подписи

// Claims структура с UUID как UserID
type Claims struct {
	UserID uuid.UUID `json:"sub"`  // ID пользователя как UUID
	Role   string    `json:"role"` // Роль (если нужно)
	jwt.RegisteredClaims
}

// GenerateAccessToken создает access-токен с UUID
func GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(), // JWT требует string в subject
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

// GenerateRefreshToken — токен с subject как UUID (в виде строки)
func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

// ParseAccessToken возвращает Claims c UUID внутри
// ParseAccessToken возвращает Claims c UUID внутри, даже если токен истёк (для refresh flow)
func ParseAccessToken(tokenStr string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithLeeway(0), jwt.WithValidMethods([]string{"HS256"}), jwt.WithoutClaimsValidation())

	token, _, err := parser.ParseUnverified(tokenStr, &Claims{}) // сначала без валидации (например, чтобы вытянуть ID)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid access token claims")
	}

	// Теперь валидируем вручную с учетом того, что он может быть просрочен
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token is expired")
	}

	// Проверка подписи и валидности токена
	validToken, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !validToken.Valid {
		return nil, errors.New("invalid token signature")
	}

	return claims, nil
}

// ParseRefreshToken возвращает UUID из refresh токена
func ParseRefreshToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid refresh token")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID in token subject")
	}

	return id, nil
}
