package qualificationstore

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

func (q QualificationStore) AddRequirement(r types.Requirement) (types.Requirement, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID...")
	r.ID = uuid.NewString()

	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args...")
	if err := checkRequirementForMissingArgs(r); err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "Required arguments missing", slog.String("error", err.Error()))
		return types.Requirement{}, err
	}
	return r, q.provider.AddRequirement(r)
}

func (q QualificationStore) GetRequirement(id string) (types.Requirement, error) {
	return q.provider.GetRequirement(id)
}

func (q QualificationStore) GetAllRequirements() ([]types.Requirement, error) {
	reqs, err := q.provider.GetAllRequirements()
	if err != nil {
		return nil, err
	}
	return reqs, nil
}

func (q QualificationStore) UpdateRequirement(r types.Requirement) (types.Requirement, error) {
	err := q.provider.UpdateRequirement(r)
	if err != nil {
		return types.Requirement{}, err
	}
	return r, nil
}

func (q QualificationStore) DeleteRequirement(id string) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking if requirement is assigned to any qualifications")
	quals, err := q.provider.GetQualificationIDsForRequirement(id)
	if err != nil {
		return err
	}
	if len(quals) > 0 {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "Requirement is still assigned to qualifications", slog.Any("qualification_ids", quals))
		return fmt.Errorf("%w: %v", backend.ErrRequirementInUse, quals)
	}
	return q.provider.DeleteRequirement(id)
}
