package backend

import (
	"PORTal/types"
	"fmt"
)

func CheckMemberForMissingArgs(m types.Member) error {
	var errors []string
	if m.ID == "" {
		errors = append(errors, "ID")
	}
	if m.FirstName == "" {
		errors = append(errors, "FirstName")
	}
	if m.LastName == "" {
		errors = append(errors, "LastName")
	}
	if m.Grade == "" {
		errors = append(errors, "Grade")
	}
	if m.Username == "" {
		errors = append(errors, "Username")
	}
	if m.Password == "" {
		errors = append(errors, "Password")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, errors)
	}
	return nil
}

func CheckQualificationForMissingArgs(q types.Qualification) error {
	var missing []string
	if q.ID == "" {
		missing = append(missing, "ID")
	}
	if q.Name == "" {
		missing = append(missing, "Name")
	}
	if q.Expires && q.ExpirationInterval == 0 {
		missing = append(missing, "ExpirationInterval")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, missing)
	}
	return nil
}

func CheckRequirementForMissingArgs(r types.Requirement) error {
	var errors []string
	if r.Name == "" {
		errors = append(errors, "Name")
	}
	if r.DaysValidFor == 0 {
		errors = append(errors, "DaysValidFor")
	}
	if r.Reference == "" {
		errors = append(errors, "Reference")
	}
	if r.Type == "" {
		errors = append(errors, "Type")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, errors)
	}
	return nil
}
