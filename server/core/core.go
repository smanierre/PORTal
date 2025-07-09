package core

import (
	"PORTal/server/stores"
	"PORTal/types"
	"log/slog"
)

type Core struct {
	logger                   *slog.Logger
	memberStore              stores.MemberStore
	qualificationStore       stores.QualificationStore
	memberQualificationStore stores.MemberQualificationStore
	sessionStore             stores.SessionStore
	Ranks                    types.RankMap
	Organization             string
	Service                  string
	domain                   string
}

func New(
	logger *slog.Logger,
	memberStore stores.MemberStore,
	qualificationStore stores.QualificationStore,
	memberQualificationStore stores.MemberQualificationStore,
	sessionStore stores.SessionStore,
	ranks types.RankMap,
	organization string,
	service string,
	domain string,
) Core {
	return Core{
		logger:                   logger,
		memberStore:              memberStore,
		qualificationStore:       qualificationStore,
		memberQualificationStore: memberQualificationStore,
		sessionStore:             sessionStore,
		Ranks:                    ranks,
		Organization:             organization,
		Service:                  service,
		domain:                   domain,
	}
}
