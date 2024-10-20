package api

import (
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
)

func (s Server) frontendHandler(source string) func(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(source))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fullPath := source + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					s.logger.LogAttrs(r.Context(), slog.LevelError, "Error when statting file", slog.String("error", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				r.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, r)
	}
}
