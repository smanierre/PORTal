package api

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

func (a api) ValidateSessionCookie(r *http.Request, w http.ResponseWriter) (types.Member, error) {
	// Check for session
	sessionCookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, setting cookie to empty cookie")
		sessionCookie = &http.Cookie{}
	}
	m, err := a.sessionStore.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
	if errors.Is(err, backend.ErrSessionValidationFailed) {
		a.logger.LogAttrs(r.Context(), slog.LevelDebug, "Session validation failed, removing cookie ")
		a.RemoveCookie(w, SessionCookieName)
		return types.Member{}, err
	}
	return m, nil
}

func (a api) MakeCookie(name, value string, expiration time.Time) *http.Cookie {
	secure := !a.dev
	var sameSite http.SameSite
	if !secure {
		sameSite = http.SameSiteLaxMode
	} else {
		sameSite = http.SameSiteStrictMode
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   a.domain,
		Expires:  expiration,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	}
}

func (a api) MakeSessionCookie(value string, expiration time.Time) *http.Cookie {
	return a.MakeCookie(SessionCookieName, value, expiration)
}

func (a api) RemoveCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  a.domain,
		Expires: time.Now(),
		Path:    "/",
	})
}

func (a api) RemoveSessionCookie(w http.ResponseWriter) {
	a.RemoveCookie(w, SessionCookieName)
}

func (a api) getSessionId(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
