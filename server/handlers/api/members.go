package api

import (
	"PORTal/server/core"
	"PORTal/types"
	"encoding/json"
	"log/slog"
	"net/http"
)

type PotentialSupervisorsResponse struct {
	PotentialSupervisors []types.Member `json:"potential_supervisors"`
}

func PotentialSupervisors(logger *slog.Logger, c core.Core) http.Handler {
	logger = logger.With(slog.String("Source", "API/PotentialSupervisors"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := c.AdminMemberPage(r)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error getting admin page data", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(PotentialSupervisorsResponse{PotentialSupervisors: data.PotentialSupervisors})
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error encoding response data", slog.String("error", err.Error()))
		}
		return
	})
}
