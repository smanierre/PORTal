package api_test

import (
	"PORTal/api"
	"PORTal/types"
)

var _ api.Backend = (*mockBackend)(nil)

func newMockBackend() *mockBackend {
	return &mockBackend{
		addMemberOverride:            func(m types.Member) (types.Member, error) { return types.Member{}, nil },
		getMemberOverride:            func(id string) (types.Member, error) { return types.Member{}, nil },
		getAllMembersOverride:        func() ([]types.Member, error) { return []types.Member{}, nil },
		getSubordinatesOverride:      func(id string) ([]types.Member, error) { return nil, nil },
		updateMemberOverride:         func(m types.Member) (types.Member, error) { return types.Member{}, nil },
		deleteMemberOverride:         func(id string) error { return nil },
		addQualificationOverride:     func(q types.Qualification) (types.Qualification, error) { return types.Qualification{}, nil },
		getQualificationOverride:     func(id string) (types.Qualification, error) { return types.Qualification{}, nil },
		getAllQualificationsOverride: func() ([]types.Qualification, error) { return []types.Qualification{}, nil },
		updateQualificationOverride: func(q types.Qualification, forceUpdateExpiration bool) (types.Qualification, error) {
			return types.Qualification{}, nil
		},
		getRequirementOverride:            func(id string) (types.Requirement, error) { return types.Requirement{}, nil },
		deleteQualificationOverride:       func(id string) error { return nil },
		addRequirementOverride:            func(r types.Requirement) (types.Requirement, error) { return types.Requirement{}, nil },
		getAllRequirementsOverride:        func() ([]types.Requirement, error) { return []types.Requirement{}, nil },
		updateRequirementOverride:         func(r types.Requirement) (types.Requirement, error) { return types.Requirement{}, nil },
		deleteRequirementOverride:         func(id string) error { return nil },
		assignMemberQualificationOverride: func(qualID, memberID string) error { return nil },
		getMemberQualificationOverride:    func(memberID, qualID string) (types.Qualification, error) { return types.Qualification{}, nil },
		getMemberQualificationsOverride:   func(memberID string) ([]types.Qualification, error) { return nil, nil },
		removeMemberQualificationOverride: func(memberID, qualId string) error { return nil },
		addReferenceOverride:              func(r types.Reference) (types.Reference, error) { return types.Reference{}, nil },
		getReferenceOverride:              func(id string) (types.Reference, error) { return types.Reference{}, nil },
		getReferencesOverride:             func() ([]types.Reference, error) { return nil, nil },
		updateReferenceOverride:           func(r types.Reference, overrideNoVolume bool) (types.Reference, error) { return types.Reference{}, nil },
		deleteReferenceOverride:           func(id string) error { return nil },
		addSessionOverride:                func(memberID, userAgent string) (types.Session, error) { return types.Session{}, nil },
		validateSessionOverride:           func(sessionID, memberID, ipAddress string) error { return nil },
		loginOverride:                     func(username, password string) (types.Member, error) { return types.Member{}, nil },
	}
}

type mockBackend struct {
	addMemberOverride       func(m types.Member) (types.Member, error)
	getMemberOverride       func(id string) (types.Member, error)
	getAllMembersOverride   func() ([]types.Member, error)
	getSubordinatesOverride func(id string) ([]types.Member, error)
	updateMemberOverride    func(m types.Member) (types.Member, error)
	deleteMemberOverride    func(id string) error

	addQualificationOverride     func(q types.Qualification) (types.Qualification, error)
	getQualificationOverride     func(id string) (types.Qualification, error)
	getAllQualificationsOverride func() ([]types.Qualification, error)
	updateQualificationOverride  func(q types.Qualification, forceUpdateExpiration bool) (types.Qualification, error)
	deleteQualificationOverride  func(id string) error

	addRequirementOverride     func(r types.Requirement) (types.Requirement, error)
	getRequirementOverride     func(id string) (types.Requirement, error)
	getAllRequirementsOverride func() ([]types.Requirement, error)
	updateRequirementOverride  func(r types.Requirement) (types.Requirement, error)
	deleteRequirementOverride  func(id string) error

	assignMemberQualificationOverride func(qualID, memberID string) error
	getMemberQualificationOverride    func(memberID, qualID string) (types.Qualification, error)
	getMemberQualificationsOverride   func(memberID string) ([]types.Qualification, error)
	removeMemberQualificationOverride func(memberID, qualID string) error

	addReferenceOverride    func(r types.Reference) (types.Reference, error)
	getReferenceOverride    func(id string) (types.Reference, error)
	getReferencesOverride   func() ([]types.Reference, error)
	updateReferenceOverride func(r types.Reference, overrideNoVolume bool) (types.Reference, error)
	deleteReferenceOverride func(id string) error

	addSessionOverride      func(memberID, userAgent string) (types.Session, error)
	validateSessionOverride func(sessionID, memberID, ipAddress string) error
	loginOverride           func(username, password string) (types.Member, error)
}

func (m *mockBackend) AddMember(me types.Member) (types.Member, error) {
	return m.addMemberOverride(me)
}

func (m *mockBackend) GetMember(id string) (types.Member, error) {
	return m.getMemberOverride(id)
}

func (m *mockBackend) GetAllMembers() ([]types.Member, error) {
	return m.getAllMembersOverride()
}

func (m *mockBackend) GetSubordinates(id string) ([]types.Member, error) {
	return m.getSubordinatesOverride(id)
}

func (m *mockBackend) UpdateMember(me types.Member) (types.Member, error) {
	return m.updateMemberOverride(me)
}

func (m *mockBackend) DeleteMember(id string) error {
	return m.deleteMemberOverride(id)
}

func (m *mockBackend) AddQualification(q types.Qualification) (types.Qualification, error) {
	return m.addQualificationOverride(q)
}

func (m *mockBackend) GetQualification(id string) (types.Qualification, error) {
	return m.getQualificationOverride(id)
}

func (m *mockBackend) GetAllQualifications() ([]types.Qualification, error) {
	return m.getAllQualificationsOverride()
}

func (m *mockBackend) UpdateQualification(q types.Qualification, forceUpdateExpiration bool) (types.Qualification, error) {
	return m.updateQualificationOverride(q, forceUpdateExpiration)
}

func (m *mockBackend) DeleteQualification(id string) error {
	return m.deleteQualificationOverride(id)
}

func (m *mockBackend) AddRequirement(r types.Requirement) (types.Requirement, error) {
	return m.addRequirementOverride(r)
}

func (m *mockBackend) GetRequirement(id string) (types.Requirement, error) {
	return m.getRequirementOverride(id)
}

func (m *mockBackend) GetAllRequirements() ([]types.Requirement, error) {
	return m.getAllRequirementsOverride()
}

func (m *mockBackend) UpdateRequirement(r types.Requirement) (types.Requirement, error) {
	return m.updateRequirementOverride(r)
}

func (m *mockBackend) DeleteRequirement(id string) error {
	return m.deleteRequirementOverride(id)
}

func (m *mockBackend) AssignMemberQualification(qualID, memberID string) error {
	return m.assignMemberQualificationOverride(qualID, memberID)
}

func (m *mockBackend) GetMemberQualification(memberID, qualID string) (types.Qualification, error) {
	return m.getMemberQualificationOverride(memberID, qualID)
}

func (m *mockBackend) GetMemberQualifications(memberID string) ([]types.Qualification, error) {
	return m.getMemberQualificationsOverride(memberID)
}

func (m *mockBackend) RemoveMemberQualification(memberID, qualID string) error {
	return m.removeMemberQualificationOverride(memberID, qualID)
}

func (m *mockBackend) AddReference(r types.Reference) (types.Reference, error) {
	return m.addReferenceOverride(r)
}

func (m *mockBackend) GetReference(id string) (types.Reference, error) {
	return m.getReferenceOverride(id)
}

func (m *mockBackend) GetReferences() ([]types.Reference, error) {
	return m.getReferencesOverride()
}

func (m *mockBackend) UpdateReference(r types.Reference, overrideNoVolume bool) (types.Reference, error) {
	return m.updateReferenceOverride(r, overrideNoVolume)
}

func (m *mockBackend) DeleteReference(id string) error {
	return m.deleteReferenceOverride(id)
}

func (m *mockBackend) AddSession(memberID, userAgent string) (types.Session, error) {
	return m.addSessionOverride(memberID, userAgent)
}

func (m *mockBackend) ValidateSession(sessionID, memberID, ipAddress string) error {
	return m.validateSessionOverride(sessionID, memberID, ipAddress)
}

func (m *mockBackend) Login(username, password string) (types.Member, error) {
	return m.loginOverride(username, password)
}
