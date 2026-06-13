package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"carro-ideal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var database *sql.DB

func Connect(cfg *config.Config) error {
	dsn := cfg.DatabaseURL()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	database = db
	return nil
}

func GetDB() *sql.DB {
	return database
}

func RunMigrations(databaseURL string) error {
	if database == nil {
		return fmt.Errorf("database has not been initialized")
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migrate driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func GetDBHealth(ctx context.Context) error {
	if database == nil {
		return fmt.Errorf("database is not initialized")
	}
	return database.PingContext(ctx)
}
