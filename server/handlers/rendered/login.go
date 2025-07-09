package rendered

import (
	"PORTal/backend"
	"PORTal/server/core"
	"PORTal/templates"
	"PORTal/templates/pages"
	"PORTal/templates/pages/dashboard"
	"PORTal/templates/pages/errorpages"
	"errors"
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

func LoginGetHandler(logger *slog.Logger, c core.Core) http.Handler {
	logger = logger.With(slog.String("source", "LoginGetHandler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Rendering login page")
		loginTpl := pages.Login(c.Organization)
		if r.Header.Get("Hx-Request") != "true" {
			tld := core.TopLevelData{Organization: c.Organization}
			err := templates.Root(tld, loginTpl).Render(r.Context(), w)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login page", slog.String("error", err.Error()))
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
				return
			}
		} else {
			err := loginTpl.Render(r.Context(), w)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login fragment", slog.String("err", err.Error()))
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
			}
		}
	})
}

func LoginPostHandler(logger *slog.Logger, c core.Core) http.Handler {
	logger = logger.With(slog.String("source", "LoginPostHandler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Handling login request")
		data, err := c.Login(r)
		if err != nil && errors.Is(err, backend.ErrAuthenticationFailed) {
			RenderLoginError(w, r, logger, http.StatusUnauthorized, pages.LoginAuthError())
			return
		} else if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
			RenderLoginError(w, r, logger, http.StatusBadRequest, pages.LoginCustomError(err.Error()))
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Successfully logged in member")
		w.Header().Set("HX-Push-URL", "/dashboard")
		http.SetCookie(w, c.MakeSessionCookie(data.SessionID, data.SessionExpiration))
		dashboardData, err := c.Dashboard(r)
		if err != nil {
			RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
			return
		}
		content := dashboard.Dashboard(dashboardData)
		tld, err := c.TopLevel(r, &data.Member)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error getting top level data", slog.String("error", err.Error()))
			RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
			return
		}
		err = templates.Root(tld, content).Render(r.Context(), w)
	})
}

func RenderLoginError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, statusCode int, template templ.Component) {
	w.Header().Set("HX-Retarget", "#loginError")
	w.Header().Set("HX-Reswap", "outerHTML")
	w.WriteHeader(statusCode)
	err := template.Render(r.Context(), w)
	if err != nil {
		logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login error fragment", slog.String("err", err.Error()))
	}
}
