package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func GetApiSpec(logger *slog.Logger) http.Handler {
	logger = logger.With(slog.String("Source", "GetApiSpec"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spec, err := GetSwagger()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error getting API spec", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(spec)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error encoding API spec", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
