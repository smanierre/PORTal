package memberstore

import (
	"PORTal/backend"
	"PORTal/types"
	"fmt"
)

func checkMemberForMissingArgs(m types.Member) error {
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
		return fmt.Errorf("%w: %s", backend.ErrMissingArgs, errors)
	}
	return nil
}
