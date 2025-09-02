package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a api) PostApiLogin(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "Error decoding login request body", slog.String("err", err.Error()))
		_ = json.NewEncoder(w).Encode(Error{Message: "Invalid request body."})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := a.memberStore.Login(request.Username, request.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(Error{Message: "Invalid username or password"})
		return
	}
	sessionID, expiration := a.sessionStore.CreateSession(m.ID, r.UserAgent(), r.RemoteAddr)
	c := a.MakeSessionCookie(sessionID, expiration)
	http.SetCookie(w, c)
	_ = json.NewEncoder(w).Encode(memberToApiMember(m))
}

func (a api) PostApiLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelInfo, "No session cookie found during logout, returning")
		return
	}
	a.sessionStore.DeleteSession(sessionCookie.Value)
	a.RemoveSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}
