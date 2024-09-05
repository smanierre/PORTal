package sqlite

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func (p Provider) AddSession(sessionID, memberID, userAgent string, expiration time.Time) (types.Session, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to add session for member",
		slog.String("member_id", memberID),
		slog.String("user_agent", userAgent),
		slog.String("expiration", expiration.String()),
	)
	tx, err := p.Db.Begin()
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error creating transaction when trying to add session", slog.String("error", err.Error()))
		return types.Session{}, err
	}
	_, err = tx.Exec(insertSessionQuery, sessionID, expiration, userAgent)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting session into database, rolling back", slog.String("error", err.Error()))
		err = tx.Rollback()
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return types.Session{}, err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Transaction rolled back successfully")
		return types.Session{}, err
	}
	_, err = tx.Exec(insertMemberSessionQuery, memberID, sessionID)
	if err != nil {
		var causeError error
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Foreign key violation, must be user or previous exec would have failed")
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Setting cause error to types.ErrMemberNotFound")
			causeError = backend.ErrMemberNotFound
		}
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting member_session into database, rolling back", slog.String("error", err.Error()))
		err = tx.Rollback()
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return types.Session{}, fmt.Errorf("%w: original error: %w", err, causeError)
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Transaction rolled back successfully")
		return types.Session{}, causeError
	}
	err = tx.Commit()
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error committing transaction", slog.String("error", err.Error()))
		err = tx.Rollback()
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return types.Session{}, err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Transaction rolled back successfully")
		return types.Session{}, err
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully created new session for member", slog.String("member_id", memberID))
	return types.Session{
		SessionID: sessionID,
		UserAgent: userAgent,
		Expires:   expiration,
	}, nil
}

func (p Provider) GetSession(id string) (types.Session, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting session from database", slog.String("session_id", id))
	var s types.Session
	row := p.Db.QueryRow(getSessionQuery, id)
	err := row.Scan(&s.SessionID, &s.Expires, &s.UserAgent)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting session from database", slog.String("error", err.Error()))
		return types.Session{}, err
	}
	return s, nil
}

func (p Provider) DeleteSession(id string) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting session from database", slog.String("session_id", id))
	_, err := p.Db.Exec(deleteSessionQuery, id)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting session from database", slog.String("error", err.Error()))
	}
}

func (p Provider) CheckForMemberSession(memberID, sessionID string) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for session for user",
		slog.String("session_id", sessionID),
		slog.String("user_id", memberID),
	)
	row := p.Db.QueryRow(getMemberSessionQuery, memberID, sessionID)
	var unused, unused2 string
	if err := row.Scan(&unused, &unused2); err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "No session found for user")
		return err
	}
	return nil
}
