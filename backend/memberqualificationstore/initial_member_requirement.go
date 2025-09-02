package memberqualificationstore

import (
	"PORTal/types"
)

func (m MemberQualificationStore) GetInitialMemberRequirement(memberID, requirementID string) (types.InitialMemberRequirement, error) {
	return m.memberRequirementProvider.GetInitialMemberRequirement(memberID, requirementID)
}

func (m MemberQualificationStore) GetInitialMemberRequirements(memberID string, qualificationID *string) ([]types.InitialMemberRequirement, error) {
	if qualificationID == nil {
		return m.memberRequirementProvider.GetInitialMemberRequirements(memberID)
	}
	return m.memberRequirementProvider.GetInitialMemberRequirementsForQualification(memberID, *qualificationID)
}

func (m MemberQualificationStore) CompleteInitialMemberRequirement(memberID, requirementID, completedBy string) error {
	return m.memberRequirementProvider.CompleteInitialMemberRequirement(m.clock.Now(), completedBy, memberID, requirementID)
}
