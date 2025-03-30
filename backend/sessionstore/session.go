package sessionstore

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

func (s SessionStore) CreateSession(memberID, userAgent, ipAddress string) (string, time.Time) {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to create session", slog.String("memberID", memberID))
	sessionId := uuid.NewString()
	expiration := s.clock.Now().Add(s.timeout).UTC()
	err := s.provider.CreateSession(memberID, sessionId, userAgent, ipAddress, expiration)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Session creation failed")
		return "", time.Now()
	}
	return sessionId, expiration
}

func (s SessionStore) ValidateSession(sessionID, userAgent, ipAddress string) (types.Member, error) {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to validate session")
	session, err := s.provider.GetSession(sessionID)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelWarn, "Session validation failed due to database")
		return types.Member{}, fmt.Errorf("%w: %s", backend.ErrSessionValidationFailed, err.Error())
	}
	if userAgent != session.UserAgent {
		s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Session validation failed due to user agent mismatch")
		return types.Member{}, backend.ErrSessionValidationFailed
	}
	if ipAddress != session.IpAddress {
		s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Session validation failed due to IP mismatch")
		return types.Member{}, backend.ErrSessionValidationFailed
	}
	return s.memberStore.GetMemberFromSession(sessionID)
}

func (s SessionStore) DeleteSession(sessionID string) {
	err := s.provider.DeleteSession(sessionID)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Session deletion failed", slog.String("error", err.Error()))
	}
}
