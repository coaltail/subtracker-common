package testing

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	pgtc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBCleanupFunc func(ctx context.Context) error

func CreatePostgresTestDB(ctx context.Context, initSQLPath string) (*sqlx.DB, DBCleanupFunc, error) {
	container, err := runPostgresContainer(ctx, initSQLPath)
	if err != nil {
		return nil, nil, fmt.Errorf("create test container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get connection string: %w", err)
	}
	connStr += " sslmode=disable"

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	cleanup := func(ctx context.Context) error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close database: %w", err)
		}
		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate container: %w", err)
		}
		return nil
	}

	return db, cleanup, nil
}

func runPostgresContainer(ctx context.Context, initSQLPath string) (*pgtc.PostgresContainer, error) {
	container, err := pgtc.Run(
		ctx,
		"postgres:16-alpine",
		pgtc.WithInitScripts(initSQLPath),
		pgtc.WithDatabase("postgres"),
		pgtc.WithUsername("postgres"),
		pgtc.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("run container: %w", err)
	}
	return container, nil
}
