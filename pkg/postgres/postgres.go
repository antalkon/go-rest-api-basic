// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

// Config — параметры подключения
type Config struct {
	URL      string
	PoolSize int
}

// Postgres обертка над pgxpool
type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

// New создаёт новое соединение с базой
func New(cfg Config) (*Postgres, error) {
	var (
		pool *pgxpool.Pool
		err  error
	)

	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse config: %w", err)
	}
	poolCfg.MaxConns = int32(cfg.PoolSize)

	attempts := defaultConnAttempts

	for attempts > 0 {
		pool, err = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if err == nil {
			break
		}

		time.Sleep(defaultConnTimeout)
		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres: connect failed after retries: %w", err)
	}

	return &Postgres{
		Pool:    pool,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// Close закрывает пул соединений
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
