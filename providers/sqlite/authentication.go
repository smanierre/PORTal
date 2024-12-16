package sqlite

import (
	"PORTal/types"
	"context"
	"log/slog"
	"time"
)

func (p Provider) CreateSession(memberID, sessionID, userAgent, ipAddress string, expiration time.Time) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating session in database")
	_, err := p.Db.Exec(insertSessionQuery, sessionID, userAgent, ipAddress, expiration)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Failed to create session in database", slog.String("error", err.Error()))
		return err
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Created session in database, creating mapping to user")
	_, err = p.Db.Exec(insertMemberSessionQuery, memberID, sessionID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Failed to create member_session in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p Provider) GetSession(sessionID string) (types.Session, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting session from database", slog.String("session_id", sessionID))
	row := p.Db.QueryRow(getSessionQuery, sessionID)
	var s types.Session
	err := row.Scan(&s.SessionID, &s.Expires, &s.UserAgent, &s.IpAddress)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting session from database", slog.String("error", err.Error()))
		return types.Session{}, err
	}
	return s, nil
}
