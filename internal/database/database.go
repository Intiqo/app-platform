package database

import (
	"context"
	"fmt"
	"log"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Intiqo/app-platform/internal/pkg/config"
)

// NewDB returns a new database connection pool
func NewDB(cfg config.AppConfig) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DatabaseUsername, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)

	// Create a database connection pool
	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	// Register the uuid type
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	// Create the connection pool
	dbp, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return dbp
}
