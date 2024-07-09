package sqlite

import (
	"PORTal/types"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net"
	"strings"
	"time"
)

func (b *Backend) AddSession(ipAddress, memberID string) (string, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to add session for member", slog.String("member_id", memberID), slog.String("ip_address", ipAddress))
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "validating IP address")
	if ip := net.ParseIP(ipAddress); ip == nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Invalid IP address supplied by client", slog.String("ip_address", ipAddress))
		return "", types.ErrInvalidIP
	}

	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding session to database", slog.String("ip_address", ipAddress))
	id := uuid.NewString()
	expirationTime := time.Now().Add(24 * 7 * time.Hour)
	tx, err := b.Db.Begin()
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error creating transaction when trying to add session", slog.String("error", err.Error()))
		return "", err
	}
	_, err = tx.Exec(insertSessionQuery, id, expirationTime, ipAddress)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting session into database, rolling back", slog.String("error", err.Error()))
		err = tx.Rollback()
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return "", err
		}
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Transaction rolled back successfully")
		return "", err
	}
	_, err = tx.Exec(insertMemberSessionQuery, memberID, id)
	if err != nil {
		var causeError error
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Foreign key violation, must be user or previous exec would have failed")
			b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Setting cause error to types.ErrMemberNotFound")
			causeError = types.ErrMemberNotFound
		}
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting member_session into database, rolling back", slog.String("error", err.Error()))
		err = tx.Rollback()
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return "", fmt.Errorf("%w: original error: %w", err, causeError)
		}
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Transaction rolled back successfully")
		return "", causeError
	}
	err = tx.Commit()
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error committing transaction", slog.String("error", err.Error()))
		err = tx.Rollback()
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return "", err
		}
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Transaction rolled back successfully")
		return "", err
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully created new session for member", slog.String("member_id", memberID))
	return id, nil
}

func (b *Backend) ValidateSession(sessionID, memberID, ipAddress string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Validating session",
		slog.String("session_id", sessionID),
		slog.String("member_id", memberID),
		slog.String("ip_address", ipAddress),
	)
	row := b.Db.QueryRow(getSessionForMemberQuery, sessionID, memberID)
	var session types.Session
	err := row.Scan(&session.SessionID, &session.Expires, &session.IPAddress)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "No session found for that session id member id combo")
		return fmt.Errorf("%w: no session found with that sessionID memberID combination", types.ErrSessionValidationFailed)
	} else if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning session into struct", slog.String("error", err.Error()))
		return err
	}
	if session.IPAddress != ipAddress {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Provided IP address doesn't match IP of current session")
		return fmt.Errorf("%w: provided IP address doesn't match existing session", types.ErrSessionValidationFailed)
	}
	return nil
}

func (b *Backend) Login(username, password string) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting current member", slog.String("username", username))
	member, err := b.GetMemberByUsername(username)
	if err != nil {
		return types.Member{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(member.Hash), []byte(password))
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Unable to validate supplied password", slog.String("error", err.Error()))
		return types.Member{}, types.ErrPasswordAuthenticationFailed
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully validated password, returning user")
	return member, nil
}
