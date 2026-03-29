package database

import (
	"context"
	"embed"
	"fmt"
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
	dbMigrations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Migration])
	if err != nil {
		return err
	}

	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

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
				slog.Info(fmt.Sprintf("Skipping migration %s as it was already ran", m.Name()))
				break
			}
		}

		if alreadyRan {
			continue
		}

		// Execute the migration
		migrationSql, err := migrationFiles.ReadFile("migrations/" + m.Name())
		if err != nil {
			slog.Error("Error reading migration file", m.Name(), err)
			return err
		}
		_, err = tx.Exec(context.Background(), string(migrationSql))
		if err != nil {
			slog.Error("Error executing migration", m.Name(), err)
			return err
		}

		// Update the migrations table
		_, err = tx.Exec(context.Background(), `
		INSERT INTO migrations (name) VALUES ($1)`, m.Name())
		if err != nil {
			slog.Error("Error executing migration table update", m.Name(), err)
			return err
		}

		slog.Info(fmt.Sprintf("Migration %s applied", m.Name()))
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
