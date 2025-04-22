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

/*
Initial requirements have the following rules:
Can't be Proficiency type.
All: Name, reference, and Type. Notes are optional
Qualification type: Must have QualificationID, can't have Grade or DaysValidFor
WBT Type: Can't have QualificationID, Grade, or DaysValidFor
Grade: Must have Grade, can't have QualificationID or DaysValidFor
*/
func validateInitialRequirement(r types.Requirement) error {
	var errors []string
	switch r.Type {
	case types.QualificationType:
		if r.QualificationID == "" {
			errors = append(errors, "Missing QualificationID")
		}
		if r.Grade != "" {
			errors = append(errors, "Grade should be blank")
		}
		if r.DaysValidFor != 0 {
			errors = append(errors, "DaysValidFor should be 0")
		}
		break
	case types.WbtType:
		if r.QualificationID != "" {
			errors = append(errors, "QualificationID should be blank")
		}
		if r.Grade != "" {
			errors = append(errors, "Grade should be blank")
		}
		if r.DaysValidFor != 0 {
			errors = append(errors, "DaysValidFor should be 0")
		}
		break
	case types.GradeType:
		if r.Grade == "" {
			errors = append(errors, "Missing Grade")
		}
		if r.QualificationID != "" {
			errors = append(errors, "QualificationID should be blank")
		}
		if r.DaysValidFor != 0 {
			errors = append(errors, "DaysValidFor should be 0")
		}
		break
	default:
		return fmt.Errorf("%w: Invalid type for initial requirement: %s", backend.ErrValidation, r.Type)
	}
	if r.Name == "" {
		errors = append(errors, "Missing Name")
	}
	if r.Reference == "" {
		errors = append(errors, "Missing Reference")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", backend.ErrValidation, errors)
	}
	return nil
}

/*
Recurring requirements have the following rules:
Can only be WbtType or ProficiencyType
All: Name, Reference, Type, and DaysValidFor.
WbtType: Can't have QualificationID or Grade. Notes optional
ProficiencyType: Notes required. Can't have QualificationID or Grade.
*/
func validateRecurringRequirement(r types.Requirement) error {
	var errors []string
	switch r.Type {
	case types.WbtType:
		if r.QualificationID != "" {
			errors = append(errors, "QualificationID should be blank")
		}
		if r.Grade != "" {
			errors = append(errors, "Grade should be blank")
		}
		break
	case types.ProficiencyType:
		if r.Notes == "" {
			errors = append(errors, "Missing Notes")
		}
		if r.QualificationID != "" {
			errors = append(errors, "QualificationID should be blank")
		}
		if r.Grade != "" {
			errors = append(errors, "Grade should be blank")
		}
		break
	default:
		return fmt.Errorf("%w: Invalid type for recurring requirement: %s", backend.ErrValidation, r.Type)
	}
	if r.Name == "" {
		errors = append(errors, "Missing Name")
	}
	if r.Reference == "" {
		errors = append(errors, "Missing Reference")
	}
	if r.DaysValidFor == 0 {
		errors = append(errors, "Missing DaysValidFor")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", backend.ErrValidation, errors)
	}
	return nil
}
