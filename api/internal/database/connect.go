package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseUrl string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatal("Failed to create connection pool", err.Error())
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping database", err.Error())
	}

	return conn
}
