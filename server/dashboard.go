package server

import (
	"PORTal/templates"
	"PORTal/templates/pages"
	"PORTal/types"
	"log"
	"log/slog"
	"net/http"
)

func (s Server) DashboardGetHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Context().Value(MemberContextKey)
	member, ok := m.(types.Member)
	if !ok {
		//TODO: FIX THIS
		log.Println("hmmmm")
	}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/dashboard")
		err := pages.Dashboard().Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard template", slog.String("error", err.Error()))
			return
		}
	} else {
		err := templates.Root(templates.NavData{
			Show:         true,
			OobSwap:      false,
			Member:       member,
			Subordinates: false,
		},
			pages.Dashboard()).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard page", slog.String("error", err.Error()))
		}
	}
}
