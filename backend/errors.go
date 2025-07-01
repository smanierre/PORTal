package backend

import "errors"

var (
	ErrAuthenticationFailed               = errors.New("unable to authenticate user")
	ErrDuplicateRequirement               = errors.New("requirement with that name already exists")
	ErrDuplicateUsername                  = errors.New("member with that username already exists")
	ErrInvalidQualExpiration              = errors.New("invalid expiration length for qualification")
	ErrInitialMemberRequirementNotFound   = errors.New("initial member requirement not found")
	ErrInvalidRequirementCompletionDate   = errors.New("invalid requirement completion date")
	ErrRecurringMemberRequirementNotFound = errors.New("recurring member requirement not found")
	ErrMemberDisabled                     = errors.New("member is disabled")
	ErrMemberNotFound                     = errors.New("member with that id not found")
	ErrMemberQualificationNotFound        = errors.New("member with given qualification not found")
	ErrMissingArgs                        = errors.New("missing required arguments")
	ErrPasswordTooLong                    = errors.New("password exceeds maximum length of 72 characters")
	ErrQualificationAlreadyAssigned       = errors.New("qualification already assigned to member")
	ErrQualificationNotFound              = errors.New("qualification with that id not found")
	ErrReferenceNotFound                  = errors.New("unable to find reference with given id")
	ErrRequirementInUse                   = errors.New("requirement is assigned to qualification")
	ErrRequirementNotFound                = errors.New("requirement with that identifier not found")
	ErrSessionNotFound                    = errors.New("no session with that ID found")
	ErrSessionValidationFailed            = errors.New("session validation failed")
	ErrSupervisorNotFound                 = errors.New("supervisor with that ID not found")
	ErrValidation                         = errors.New("validation failed")
	ErrWeakPassword                       = errors.New("supplied password doesn't meet requirements")
)
