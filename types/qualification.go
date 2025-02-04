package types

import (
	"fmt"
	"log/slog"
	"time"
)

type RequirementType string

const (
	QualificationType RequirementType = "Qualification"
	WbtType           RequirementType = "WBT"
	GradeType         RequirementType = "Grade"
)

var Never time.Time

func init() {
	var err error
	Never, err = time.Parse(time.DateOnly, "9999-12-31")
	if err != nil {
		panic(fmt.Sprintf("Error initializing never value: %s", err.Error()))
	}
}

type Qualification struct {
	ID                    string
	Name                  string
	InitialRequirements   []Requirement
	RecurringRequirements []Requirement
	Notes                 string
	Expires               bool
	ExpirationInterval    time.Duration
}

func (q Qualification) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s, Name: %s, Notes: %s, Expires: %t, Expiration Days: %d, Initial Requirements: %v, Recurring Requirements: %v",
		q.ID, q.Name, q.Notes, q.Expires, q.ExpirationInterval, q.InitialRequirements, q.RecurringRequirements))
}

type Requirement struct {
	ID           string
	Name         string
	Reference    string
	Notes        string
	DaysValidFor int
	Type         RequirementType
}

func (r Requirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Name: %s Notes: %s DaysValidFor: %d", r.ID, r.Name, r.Notes, r.DaysValidFor))
}

type MemberRequirement struct {
	MemberID string
	Requirement
	Completed     bool
	CompletedDate time.Time
}
