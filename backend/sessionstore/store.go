package sessionstore

import (
	"PORTal/backend"
	"PORTal/server/stores"
	"PORTal/types"
	"log/slog"
	"time"
)

type SessionStore struct {
	memberStore stores.MemberStore
	provider    SessionProvider
	logger      *slog.Logger
	clock       backend.Clock
	timeout     time.Duration
}

func New(provider SessionProvider, memberStore stores.MemberStore, clock backend.Clock, timeout time.Duration, logger *slog.Logger) SessionStore {
	logger = logger.With(slog.String("source", "SessionStore"))
	if clock == nil {
		clock = backend.RealClock{}
	}
	return SessionStore{
		memberStore: memberStore,
		provider:    provider,
		logger:      logger,
		clock:       clock,
		timeout:     timeout,
	}
}

type SessionProvider interface {
	CreateSession(memberID, sessionID, userAgent, ipAddress string, expiration time.Time) error
	GetSession(sessionID string) (types.Session, error)
	DeleteSession(sessionID string) error
}
