package testutils

import (
	"PORTal/types"
	"time"
)

type MockSessionStore struct {
	CreateSessionOverride   func(memberID, userAgent, ipAddress string) (string, time.Time)
	ValidateSessionOverride func(sessionID, userAgent, ipAddress string) (types.Member, error)
	DeleteSessionOverride   func(sessionID string)
}

func (m *MockSessionStore) CreateSession(memberID, userAgent, ipAddress string) (string, time.Time) {
	return m.CreateSessionOverride(memberID, userAgent, ipAddress)
}
func (m *MockSessionStore) ValidateSession(sessionID, userAgent, ipAddress string) (types.Member, error) {
	return m.ValidateSessionOverride(sessionID, userAgent, ipAddress)
}
func (m *MockSessionStore) DeleteSession(sessionID string) {
	m.DeleteSessionOverride(sessionID)
}

type MockMemberStore struct {
	LoginOverride                   func(username, password string) (types.Member, error)
	AddMemberOverride               func(mem types.Member) (types.Member, error)
	GetMemberOverride               func(identifier string) (types.Member, error)
	GetPotentialSupervisorsOverride func(mem types.Member, grade types.Grade) ([]types.Member, error)
	GetMemberFromSessionOverride    func(sessionID string) (types.Member, error)
	GetDisabledMemberOverride       func(identifier string) (types.Member, error)
	GetAllMembersOverride           func() ([]types.Member, error)
	GetDisabledMembersOverride      func() ([]types.Member, error)
	GetSubordinatesOverride         func(memberID string) []types.Member
	UpdateMemberOverride            func(mem types.Member) (types.Member, error)
	DeleteMemberOverride            func(id string) error
	DisableMemberOverride           func(id string) error
	EnableMemberOverride            func(id string) error
}

func (m *MockMemberStore) Login(username, password string) (types.Member, error) {
	return m.LoginOverride(username, password)
}
func (m *MockMemberStore) AddMember(mem types.Member) (types.Member, error) {
	return m.AddMemberOverride(mem)
}
func (m *MockMemberStore) GetMember(identifier string) (types.Member, error) {
	return m.GetMemberOverride(identifier)
}
func (m *MockMemberStore) GetPotentialSupervisors(mem types.Member, grade types.Grade) ([]types.Member, error) {
	return m.GetPotentialSupervisorsOverride(mem, grade)
}
func (m *MockMemberStore) GetMemberFromSession(sessionID string) (types.Member, error) {
	return m.GetMemberFromSessionOverride(sessionID)
}
func (m *MockMemberStore) GetDisabledMember(identifier string) (types.Member, error) {
	return m.GetDisabledMemberOverride(identifier)
}
func (m *MockMemberStore) GetAllMembers() ([]types.Member, error) {
	return m.GetAllMembersOverride()
}
func (m *MockMemberStore) GetDisabledMembers() ([]types.Member, error) {
	return m.GetDisabledMembersOverride()
}
func (m *MockMemberStore) GetSubordinates(memberID string) []types.Member {
	return m.GetSubordinatesOverride(memberID)
}
func (m *MockMemberStore) UpdateMember(mem types.Member) (types.Member, error) {
	return m.UpdateMemberOverride(mem)
}
func (m *MockMemberStore) DeleteMember(id string) error {
	return m.DeleteMemberOverride(id)
}
func (m *MockMemberStore) DisableMember(id string) error {
	return m.DisableMemberOverride(id)
}
func (m *MockMemberStore) EnableMember(id string) error {
	return m.EnableMemberOverride(id)
}

type MockQualificationStore struct {
	AddQualificationOverride                 func(q types.Qualification) (types.Qualification, error)
	AssignRequirementToQualificationOverride func(qualificationID, requirementID string, initial bool) error
	GetQualificationOverride                 func(identifier string) (types.Qualification, error)
	GetAllQualificationsOverride             func() ([]types.Qualification, error)
	UpdateQualificationOverride              func(q types.Qualification) (types.Qualification, error)
	DeleteQualificationOverride              func(identifier string) error
	AddRequirementOverride                   func(r types.Requirement) (types.Requirement, error)
	GetRequirementOverride                   func(identifier string) (types.Requirement, error)
	GetAllRequirementsOverride               func() ([]types.Requirement, error)
	UpdateRequirementOverride                func(r types.Requirement) (types.Requirement, error)
	DeleteRequirementOverride                func(identifier string) error
}

func (m MockQualificationStore) AddQualification(q types.Qualification) (types.Qualification, error) {
	return m.AddQualificationOverride(q)
}
func (m MockQualificationStore) AssignRequirementToQualification(qualificationID, requirementID string, initial bool) error {
	return m.AssignRequirementToQualificationOverride(qualificationID, requirementID, initial)
}
func (m MockQualificationStore) GetQualification(id string) (types.Qualification, error) {
	return m.GetQualificationOverride(id)
}
func (m MockQualificationStore) GetAllQualifications() ([]types.Qualification, error) {
	return m.GetAllQualificationsOverride()
}
func (m MockQualificationStore) UpdateQualification(q types.Qualification) (types.Qualification, error) {
	return m.UpdateQualificationOverride(q)
}
func (m MockQualificationStore) DeleteQualification(id string) error {
	return m.DeleteQualificationOverride(id)
}
func (m MockQualificationStore) AddRequirement(r types.Requirement) (types.Requirement, error) {
	return m.AddRequirementOverride(r)
}
func (m MockQualificationStore) GetRequirement(id string) (types.Requirement, error) {
	return m.GetRequirementOverride(id)
}
func (m MockQualificationStore) GetAllRequirements() ([]types.Requirement, error) {
	return m.GetAllRequirementsOverride()
}
func (m MockQualificationStore) UpdateRequirement(r types.Requirement) (types.Requirement, error) {
	return m.UpdateRequirementOverride(r)
}
func (m MockQualificationStore) DeleteRequirement(id string) error {
	return m.DeleteRequirementOverride(id)
}
