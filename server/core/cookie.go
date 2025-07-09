package core

import (
	"PORTal/backend"
	"PORTal/types"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const SessionCookieName = "session_id"

func (c Core) ValidateSessionCookie(r *http.Request, w http.ResponseWriter) (types.Member, error) {
	// Check for session
	sessionCookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		c.logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, setting cookie to empty cookie")
		sessionCookie = &http.Cookie{}
	}
	m, err := c.sessionStore.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
	if errors.Is(err, backend.ErrSessionValidationFailed) {
		c.logger.LogAttrs(r.Context(), slog.LevelDebug, "Session validation failed, removing cookie ")
		c.RemoveCookie(w, SessionCookieName)
		return types.Member{}, err
	}
	return m, nil
}

func (c Core) MakeCookie(name, value string, expiration time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   c.domain,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (c Core) MakeSessionCookie(value string, expiration time.Time) *http.Cookie {
	return c.MakeCookie(SessionCookieName, value, expiration)
}

func (c Core) RemoveCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  c.domain,
		Expires: time.Now(),
		Path:    "/",
	})
}

func (c Core) RemoveSessionCookie(w http.ResponseWriter) {
	c.RemoveCookie(w, SessionCookieName)
}

func (c Core) getSessionId(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
