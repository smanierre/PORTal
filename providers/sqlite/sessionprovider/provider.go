package sessionprovider

import (
	"PORTal/providers/sqlite"
	"database/sql"
	"log/slog"
)

type SessionProvider struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(dbFile string, logger *slog.Logger) (SessionProvider, error) {
	l := logger.With(slog.String("source", "SessionProvider"))
	db, err := sqlite.OpenDB(l, dbFile)
	if err != nil {
		return SessionProvider{}, err
	}
	return SessionProvider{db: db, logger: l}, nil
}
