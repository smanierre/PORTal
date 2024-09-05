package backend

import "errors"

var (
	ErrAuthenticationFailed         = errors.New("unable to authenticate user")
	ErrBadUpdate                    = errors.New("supplied update values are invalid")
	ErrDuplicateReference           = errors.New("reference with that name already exists")
	ErrDuplicateRequirement         = errors.New("requirement with that name already exists")
	ErrDuplicateUsername            = errors.New("member with that username already exists")
	ErrInvalidQualExpiration        = errors.New("invalid expiration length for qualification")
	ErrMemberNotFound               = errors.New("member with that id not found")
	ErrMemberQualificationNotFound  = errors.New("member with given qualification not found")
	ErrMissingArgs                  = errors.New("missing required arguments")
	ErrPasswordTooLong              = errors.New("password exceeds maximum length of 72 characters")
	ErrQualificationAlreadyAssigned = errors.New("qualification already assigned to member")
	ErrQualificationNotFound        = errors.New("qualification with that id not found")
	ErrReferenceNotFound            = errors.New("unable to find reference with given id")
	ErrRequirementInUse             = errors.New("requirement is assigned to qualification")
	ErrRequirementNotFound          = errors.New("requirement with that identifier not found")
	ErrSessionValidationFailed      = errors.New("failed to validate session for member")
	ErrSupervisorNotFound           = errors.New("supervisor with that ID not found")
	ErrWeakPassword                 = errors.New("supplied password doesn't meet requirements")
)
