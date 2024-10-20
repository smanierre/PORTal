package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func (s Server) login(w http.ResponseWriter, r *http.Request) {
	var res LoginResponse
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Deserializing body into types.Credentials")
	var creds Credentials
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
	res.Member = member.ToApiMember()
	res.Qualifications, err = s.backend.GetMemberQualifications(res.Member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	subordinates, err := s.backend.GetSubordinates(res.Member.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, subordinate := range subordinates {
		res.Subordinates = append(res.Subordinates, subordinate.ToApiMember())
	}
	token, err := createToken(member, s.config.JWTExpiration*time.Hour, []byte(s.config.JWTSecret))
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error creating JWT, still logging in", slog.String("error", err.Error()))
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating identity cookie")
	cookie := &http.Cookie{
		Name:     "identity",
		Value:    token,
		Path:     "/api",
		Domain:   s.config.Domain,
		Expires:  time.Now().Add(s.config.JWTExpiration * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending response back to client", slog.Any("response", res))
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error serializing response to client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) logout(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Clearing identity cookie for member")
	http.SetCookie(w, &http.Cookie{
		Name:    JWTCookieName,
		Path:    "/api",
		Domain:  s.config.Domain,
		Expires: time.Now(),
	})
}

func (s Server) checkAdmin(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Validating member's admin permissions")
	tokenCookie, err := r.Cookie(JWTCookieName)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error getting token cookie", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := validateToken(tokenCookie.Value, s.jwtKeyFunc, s.logger)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	customClaims, ok := token.Claims.(*CustomClaims)
	if !ok {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error casting claims to CustomClaims")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !customClaims.Admin {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "User is not admin", slog.String("member_id", customClaims.Subject))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
