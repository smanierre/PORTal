package memberqualificationstore

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (m MemberQualificationStore) GetRecurringMemberRequirement(memberID, requirementID string) (types.RecurringMemberRequirement, error) {
	return m.memberRequirementProvider.GetRecurringMemberRequirement(memberID, requirementID)
}

func (m MemberQualificationStore) GetRecurringMemberRequirements(memberID string, qualificationID *string) ([]types.RecurringMemberRequirement, error) {
	if qualificationID == nil {
		return m.memberRequirementProvider.GetRecurringMemberRequirements(memberID)
	}
	return m.memberRequirementProvider.GetRecurringMemberRequirementsForQualification(memberID, *qualificationID)
}

func (m MemberQualificationStore) CompleteRecurringMemberRequirement(memberID, requirementID, completedBy string) error {
	memberRequirement, err := m.GetRecurringMemberRequirement(memberID, requirementID)
	if err != nil {
		return err
	}
	for _, completion := range memberRequirement.CompletionHistory {
		if m.clock.Now().Before(completion.CompletionDate) {
			m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Requirement was completed after provided date, provide a valid date")
			return backend.ErrInvalidRequirementCompletionDate
		}
	}
	return m.memberRequirementProvider.CompleteRecurringMemberRequirement(uuid.NewString(), memberRequirement.ID, completedBy, m.clock.Now())
}
