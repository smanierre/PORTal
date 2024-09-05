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

const SessionExpiration = 24 * 7 * time.Hour

func (b Backend) AddSession(memberID, userAgent string) (types.Session, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating new session id")
	var missingArgs []string
	if memberID == "" {
		missingArgs = append(missingArgs, "memberId")
	}
	if userAgent == "" {
		missingArgs = append(missingArgs, "userAgent")
	}
	if len(missingArgs) > 0 {
		return types.Session{}, fmt.Errorf("%w: %s", ErrMissingArgs, missingArgs)
	}
	return b.authenticationProvider.AddSession(uuid.NewString(), memberID, userAgent, b.clock.Now().Add(SessionExpiration))
}

func (b Backend) ValidateSession(sessionID, memberID, userAgent string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Validating session",
		slog.String("session_id", sessionID),
		slog.String("member_id", memberID),
		slog.String("user_agent", userAgent),
	)
	err := b.authenticationProvider.CheckForMemberSession(memberID, sessionID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Unable to find session for given member")
		return ErrSessionValidationFailed
	}
	session, err := b.authenticationProvider.GetSession(sessionID)
	if err != nil {
		return ErrSessionValidationFailed
	}
	if session.UserAgent != userAgent {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided user agent doesn't match existing session, deleting session")
		b.authenticationProvider.DeleteSession(sessionID)
		return ErrSessionValidationFailed
	}
	if b.clock.Now().After(session.Expires) {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Session is expired, deleting session")
		b.authenticationProvider.DeleteSession(sessionID)
		return ErrSessionValidationFailed
	}
	return nil
}

func (b Backend) Login(username, password string) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to login member", slog.String("username", username))
	member, err := b.memberProvider.GetMember(username, ByUsername)
	if err != nil {
		return types.Member{}, ErrAuthenticationFailed
	}
	if err = bcrypt.CompareHashAndPassword([]byte(member.Hash), []byte(password)); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Password validation failed")
		return types.Member{}, ErrAuthenticationFailed
	}
	return member, nil
}
