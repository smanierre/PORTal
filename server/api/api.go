//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/config.yaml ../../api.yaml
package api

import (
	"PORTal/server/stores"
	"PORTal/types"
	"log/slog"
)

type api struct {
	logger                   *slog.Logger
	memberStore              stores.MemberStore
	qualificationStore       stores.QualificationStore
	memberQualificationStore stores.MemberQualificationStore
	sessionStore             stores.SessionStore
	domain                   string
	organization             string
	dev                      bool
	ranks                    types.RankMap
}
