package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

// Opens the connection to the data base
func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}
	
	fmt.Println("Connected to database")
	return db, nil
}


func MigrateFS(db *sql.DB, migrationsFs fs.FS, dir string) error {
	// This tells Goose to use an embedded file system instead of reading SQL files from disk.
	goose.SetBaseFS(migrationsFs)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Migrate: %w", err)
	}

	// Regardless of embedded file system or not, goose still needs to know the path.
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}


	return nil
}