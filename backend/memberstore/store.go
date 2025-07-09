package memberstore

import (
	"PORTal/backend"
	"PORTal/types"
	"log/slog"
	"time"
)

type MemberStore struct {
	provider MemberProvider
	hashCost int
	logger   *slog.Logger
}

func New(provider MemberProvider, hashCost int, logger *slog.Logger) MemberStore {
	logger = logger.With(slog.String("source", "memberStore"))
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
	GetSubordinates(memberID string) []types.Member
	RemoveSubordinates(memberID string) error
	UpdateMember(member types.Member) error
	DisableMember(id string) error
	EnableMember(id string) error
	AssignMemberQualification(memberID, qualificationID string, dateAssigned time.Time, assignedBy string) error
	GetMemberQualification(memberID, qualificationID string) (types.MemberQualification, error)
	GetQualificationsForMember(memberID string) ([]types.MemberQualification, error)
	GetAllMemberQualifications() ([]types.MemberQualification, error)
	RemoveMemberQualification(memberID, qualificationID string) error
}
