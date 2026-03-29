package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func GetMigrations(pool *pgxpool.Pool) ([]Migration, error) {
	rows, err := pool.Query(context.Background(), "SELECT * FROM migrations ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	currentMigrations, err := pgx.CollectRows(rows, pgx.RowTo[Migration])
	if err != nil {
		return nil, err
	}

	return currentMigrations, nil

}
