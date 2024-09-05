package backend

import (
	"PORTal/types"
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

type ProviderMethod int

const (
	ById ProviderMethod = iota
	ByUsername
)

const MinimumPwLength = 8

var BcryptCost int

type Backend struct {
	memberProvider         MemberProvider
	qualificationProvider  QualificationProvider
	requirementProvider    RequirementProvider
	authenticationProvider AuthenticationProvider
	clock                  Clock
	logger                 *slog.Logger
}

type MemberProvider interface {
	AddMember(m types.Member) error
	GetMember(identifier string, method ProviderMethod) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	UpdateMember(member types.Member) error
	DeleteMember(identifier string, method ProviderMethod) error
	AssignMemberQualification(memberID, qualificationID string) error
	GetMemberQualification(memberID, qualificationID string) (types.Qualification, error)
	GetMemberQualifications(memberID string) ([]types.Qualification, error)
	RemoveMemberQualification(memberID, qualificationID string) error
}

type QualificationProvider interface {
	AddQualification(q types.Qualification) error
	GetQualification(id string) (types.Qualification, error)
	GetAllQualifications() ([]types.Qualification, error)
	UpdateQualification(q types.Qualification) error
	DeleteQualification(id string) error
}

type RequirementProvider interface {
	AddRequirement(r types.Requirement) error
	GetRequirement(id string) (types.Requirement, error)
	GetAllRequirements() ([]types.Requirement, error)
	GetQualificationIDsForRequirement(requirementID string) ([]string, error)
	UpdateRequirement(r types.Requirement) error
	DeleteRequirement(id string) error
	AddReference(r types.Reference) error
	GetReference(id string) (types.Reference, error)
	GetReferences() ([]types.Reference, error)
	UpdateReference(r types.Reference) error
	DeleteReference(id string) error
}

type AuthenticationProvider interface {
	AddSession(sessionID, memberID, userAgent string, expiration time.Time) (types.Session, error)
	GetSession(id string) (types.Session, error)
	CheckForMemberSession(memberID, sessionID string) error
	DeleteSession(id string)
}

type Clock interface {
	Now() time.Time
}

type Options struct {
	BcryptCost int
	Clock      Clock
}

type realTime struct{}

func (r realTime) Now() time.Time {
	return time.Now()
}

func New(logger *slog.Logger, memberProvider MemberProvider, qualificationProvider QualificationProvider,
	requirementProvider RequirementProvider, authenticationProvider AuthenticationProvider, opts *Options) Backend {
	if opts == nil {
		opts = &Options{}
	}
	if opts.BcryptCost != 0 {
		BcryptCost = opts.BcryptCost
	} else {
		BcryptCost = 16
	}
	if opts.Clock == nil {
		opts.Clock = realTime{}
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Using bcrypt cost: %d", BcryptCost))
	return Backend{
		memberProvider:         memberProvider,
		qualificationProvider:  qualificationProvider,
		requirementProvider:    requirementProvider,
		authenticationProvider: authenticationProvider,
		clock:                  opts.Clock,
		logger:                 logger,
	}
}
