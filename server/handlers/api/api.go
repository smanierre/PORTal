package api

import (
	"PORTal/server/stores"
	"log/slog"
)

type api struct {
	logger      *slog.Logger
	memberStore stores.MemberStore
}
