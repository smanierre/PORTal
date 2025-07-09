package core

import (
	"PORTal/types"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type LoginData struct {
	Member            types.Member
	SessionID         string
	SessionExpiration time.Time
}

func (c Core) Login(r *http.Request) (LoginData, error) {
	err := r.ParseForm()
	if err != nil {
		c.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
		return LoginData{}, err
	}
	member, err := c.memberStore.Login(r.Form.Get("username"), r.Form.Get("password"))
	if err != nil {
		return LoginData{}, err
	}
	sessionId, expiration := c.sessionStore.CreateSession(member.ID, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])

	return LoginData{
		Member:            member,
		SessionID:         sessionId,
		SessionExpiration: expiration,
	}, nil
}

func (c Core) Logout(r *http.Request) {
	sessionID, err := c.getSessionId(r)
	if err != nil {
		c.logger.LogAttrs(r.Context(), slog.LevelInfo, "Failed to get sessionID, skipping deletion")
		return
	}
	c.sessionStore.DeleteSession(sessionID)
}
