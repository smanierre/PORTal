package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"strings"
)

func New(logger *slog.Logger, dbFile string, expectedVersion float64) (*Backend, error) {
	l := logger.With(slog.String("source", "sqlite3_backend"))
	b := &Backend{
		logger: l,
		Db:     nil,
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "Connecting to database...")
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFile))
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error opening sqlite database", slog.String("error", err.Error()))
		return nil, err
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "Successfully connected to database")
	l.LogAttrs(context.Background(), slog.LevelInfo, "Checking database structure...")
	version, err := checkDB(db)
	if err != nil && strings.Contains(err.Error(), "no such table: version") {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Structure not detected, creating...")
		err := createDBStructure(db)
		if err != nil {
			l.LogAttrs(context.Background(), slog.LevelError, "Error creating database structure", slog.String("error", err.Error()))
			return nil, err
		}
		l.LogAttrs(context.Background(), slog.LevelInfo, "Successfully created database structure")
	} else if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error checking database", slog.String("error", err.Error()))
		return nil, err
	}
	if version != expectedVersion {
		l.LogAttrs(context.Background(), slog.LevelWarn, "Different Db version detected vs expected, upgrades not implemented yet")
		b.Db = db
		return b, nil
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "Found correct structure and version")
	b.Db = db
	return b, nil
}

type Backend struct {
	logger *slog.Logger
	Db     *sql.DB
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
