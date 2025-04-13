package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Provider struct {
	logger *slog.Logger
	Db     *sql.DB
}

var validatedStructure bool = false

func OpenDB(logger *slog.Logger, dbFile string) (*sql.DB, error) {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Connecting to database...")
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFile))
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Error opening sqlite database", slog.String("error", err.Error()))
		return nil, err
	}
	if !validatedStructure {
		_, err = checkDB(db)
		if err != nil {
			if strings.Contains(err.Error(), "no such table: versions") {
				err = createDBStructure(db)
				if err != nil {
					logger.LogAttrs(context.Background(), slog.LevelError, "Error creating database structure", slog.String("error", err.Error()))
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		validatedStructure = true
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully connected to database")
	return db, nil
}

func checkDB(db *sql.DB) (float64, error) {
	rows, err := db.Query("SELECT * FROM versions;")
	if err != nil {
		return -1, fmt.Errorf("error selecting versions from database: %w", err)
	}
	var v float64
	var versions []float64
	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			return -1, fmt.Errorf("error scanning version into int: %w", err)
		}
		versions = append(versions, v)
	}
	v = 0
	for _, version := range versions {
		if version > v {
			v = version
		}
	}
	return v, nil
}

func createDBStructure(db *sql.DB) error {
	_, err := db.Exec(createStructureQuery)
	if err != nil {
		return err
	}
	return nil
}
