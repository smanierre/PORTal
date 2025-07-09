package rendered

import (
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

func RenderErrorPage(w http.ResponseWriter, r *http.Request, logger *slog.Logger, statusCode int, template templ.Component) {
	w.Header().Set("Hx-Retarget", "body")
	w.Header().Set("Hx-Swap", "innerHTML")
	w.WriteHeader(statusCode)
	err := template.Render(r.Context(), w)
	if err != nil {
		logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering GenericISE page", slog.String("err", err.Error()))
	}
}
