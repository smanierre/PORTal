package server

import (
	"PORTal/templates"
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
		err := s.templateRepo.RenderFragment(w, "dashboard", "content", nil)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard template", slog.String("error", err.Error()))
			return
		}
	} else {
		err := s.templateRepo.Render(w, "dashboard", &templates.TplData{
			NavData: templates.NavData{
				Show:         true,
				OobSwap:      false,
				Member:       member,
				Subordinates: false,
			},
			ContentData: nil})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard page", slog.String("error", err.Error()))
		}
	}
}
