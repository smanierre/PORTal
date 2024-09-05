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
	if m.Rank == "" {
		errors = append(errors, "Rank")
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
	if q.Expires && q.ExpirationDays == 0 {
		missing = append(missing, "ExpirationDays")
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
	if r.Description == "" {
		errors = append(errors, "Description")
	}
	if err := CheckReferenceForMissingArgs(r.Reference); err != nil {
		errors = append(errors, fmt.Sprintf("Reference: %s", err.Error()))
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, errors)
	}
	return nil
}

func CheckReferenceForMissingArgs(r types.Reference) error {
	var errors []string
	if r.ID == "" {
		errors = append(errors, "ID")
	}
	if r.Name == "" {
		errors = append(errors, "Name")
	}
	if r.Paragraph == "" {
		errors = append(errors, "Paragraph")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, errors)
	}
	return nil
}
