package types

import (
	"fmt"
	"log/slog"
	"time"
)

var Never time.Time
var Day = time.Hour * 24

func init() {
	var err error
	Never, err = time.Parse(time.DateOnly, "9999-12-31")
	if err != nil {
		panic(fmt.Sprintf("Error initializing never value: %s", err.Error()))
	}
}

type Qualification struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	InitialRequirements   []Requirement `json:"initial_requirements"`
	RecurringRequirements []Requirement `json:"recurring_requirements,omitempty"`
	Notes                 string        `json:"notes,omitempty"`
	References            []Reference   `json:"references,omitempty"`
	Expires               bool          `json:"expires"`
	ExpirationDays        int           `json:"expiration_days,omitempty"`
}

func (q Qualification) CheckForMissingArgs() error {
	missing := []string{}
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

type MemberQualification struct {
	MemberID      string `json:"member_id"`
	Qualification `json:"qualification"`
	Active        bool      `json:"active"`
	ActiveDate    time.Time `json:"active_date"`
}

func (m MemberQualification) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("MemberID: %s QualificationID: %s Active: %t ActiveDate: %s", m.MemberID, m.Qualification.ID, m.Active, m.ActiveDate.String()))
}

type Requirement struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Notes        string `json:"notes,omitempty"`
	DaysValidFor int    `json:"days_valid_for,omitempty"`
}

func (r Requirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Name: %s Description: %s Notes: %s DaysValidFor: %d", r.ID, r.Name, r.Description, r.Notes, r.DaysValidFor))
}

func (r Requirement) CheckForMissingArgs() error {
	var errors []string
	if r.ID == "" {
		errors = append(errors, "ID")
	}
	if r.Name == "" {
		errors = append(errors, "Name")
	}
	if r.DaysValidFor == 0 {
		errors = append(errors, "DaysValidFor")
	}
	if r.Description == "" {
		errors = append(errors, "Description")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, errors)
	}
	return nil
}

type MemberRequirement struct {
	MemberID      string `json:"member_id"`
	Requirement   `json:"requirement"`
	Completed     bool      `json:"completed"`
	CompletedDate time.Time `json:"completed_date,omitempty"`
}

type Reference struct {
	Name      string `json:"name"`
	Volume    int    `json:"volume"`
	Paragraph string `json:"paragraph"`
}
