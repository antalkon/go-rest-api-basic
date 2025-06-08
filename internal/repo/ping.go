package repository

import (
	"backend/internal/models"
	"backend/pkg/postgres"
	"context"
)

type PingRepository struct {
	db *postgres.Postgres
}

func NewPingRepository(db *postgres.Postgres) *PingRepository {
	if db == nil {
		panic("PingRepository: nil db passed")
	}
	return &PingRepository{db: db}
}

// SavePing сохраняет IP клиента при запросе
func (r *PingRepository) SavePing(ctx context.Context, ip string) error {
	query := `INSERT INTO pings (ip) VALUES ($1)`
	_, err := r.db.Pool.Exec(ctx, query, ip)
	return err
}

// GetAll возвращает все пинги
func (r *PingRepository) GetAll(ctx context.Context) ([]models.Ping, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id, timestamp, ip FROM pings ORDER BY timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Ping
	for rows.Next() {
		var p models.Ping
		if err := rows.Scan(&p.ID, &p.Timestamp, &p.IP); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}
