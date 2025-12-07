package postgres

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	goosev3 "github.com/pressly/goose/v3"
)

func New(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool, dir string, logger *log.Logger) error {
	var db *sql.DB = stdlib.OpenDB(*pool.Config().ConnConfig)
	defer func(db *sql.DB) {
		if db == nil {
			return
		}
		if err := db.Close(); err != nil {
			if logger != nil {
				logger.Printf("postgres.Migrate: failed to close DB handle: %v", err)
			} else {
				log.Printf("postgres.Migrate: failed to close DB handle: %v", err)
			}
		}
	}(db)

	goosev3.SetBaseFS(nil)
	goosev3.SetLogger(logger)

	if err := goosev3.SetDialect("postgres"); err != nil {
		return err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if err := goosev3.UpContext(ctx, db, absDir); err != nil {
		return err
	}
	return nil
}
