package sessionprovider

import (
	"PORTal/types"
	"context"
	"log/slog"
	"time"
)

func (s SessionProvider) CreateSession(memberID, sessionID, userAgent, ipAddress string, expiration time.Time) error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating session in database")
	_, err := s.db.Exec(insertSessionQuery, sessionID, userAgent, ipAddress, expiration)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Failed to create session in database", slog.String("error", err.Error()))
		return err
	}
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Created session in database, creating mapping to user")
	_, err = s.db.Exec(insertMemberSessionQuery, memberID, sessionID)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Failed to create member_session in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s SessionProvider) GetSession(sessionID string) (types.Session, error) {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting session from database", slog.String("session_id", sessionID))
	row := s.db.QueryRow(getSessionQuery, sessionID)
	var session types.Session
	err := row.Scan(&session.SessionID, &session.Expires, &session.UserAgent, &session.IpAddress)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting session from database", slog.String("error", err.Error()))
		return types.Session{}, err
	}
	return session, nil
}

func (s SessionProvider) DeleteSession(sessionID string) error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting session from database", slog.String("session_id", sessionID))
	_, err := s.db.Exec(deleteSessionQuery, sessionID)
	if err != nil {
		return err
	}
	return nil
}
