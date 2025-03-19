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

func GetRequirementTypes() []RequirementType {
	return []RequirementType{
		QualificationType, WbtType, GradeType,
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

func (q Qualification) GetID() string {
	return q.ID
}

func (q Qualification) Display() string {
	return q.Name
}

type Requirement struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Reference       string          `json:"reference"`
	QualificationID string          `json:"qualification_id"`
	Grade           Grade           `json:"grade"`
	Notes           string          `json:"notes"`
	DaysValidFor    int             `json:"days_valid_for"`
	Type            RequirementType `json:"type"`
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
