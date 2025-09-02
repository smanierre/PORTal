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
All: Reference, and Type. Notes are optional
Qualification type: Must have QualificationID, can't have Grade or DaysValidFor
WBT Type: Can't have QualificationID, Grade, or DaysValidFor, must have Name
Grade: Must have Grade, can't have QualificationID or DaysValidFor

As long as all the required information is present, any irrelevant information will be zeroed out.
Otherwise, an error will be returned with the missing fields
*/
func normalizeInitialRequirement(r types.Requirement) (types.Requirement, error) {
	var errors []string
	switch r.Type {
	case types.QualificationType:
		if r.QualificationID == "" {
			errors = append(errors, "Missing QualificationID")
		}
		r.Grade = ""
		r.DaysValidFor = 0
		break
	case types.WbtType:
		if r.Name == "" {
			errors = append(errors, "Missing Name")
		}
		r.QualificationID = ""
		r.Grade = ""
		r.DaysValidFor = 0
		break
	case types.GradeType:
		if !isValidGrade(r.Grade) {
			errors = append(errors, "Invalid Grade")
		}
		r.QualificationID = ""
		r.DaysValidFor = 0
		break
	default:
		return types.Requirement{}, fmt.Errorf("%w: Invalid type for initial requirement: %s", backend.ErrValidation, r.Type)
	}
	if r.Reference == "" {
		errors = append(errors, "Missing Reference")
	}
	if len(errors) > 0 {
		return types.Requirement{}, fmt.Errorf("%w: %s", backend.ErrValidation, errors)
	}
	return r, nil
}

func isValidGrade(grade types.Grade) bool {
	return grade == types.E1 ||
		grade == types.E2 ||
		grade == types.E3 ||
		grade == types.E4 ||
		grade == types.E5 ||
		grade == types.E6 ||
		grade == types.E7 ||
		grade == types.E8 ||
		grade == types.E9
}

/*
Recurring requirements have the following rules:
Can only be WbtType or ProficiencyType
All: Name, Reference, Type, and DaysValidFor.
WbtType: Can't have QualificationID or Grade. Notes optional
ProficiencyType: Notes required. Can't have QualificationID or Grade.
*/
func normalizeRecurringRequirement(r types.Requirement) (types.Requirement, error) {
	var errors []string
	switch r.Type {
	case types.WbtType:
		r.QualificationID = ""
		r.Grade = ""
		break
	case types.ProficiencyType:
		if r.Notes == "" {
			errors = append(errors, "Missing Notes")
		}
		r.QualificationID = ""
		r.Grade = ""
		break
	default:
		return types.Requirement{}, fmt.Errorf("%w: Invalid type for recurring requirement: %s", backend.ErrValidation, r.Type)
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
		return types.Requirement{}, fmt.Errorf("%w: %s", backend.ErrValidation, errors)
	}
	return r, nil
}
