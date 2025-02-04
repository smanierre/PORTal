package backend

import (
	"PORTal/types"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

func (b Backend) Login(username, password string) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to login member", slog.String("username", username))
	member, err := b.memberProvider.GetMember(username, ByUsername)
	if err != nil || member.Disabled {
		return types.Member{}, ErrAuthenticationFailed
	}
	if err = bcrypt.CompareHashAndPassword([]byte(member.Hash), []byte(password)); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Password validation failed", slog.String("error", err.Error()))
		return types.Member{}, ErrAuthenticationFailed
	}
	return member, nil
}

func (b Backend) CreateSession(memberID, userAgent, ipAddress string) (string, time.Time) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to create session", slog.String("memberID", memberID))
	sessionId := uuid.NewString()
	expiration := b.clock.Now().Add(time.Duration(b.config.SessionTimeout) * time.Hour).UTC()
	err := b.authenticationProvider.CreateSession(memberID, sessionId, userAgent, ipAddress, expiration)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Session creation failed")
		return "", time.Now()
	}
	return sessionId, expiration
}

func (b Backend) ValidateSession(sessionID, userAgent, ipAddress string) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to validate session")
	session, err := b.authenticationProvider.GetSession(sessionID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Session validation failed due to database")
		return types.Member{}, fmt.Errorf("%w: %s", ErrSessionValidationFailed, err.Error())
	}
	if userAgent != session.UserAgent {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Session validation failed due to user agent mismatch")
		return types.Member{}, ErrSessionValidationFailed
	}
	if ipAddress != session.IpAddress {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Session validation failed due to IP mismatch")
		return types.Member{}, ErrSessionValidationFailed
	}
	return b.memberProvider.GetMemberFromSession(sessionID)
}

func (b Backend) DeleteSession(sessionID string) {
	err := b.authenticationProvider.DeleteSession(sessionID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Session deletion failed", slog.String("error", err.Error()))
	}
}
