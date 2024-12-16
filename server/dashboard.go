package server

import (
	"PORTal/templates"
	"PORTal/types"
	"fmt"
	"log"
	"net/http"
)

func (s Server) DashboardGetHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Context().Value(MemberContextKey)
	member, ok := m.(types.Member)
	if !ok {
		log.Println("hmmmm")
	}
	if checkHTMXRequest(r) {
		s.templateRepo.RenderFragment(w, "dashboard", "content", nil)
	} else {
		s.templateRepo.Render(w, "dashboard", &templates.TplData{
			NavData: templates.NavData{
				Show:         true,
				DisplayName:  fmt.Sprintf("%s %s %s", member.GetRank(s.config.Service), member.FirstName, member.LastName),
				OobSwap:      false,
				Member:       member,
				Subordinates: false,
			},
			ContentData: nil})
	}
}
