package qualificationstore

import (
	"PORTal/types"
	"log/slog"
)

type QualificationStore struct {
	provider QualificationProvider
	logger   *slog.Logger
}

func New(provider QualificationProvider, logger *slog.Logger) QualificationStore {
	logger = logger.With(slog.String("source", "qualificationstore"))
	return QualificationStore{
		provider: provider,
		logger:   logger,
	}
}

type QualificationProvider interface {
	AddQualification(q types.Qualification) error
	AssignRequirementToQualification(qualificationID, requirementID string, initial bool) error
	GetQualification(id string) (types.Qualification, error)
	GetAllQualifications() ([]types.Qualification, error)
	UpdateQualification(q types.Qualification) error
	DeleteQualification(id string) error
	GetMemberQualification(memberID, qualificationID string) (types.Qualification, error)
	GetMemberQualifications(memberID string) ([]types.Qualification, error)
	AddRequirement(r types.Requirement) error
	GetRequirement(id string) (types.Requirement, error)
	GetAllRequirements() ([]types.Requirement, error)
	UpdateRequirement(r types.Requirement) error
	DeleteRequirement(id string) error
}
