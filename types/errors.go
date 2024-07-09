package types

import "errors"

var (
	ErrMemberNotFound               = errors.New("member with that id not found")
	ErrSupervisorNotFound           = errors.New("supervisor with that id not found")
	ErrQualificationNotFound        = errors.New("qualification with that id not found")
	ErrRequirementNotFound          = errors.New("requirement with that id not found")
	ErrMissingArgs                  = errors.New("missing required arguments")
	ErrBadUpdate                    = errors.New("invalid update")
	ErrQualificationAlreadyAssigned = errors.New("member is already assigned qualification")
	ErrMemberQualificationNotFound  = errors.New("could not find given qualification for member")
	ErrInvalidIP                    = errors.New("invalid IP address provided")
	ErrSessionValidationFailed      = errors.New("failed to validate session")
	ErrUsernameAlreadyExists        = errors.New("that username already exists in the database")
	ErrPasswordAuthenticationFailed = errors.New("unable to validate password")
)
