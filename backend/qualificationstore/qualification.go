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
	err := q.provider.AddQualification(qual)
	if err != nil {
		return types.Qualification{}, err
	}
	for i, ir := range qual.InitialRequirements {
		ir, err := q.AddRequirement(ir)
		if err != nil {
			continue
		}
		err = q.AssignRequirementToQualification(qual.ID, ir.ID, true)
		qual.InitialRequirements[i] = ir
	}
	for i, rr := range qual.RecurringRequirements {
		rr, err := q.AddRequirement(rr)
		if err != nil {
			continue
		}
		err = q.AssignRequirementToQualification(qual.ID, rr.ID, false)
		if err != nil {
			continue
		}
		qual.RecurringRequirements[i] = rr
	}
	return qual, nil
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
	for _, v := range qual.InitialRequirements {
		found := false
		for i, v2 := range newQual.InitialRequirements {
			if v.ID == v2.ID {
				found = true
				if v != v2 {
					newReq, err := q.UpdateRequirement(v)
					if err != nil {
						break
					}
					newQual.InitialRequirements[i] = newReq
				}
			}
		}
		if !found {
			r, err := q.AddRequirement(v)
			if err != nil {
				return types.Qualification{}, err
			}
			err = q.AssignRequirementToQualification(qual.ID, r.ID, true)
			if err != nil {
				continue
			}
			newQual.InitialRequirements = append(newQual.InitialRequirements, r)
		}
	}
	for _, v := range qual.RecurringRequirements {
		found := false
		for i, v2 := range newQual.RecurringRequirements {
			if v.ID == v2.ID {
				found = true
				if v != v2 {
					newReq, err := q.UpdateRequirement(v)
					if err != nil {
						break
					}
					newQual.RecurringRequirements[i] = newReq
				}
			}
		}
		if !found {
			r, err := q.AddRequirement(v)
			if err != nil {
				return types.Qualification{}, err
			}
			err = q.AssignRequirementToQualification(qual.ID, r.ID, false)
			if err != nil {
				continue
			}
			newQual.RecurringRequirements = append(newQual.InitialRequirements, r)
		}
	}
	// Need to refetch qualification to get the new requirements that were assigned to it
	newQual, err = q.provider.GetQualification(qual.ID)
	if err != nil {
		return types.Qualification{}, err
	}
	for i, v := range newQual.InitialRequirements {
		found := false
		for _, v2 := range qual.InitialRequirements {
			// If it's a new requirement, no need to worry about deleting it
			if v2.ID == "" || v.ID == v2.ID {
				found = true
			}
		}
		if !found {
			err = q.DeleteRequirement(v.ID)
			if err != nil {
				continue
			}
			newQual.InitialRequirements = append(newQual.InitialRequirements[:i], newQual.InitialRequirements[i+1:]...)
		}
	}
	for i, v := range newQual.RecurringRequirements {
		found := false
		for _, v2 := range qual.RecurringRequirements {
			// If it's a new requirement, no need to worry about deleting it
			if v2.ID == "" || v.ID == v2.ID {
				found = true
			}
		}
		if !found {
			err = q.DeleteRequirement(v.ID)
			if err != nil {
				continue
			}
			newQual.RecurringRequirements = append(newQual.RecurringRequirements[:i], newQual.RecurringRequirements[i+1:]...)
		}
	}
	return newQual, nil
}

func (q QualificationStore) DeleteQualification(id string) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification")
	return q.provider.DeleteQualification(id)
}
