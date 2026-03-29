package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseUrl string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
