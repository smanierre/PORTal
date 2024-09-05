package api

import (
	"PORTal/types"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

func (s Server) login(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Deserializing body into types.Credentials")
	var creds types.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error deserializing credentials from client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}
	member, err := s.backend.Login(creds.Username, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session, err := s.backend.AddSession(member.ID, r.UserAgent())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Skipping setting session cookie, but still logging in")
	} else {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Setting session cookie on response")
		http.SetCookie(w, &http.Cookie{
			Name:     "session-id",
			Value:    session.SessionID,
			Path:     "/api",
			Domain:   os.Getenv("DOMAIN"),
			Expires:  session.Expires,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending member back to client", slog.Any("member", member))
	err = json.NewEncoder(w).Encode(member.ToApiMember())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error serializing member to response body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) validateSession(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Checking for session-id cookie")
	cookie, err := r.Cookie("session-id")
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get session-id cookie from request", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Checking for member ID in body")
	var id types.IdJson
	err = json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error deserializing IdJson into struct", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.backend.ValidateSession(cookie.Value, id.ID, r.UserAgent())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
