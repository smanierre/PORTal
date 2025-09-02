package memberstore

import (
	"PORTal/backend"
	"PORTal/types"
	"fmt"
)

func checkMemberForInvalidArgs(m types.Member, ranks types.RankMap) error {
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
		var found bool
		for rank, _ := range ranks {
			if rank == m.Grade {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, "Grade")
		}
	}
	if m.Username == "" {
		errors = append(errors, "Username")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", backend.ErrMissingArgs, errors)
	}
	return nil
}
