package qualificationprovider

import (
	"PORTal/providers/sqlite"
	"database/sql"
	"log/slog"
)

type QualificationProvider struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(dbFile string, logger *slog.Logger) (QualificationProvider, error) {
	l := logger.With(slog.String("source", "QualificationProvider"))
	db, err := sqlite.OpenDB(l, dbFile)
	if err != nil {
		return QualificationProvider{}, err
	}
	return QualificationProvider{db: db, logger: l}, nil
}
