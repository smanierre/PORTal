package server

import (
	"PORTal/backend"
	"PORTal/server/stores"
	"PORTal/templates"
	"PORTal/templates/pages"
	"PORTal/templates/pages/dashboard"
	"PORTal/types"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func LoginDirectorHandler(logger *slog.Logger, organization string, ranks types.RankMap, memberStore stores.MemberStore, sessionStore stores.SessionStore) http.Handler {
	logger = logger.With(slog.String("route", "/"))
	loginGetHandler := LoginGetHandler(logger.With(slog.String("route", "GET /")), organization)
	loginPostHandler := LoginPostHandler(logger.With(slog.String("route", "POST /")), memberStore, sessionStore, ranks, organization)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			loginGetHandler.ServeHTTP(w, r)
		} else if r.Method == "POST" && r.URL.Path == "/" {
			loginPostHandler.ServeHTTP(w, r)
		} else {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "Requested path and method not found", slog.String("path", r.URL.Path), slog.String("method", r.Method))
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func LoginGetHandler(logger *slog.Logger, organization string) http.Handler {
	logger = logger.With(slog.String("source", "LoginGetHandler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginTpl := pages.Login(organization)
		if r.Header.Get("Hx-Request") != "true" {
			serveContentAsRoot(loginTpl, logger, w, r)
		} else {
			HandleRenderError(r.Context(), logger, loginTpl.Render(r.Context(), w))
		}
	})
}

func LoginPostHandler(logger *slog.Logger, memberStore stores.MemberStore, sessionStore stores.SessionStore, ranks types.RankMap, organization string) http.Handler {
	logger = logger.With(slog.String("source", "LoginPostHandler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Handling login request")
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
			w.Header().Set("HX-Retarget", "#loginError")
			w.Header().Set("HX-Reswap", "outerHTML")
			HandleRenderError(r.Context(), logger, pages.LoginCustomError(err.Error()).Render(r.Context(), w))
			return
		}
		member, err := memberStore.Login(r.Form.Get("username"), r.Form.Get("password"))
		if err != nil && errors.Is(err, backend.ErrAuthenticationFailed) {
			w.Header().Set("HX-Retarget", "#loginError")
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusUnauthorized)
			HandleRenderError(r.Context(), logger, pages.LoginAuthError().Render(r.Context(), w))
			return
		} else if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
			w.Header().Set("HX-Retarget", "#loginError")
			w.Header().Set("HX-Reswap", "outerHTML")
			HandleRenderError(r.Context(), logger, pages.LoginCustomError(err.Error()).Render(r.Context(), w))
			return
		}
		sessionId, expiration := sessionStore.CreateSession(member.ID, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
		w.Header().Set("HX-Push-URL", "/dashboard")
		http.SetCookie(w, MakeCookie(SessionCookieName, sessionId, expiration))
		var hasSubordinates bool
		subordinates := memberStore.GetSubordinates(member.ID)
		if len(subordinates) > 0 {
			hasSubordinates = true
		}
		displayName := fmt.Sprintf("%s %s %s", ranks[member.Grade], member.FirstName, member.LastName)
		content := dashboard.Dashboard()
		HandleRenderError(r.Context(), logger, templates.Root(true, hasSubordinates, member.Admin, displayName, organization, content).Render(r.Context(), w))
	})
}
