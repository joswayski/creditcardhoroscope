package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migration struct {
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func RunMigrations(pool *pgxpool.Pool) error {
	// Create the migrations table if it doesn't exists
	_, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS migrations (
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	if err != nil {
		return err
	}
	// Get the current migrations in the DB
	rows, err := pool.Query(context.Background(), "SELECT * FROM migrations ORDER BY created_at ASC")
	if err != nil {
		return err
	}
	dbMigrations, err := pgx.CollectRows(rows, pgx.RowTo[Migration])
	if err != nil {
		return err
	}

	// Get the current migrations in /migrations
	localMigrations, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	// Loop through the local ones, and if not in the DB, execute it
	for _, m := range localMigrations {
		alreadyRan := false
		for _, dbM := range dbMigrations {
			if m.Name() == dbM.Name {
				alreadyRan = true
				break
			}

			if alreadyRan {
				continue
			}
		}

		// Execute the migration
		migrationSql, err := migrationFiles.ReadFile("migrations/" + m.Name())
		if err != nil {
			log.Fatal("Error reading migration file", err)
		}

		_, err = pool.Exec(context.Background(), string(migrationSql))
		if err != nil {
			log.Fatal("Error executing migration", err)
		}

		slog.Info(fmt.Sprintf("Migration %s applied", m.Name()))
	}

	return err
}
