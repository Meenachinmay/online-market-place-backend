package testing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"soda-interview/foundation/config"
	"soda-interview/foundation/database/postgres"
	"soda-interview/foundation/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBContainer manages the database connection for tests.
type DBContainer struct {
	Log *logger.Logger
	DB  *pgxpool.Pool
	Cfg *config.Config
}

// NewDBContainer sets up the database for testing.
func NewDBContainer(t *testing.T) *DBContainer {
	t.Helper()

	// 1. Initialize Logger
	log := logger.New(os.Stdout, "DEBUG")

	// 2. Locate Config
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../../") // zarf/testing -> ../../ -> root
	configPath := filepath.Join(projectRoot, "foundation/config")
	
	// 3. Load Config
	os.Setenv("APP_ENVIRONMENT", "test")
	cfg, err := config.LoadWithPath(configPath, "test")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 4. Connect to Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgres.New(ctx, cfg.GetDatabaseDSN())
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// 5. Run Migrations
	migrationDir := filepath.Join(projectRoot, "business/data/schema/migrations")
	if err := postgres.Migrate(ctx, db, migrationDir, log.NewStdLogger()); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return &DBContainer{
		Log: log,
		DB:  db,
		Cfg: cfg,
	}
}

// Teardown cleans up the database connection.
func (c *DBContainer) Teardown(t *testing.T) {
	t.Helper()
	c.DB.Close()
}

// Truncate wipes data from tables to ensure a clean slate.
func (c *DBContainer) Truncate(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tables := []string{
		"transactions",
		"orders",
		"blogs",
		"products",
		"wallets",
	}

	q := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", strings.Join(tables, ", "))
	if _, err := c.DB.Exec(ctx, q); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}
