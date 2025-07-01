package memberqualificationstore

import (
	"PORTal/backend"
	"PORTal/backend/memberstore"
	"PORTal/backend/qualificationstore"
	"PORTal/types"
	"log/slog"
	"time"
)

type MemberQualificationStore struct {
	logger                    *slog.Logger
	clock                     backend.Clock
	memberProvider            memberstore.MemberProvider
	qualificationProvider     qualificationstore.QualificationProvider
	memberRequirementProvider MemberRequirementProvider
}

func New(
	memberProvider memberstore.MemberProvider,
	qualificationProvider qualificationstore.QualificationProvider,
	memberRequirementProvider MemberRequirementProvider,
	logger *slog.Logger,
	clock backend.Clock,
) MemberQualificationStore {
	if clock == nil {
		clock = backend.RealClock{}
	}
	return MemberQualificationStore{
		logger:                    logger,
		memberProvider:            memberProvider,
		clock:                     clock,
		qualificationProvider:     qualificationProvider,
		memberRequirementProvider: memberRequirementProvider,
	}
}

type MemberRequirementProvider interface {
	AssignInitialMemberRequirement(memberID, requirementID, assignedBy string) error
	GetInitialMemberRequirement(memberID, requirementID string) (types.InitialMemberRequirement, error)
	GetInitialMemberRequirementsForQualification(memberID, qualificationID string) ([]types.InitialMemberRequirement, error)
	CompleteInitialMemberRequirement(completedDate time.Time, completedBy, memberID, requirementID string) error

	AssignRecurringMemberRequirement(id, memberID, requirementID, assignedBy string) error
	GetRecurringMemberRequirement(memberID, requirementID string) (types.RecurringMemberRequirement, error)
	GetRecurringMemberRequirementsForQualification(memberID, qualificationID string) ([]types.RecurringMemberRequirement, error)
	CompleteRecurringMemberRequirement(id, memberRequirementID, completedBy string, completedDate time.Time) error
}
