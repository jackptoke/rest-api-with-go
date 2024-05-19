package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Database struct {
	Client *sqlx.DB
}

//type MyDatabase interface {
//	NewDatabase() (*Database, error)
//	Ping(context.Context) error
//	MigrateDB() error
//}

// Schema - for testing only
var schema = `
CREATE TABLE IF NOT EXISTS comments (
    ID uuid,
    Slug text,
    Author text,
    Body text
);`

func NewDatabase() (*Database, error) {
	connxString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_DB"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	fmt.Println("Connecting to database at: ", connxString)

	dbConn, err := sqlx.Connect("postgres", "host=postgres port=5432 user=postgres password=password dbname=postgres sslmode=disable timezone=UTC connect_timeout=5")
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}
	return &Database{
		Client: dbConn,
	}, nil
}

func NewTestDatabase() (*Database, error) {

	fmt.Println("Connecting to database at: ")

	dbConn, err := sqlx.Connect("sqlite3", "fiberdb.db?cache=shared&mode=rwc")

	dbConn.MustExec(schema)

	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}
	return &Database{
		Client: dbConn,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}

// MigrateDB - runs all migrations in the migrations
func (d *Database) MigrateDB() error {
	log.Info("migrating our database")
	driver, err := postgres.WithInstance(d.Client.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not instantiate database driver: %w", err)
	}
	fmt.Println("Database driver is", driver)

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres",
		driver)

	fmt.Println("Database migration is", m)

	if err != nil {
		log.Errorf("could not instantiate migration instance: %v", err)
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run migration: %w", err)
		}
		return err
	}

	fmt.Println("Successfully instantiated migration")
	return nil
}
