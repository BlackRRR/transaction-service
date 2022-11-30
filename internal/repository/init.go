package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func InitDB(cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	pgPool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pgxpool")
	}

	return pgPool, nil
}
