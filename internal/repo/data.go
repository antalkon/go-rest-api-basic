package repository

import (
	"backend/intenal/models"
	"backend/pkg/postgres"
	"context"
)

type DataRepository struct {
	db *postgres.Postgres
}

func NewDataRepository(db *postgres.Postgres) *DataRepository {
	if db == nil {
		panic("AuthRepository: nil db passed")
	}
	return &DataRepository{db: db}
}

func (r *DataRepository) GetUserData(ctx context.Context, userID string) (models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(context.Background(), query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
