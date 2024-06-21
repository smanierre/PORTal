package api_test

import (
	"PORTal/types"
)

func newMockBackend() *mockBackend {
	return &mockBackend{
		addMemberOverride:                 func(m types.Member) error { return nil },
		getMemberOverride:                 func(id string) (types.Member, error) { return types.Member{}, nil },
		getAllMembersOverride:             func() ([]types.Member, error) { return []types.Member{}, nil },
		updateMemberOverride:              func(m types.Member) error { return nil },
		deleteMemberOverride:              func(id string) error { return nil },
		addQualificationOverride:          func(q types.Qualification) error { return nil },
		getQualificationOverride:          func(id string) (types.Qualification, error) { return types.Qualification{}, nil },
		getAllQualificationsOverride:      func() ([]types.Qualification, error) { return []types.Qualification{}, nil },
		updateQualificationOverride:       func(q types.Qualification) error { return nil },
		getRequirementOverride:            func(id string) (types.Requirement, error) { return types.Requirement{}, nil },
		deleteQualificationOverride:       func(id string) error { return nil },
		addRequirementOverride:            func(r types.Requirement) error { return nil },
		getAllRequirementsOverride:        func() ([]types.Requirement, error) { return []types.Requirement{}, nil },
		updateRequirementOverride:         func(r types.Requirement) error { return nil },
		deleteRequirementOverride:         func(id string) error { return nil },
		addMemberQualificationOverride:    func(qualID, memberID string) error { return nil },
		getMemberQualificationsOverride:   func(memberID string) ([]types.MemberQualification, error) { return []types.MemberQualification{}, nil },
		updateMemberQualificationOverride: func(mq types.MemberQualification) error { return nil },
		deleteMemberQualificationOverride: func(qualID, memberID string) error { return nil },
	}
}

type mockBackend struct {
	addMemberOverride     func(m types.Member) error
	getMemberOverride     func(id string) (types.Member, error)
	getAllMembersOverride func() ([]types.Member, error)
	updateMemberOverride  func(m types.Member) error
	deleteMemberOverride  func(id string) error

	addQualificationOverride     func(q types.Qualification) error
	getQualificationOverride     func(id string) (types.Qualification, error)
	getAllQualificationsOverride func() ([]types.Qualification, error)
	updateQualificationOverride  func(q types.Qualification) error
	deleteQualificationOverride  func(id string) error

	addRequirementOverride     func(r types.Requirement) error
	getRequirementOverride     func(id string) (types.Requirement, error)
	getAllRequirementsOverride func() ([]types.Requirement, error)
	updateRequirementOverride  func(r types.Requirement) error
	deleteRequirementOverride  func(id string) error

	addMemberQualificationOverride    func(qualID, memberID string) error
	getMemberQualificationsOverride   func(memberID string) ([]types.MemberQualification, error)
	updateMemberQualificationOverride func(mq types.MemberQualification) error
	deleteMemberQualificationOverride func(qualID, memberID string) error
}

func (m *mockBackend) AddMember(me types.Member) error {
	return m.addMemberOverride(me)
}

func (m *mockBackend) GetMember(id string) (types.Member, error) {
	return m.getMemberOverride(id)
}

func (m *mockBackend) GetAllMembers() ([]types.Member, error) {
	return m.getAllMembersOverride()
}

func (m *mockBackend) UpdateMember(me types.Member) error {
	return m.updateMemberOverride(me)
}

func (m *mockBackend) DeleteMember(id string) error {
	return m.deleteMemberOverride(id)
}

func (m *mockBackend) AddQualification(q types.Qualification) error {
	return m.addQualificationOverride(q)
}

func (m *mockBackend) GetQualification(id string) (types.Qualification, error) {
	return m.getQualificationOverride(id)
}

func (m *mockBackend) GetAllQualifications() ([]types.Qualification, error) {
	return m.getAllQualificationsOverride()
}

func (m *mockBackend) UpdateQualification(q types.Qualification) error {
	return m.updateQualificationOverride(q)
}

func (m *mockBackend) DeleteQualification(id string) error {
	return m.deleteQualificationOverride(id)
}

func (m *mockBackend) AddRequirement(r types.Requirement) error {
	return m.addRequirementOverride(r)
}

func (m *mockBackend) GetRequirement(id string) (types.Requirement, error) {
	return m.getRequirementOverride(id)
}

func (m *mockBackend) GetAllRequirements() ([]types.Requirement, error) {
	return m.getAllRequirementsOverride()
}

func (m *mockBackend) UpdateRequirement(r types.Requirement) error {
	return m.updateRequirementOverride(r)
}

func (m *mockBackend) DeleteRequirement(id string) error {
	return m.deleteRequirementOverride(id)
}

func (m *mockBackend) AddMemberQualification(qualID, memberID string) error {
	return m.addMemberQualificationOverride(qualID, memberID)
}

func (m *mockBackend) GetMemberQualifications(memberID string) ([]types.MemberQualification, error) {
	return m.getMemberQualificationsOverride(memberID)
}

func (m *mockBackend) UpdateMemberQualification(mq types.MemberQualification) error {
	return m.updateMemberQualificationOverride(mq)
}

func (m *mockBackend) DeleteMemberQualification(qualID, memberID string) error {
	return m.deleteMemberQualificationOverride(qualID, memberID)
}
