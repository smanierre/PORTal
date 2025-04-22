package stores

import (
	"PORTal/types"
	"time"
)

type MemberStore interface {
	Login(username, password string) (types.Member, error)
	AddMember(m types.Member) (types.Member, error)
	GetMember(identifier string) (types.Member, error)
	GetPotentialSupervisors(m types.Member, grade types.Grade) ([]types.Member, error)
	GetMemberFromSession(sessionID string) (types.Member, error)
	GetDisabledMember(identifier string) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	GetDisabledMembers() ([]types.Member, error)
	GetSubordinates(memberID string) []types.Member
	UpdateMember(m types.Member) (types.Member, error)
	DeleteMember(id string) error
	DisableMember(id string) error
	EnableMember(id string) error
}

type QualificationStore interface {
	AddQualification(q types.Qualification) (types.Qualification, error)
	AssignRequirementToQualification(qualificationID, requirementID string, initial bool) error
	GetQualification(id string) (types.Qualification, error)
	GetAllQualifications() ([]types.Qualification, error)
	UpdateQualification(q types.Qualification) (types.Qualification, error)
	DeleteQualification(id string) error
	AddRequirement(r types.Requirement) (types.Requirement, error)
	GetRequirement(id string) (types.Requirement, error)
	GetAllRequirements() ([]types.Requirement, error)
	UpdateRequirement(r types.Requirement) (types.Requirement, error)
	DeleteRequirement(id string) error
}

type MemberQualificationStore interface {
	AssignMemberQualification(memberID, qualID string) error
	GetMemberQualification(memberID string, qualificationID string) (types.Qualification, error)
	GetMemberQualifications(memberID string) ([]types.Qualification, error)
	RemoveMemberQualification(memberID, qualificationID string) error
}

type SessionStore interface {
	CreateSession(memberID, userAgent, ipAddress string) (string, time.Time)
	ValidateSession(sessionID, userAgent, ipAddress string) (types.Member, error)
	DeleteSession(sessionID string)
}
