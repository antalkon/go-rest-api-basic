package repository

import (
	"backend/intenal/models"
	"backend/pkg/postgres"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type AuthRepository struct {
	db *postgres.Postgres
}

func NewAuthRepository(db *postgres.Postgres) *AuthRepository {
	if db == nil {
		panic("AuthRepository: nil db passed")
	}
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, uuid uuid.UUID, username, email, passwordHash string) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)`,
		uuid, username, email, passwordHash,
	)
	return err
}

func (r *AuthRepository) SaveRefresh(ctx context.Context, refresh models.RefreshToken) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO refresh_tokens (id, user_id, created_at, expires_at, revoked) VALUES ($1, $2, $3, $4, $5)`,
		refresh.ID, refresh.UserID, refresh.CreatedAt, refresh.ExpiresAt, refresh.Revoked,
	)
	return err
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx, `SELECT id, username, email, password_hash FROM users WHERE email = $1`, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
