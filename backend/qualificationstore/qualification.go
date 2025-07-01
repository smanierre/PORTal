package qualificationstore

import (
	"PORTal/types"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

func (q QualificationStore) AddQualification(qual types.Qualification) (types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID for new qualification")
	qual.ID = uuid.NewString()
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args")
	if err := checkQualificationForMissingArgs(qual); err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Missing required arguments", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	return qual, q.provider.AddQualification(qual)
}

func (q QualificationStore) AssignRequirementToQualification(qualificationID, requirementID string, initial bool) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Assigning requirement to qualification",
		slog.String("qualification_id", qualificationID), slog.String("requirement_id", requirementID))
	return q.provider.AssignRequirementToQualification(qualificationID, requirementID, initial)
}

func (q QualificationStore) GetQualification(id string) (types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualification", slog.String("id", id))
	return q.provider.GetQualification(id)
}

func (q QualificationStore) GetAllQualifications() ([]types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all qualifications")
	quals, err := q.provider.GetAllQualifications()
	if err != nil {
		return []types.Qualification{}, err
	}
	return quals, nil
}

func (q QualificationStore) UpdateQualification(qual types.Qualification) (types.Qualification, error) {
	if err := checkQualificationForMissingArgs(qual); err != nil {
		return types.Qualification{}, err
	}
	err := q.provider.UpdateQualification(qual)
	if err != nil {
		return types.Qualification{}, err
	}
	newQual, err := q.provider.GetQualification(qual.ID)
	if err != nil {
		return types.Qualification{}, err
	}
	return newQual, nil
}

func (q QualificationStore) DeleteQualification(id string) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification")
	return q.provider.DeleteQualification(id)
}
