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

type Backend struct {
	memberProvider        MemberProvider
	qualificationProvider QualificationProvider
	requirementProvider   RequirementProvider
	clock                 Clock
	logger                *slog.Logger
	config                Config
}

type MemberProvider interface {
	AddMember(m types.Member) error
	GetMember(identifier string, method ProviderMethod) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	GetSubordinates(memberID string) ([]types.Member, error)
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

type Clock interface {
	Now() time.Time
}

type Config struct {
	DbFile     string `yaml:"DbFile"`
	BcryptCost int    `yaml:"BcryptCost"`
}

type realTime struct{}

func (r realTime) Now() time.Time {
	return time.Now()
}

func New(logger *slog.Logger, memberProvider MemberProvider, qualificationProvider QualificationProvider,
	requirementProvider RequirementProvider, config Config, clock Clock) Backend {
	if clock == nil {
		clock = realTime{}
	}

	logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Using bcrypt cost: %d", config.BcryptCost))
	return Backend{
		memberProvider:        memberProvider,
		qualificationProvider: qualificationProvider,
		requirementProvider:   requirementProvider,
		clock:                 clock,
		logger:                logger,
		config:                config,
	}
}
