package server

import (
	"PORTal/backend"
	"PORTal/templates"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

func (s Server) RootGetHandler(w http.ResponseWriter, r *http.Request) {
	hx := checkHTMXRequest(r)
	if !hx {
		err := s.templateRepo.Render(w, "login", &templates.TplData{
			NavData:     templates.NavData{},
			ContentData: templates.LoginData{Organization: s.config.Organization},
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login template: %s", slog.String("error", err.Error()))
		}
	} else {
		err := s.templateRepo.RenderFragment(w, "login", "content", nil)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login fragment: %s", slog.String("error", err.Error()))
		}
	}
}

func (s Server) RootPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
	}
	m, err := s.backend.Login(r.FormValue("username"), r.FormValue("password"))
	if errors.Is(err, backend.ErrAuthenticationFailed) {
		w.Header().Set("HX-Retarget", "#loginError")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.WriteHeader(http.StatusUnauthorized)
		err = s.templateRepo.RenderFragment(w, "login", "auth_error", nil)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login error fragment", slog.String("error", err.Error()))
		}
		return
	}
	if err != nil {
		w.Header().Set("HX-Retarget", "#loginError")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.WriteHeader(http.StatusInternalServerError)
		err = s.templateRepo.RenderFragment(w, "login", "internal_error", nil)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering ISE error fragment", slog.String("error", err.Error()))
		}
		return
	}
	w.Header().Set("HX-Push-URL", "/dashboard")
	err = s.templateRepo.RenderFragment(w, "dashboard", "content", nil)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard fragment", slog.String("error", err.Error()))
	}
	subordinates, err := s.backend.GetSubordinates(m.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering generic ise error fragment", slog.String("error", err.Error()))
		}
	}
	err = s.templateRepo.RenderFragment(w, "nav", "nav", templates.NavData{
		Show:         true,
		DisplayName:  fmt.Sprintf("%s %s %s", m.GetRank(s.config.Service), m.FirstName, m.LastName),
		Admin:        m.Admin,
		OobSwap:      true,
		Member:       m,
		Subordinates: len(subordinates) > 0,
	})
}
