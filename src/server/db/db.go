package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	ctx  *context.Context
}

func Connect(ctx context.Context, pgConnectionString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, pgConnectionString)

	if err != nil {
		return nil, fmt.Errorf("unable to connect")
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db")
	}

	db := &DB{
		pool: pool,
		ctx:  &ctx,
	}

	return db, nil
}
