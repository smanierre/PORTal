package memberprovider

import (
	"PORTal/providers/sqlite"
	"database/sql"
	"log/slog"
)

type MemberProvider struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(dbFile string, logger *slog.Logger) (MemberProvider, error) {
	l := logger.With(slog.String("source", "MemberProvider"))
	db, err := sqlite.OpenDB(l, dbFile)
	if err != nil {
		return MemberProvider{}, err
	}
	return MemberProvider{db: db, logger: l}, nil
}
