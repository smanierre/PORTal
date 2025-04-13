package qualificationstore

import (
	"PORTal/backend"
	"PORTal/types"
	"fmt"
)

func checkQualificationForMissingArgs(q types.Qualification) error {
	var missing []string
	if q.ID == "" {
		missing = append(missing, "ID")
	}
	if q.Name == "" {
		missing = append(missing, "Name")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", backend.ErrMissingArgs, missing)
	}
	return nil
}

func checkRequirementForMissingArgs(r types.Requirement) error {
	var errors []string
	if (r.Type != types.QualificationType && r.Type != types.GradeType) && r.Name == "" {
		errors = append(errors, "Name")
	}
	if (r.Type != types.QualificationType && r.Type != types.GradeType) && r.DaysValidFor == 0 {
		errors = append(errors, "DaysValidFor")
	}
	if r.Reference == "" {
		errors = append(errors, "Reference")
	}
	if r.Type == "" {
		errors = append(errors, "Type")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", backend.ErrMissingArgs, errors)
	}
	return nil
}
