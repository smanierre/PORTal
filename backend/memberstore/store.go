package memberstore

import (
	"PORTal/backend"
	"PORTal/types"
	"log/slog"
)

type MemberStore struct {
	provider MemberProvider
	hashCost int
	logger   *slog.Logger
}

func New(provider MemberProvider, hashCost int, logger *slog.Logger) MemberStore {
	logger = logger.With(slog.String("source", "MemberStore"))
	return MemberStore{
		provider: provider,
		hashCost: hashCost,
		logger:   logger,
	}
}

type MemberProvider interface {
	AddMember(m types.Member) error
	GetMember(identifier string, method backend.ProviderMethod) (types.Member, error)
	GetMemberFromSession(sessionID string) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	GetDisabledMembers() ([]types.Member, error)
	GetSubordinates(memberID string) ([]types.Member, error)
	RemoveSubordinates(memberID string) error
	UpdateMember(member types.Member) error
	DeleteMember(identifier string, method backend.ProviderMethod) error
	DisableMember(id string) error
	EnableMember(id string) error
	AssignMemberQualification(memberID, qualificationID string) error
	RemoveMemberQualification(memberID, qualificationID string) error
}
