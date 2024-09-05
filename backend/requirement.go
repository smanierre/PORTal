package backend

import (
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
)

func (b Backend) AddRequirement(r types.Requirement) (types.Requirement, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args...")
	if err := CheckRequirementForMissingArgs(r); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Required arguments missing", slog.String("error", err.Error()))
		return types.Requirement{}, err
	}
	return r, b.requirementProvider.AddRequirement(r)
}

func (b Backend) GetRequirement(id string) (types.Requirement, error) {
	return b.requirementProvider.GetRequirement(id)
}

func (b Backend) GetAllRequirements() ([]types.Requirement, error) {
	return b.requirementProvider.GetAllRequirements()
}

func (b Backend) UpdateRequirement(r types.Requirement) (types.Requirement, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting existing requirement to determine updates")
	existingReq, err := b.GetRequirement(r.ID)
	if err != nil {
		return types.Requirement{}, err
	}
	existingReq = existingReq.MergeIn(r)
	err = b.requirementProvider.UpdateRequirement(existingReq)
	if err != nil {
		return types.Requirement{}, err
	}
	return r, nil
}

func (b Backend) DeleteRequirement(id string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking if requirement is assigned to any qualifications")
	quals, err := b.requirementProvider.GetQualificationIDsForRequirement(id)
	if err != nil {
		return err
	}
	if len(quals) > 0 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Requirement is still assigned to qualifications", slog.Any("qualification_ids", quals))
		return fmt.Errorf("%w: %v", ErrRequirementInUse, quals)
	}
	return b.requirementProvider.DeleteRequirement(id)
}
