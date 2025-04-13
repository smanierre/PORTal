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
	ProficiencyType   RequirementType = "Proficiency"
)

func GetInitialRequirementTypes() []RequirementType {
	return []RequirementType{
		QualificationType, WbtType, GradeType,
	}
}

func GetRecurringRequirementTypes() []RequirementType {
	return []RequirementType{
		WbtType, ProficiencyType,
	}
}

type Qualification struct {
	ID                    string
	Name                  string
	InitialRequirements   []Requirement
	RecurringRequirements []Requirement
	Notes                 string
}

func (q Qualification) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s, Name: %s, Notes: %s, Initial Requirements: %v, Recurring Requirements: %v",
		q.ID, q.Name, q.Notes, q.InitialRequirements, q.RecurringRequirements))
}

func (q Qualification) GetID() string {
	return q.ID
}

func (q Qualification) Display() string {
	return q.Name
}

func FilterQualifications(qualifications []Qualification, test func(q Qualification) bool) []Qualification {
	var result []Qualification
	for _, qualification := range qualifications {
		if test(qualification) {
			result = append(result, qualification)
		}
	}
	return result
}

type Requirement struct {
	ID              string
	Name            string
	Reference       string
	QualificationID string
	Grade           Grade
	Notes           string
	DaysValidFor    int
	Type            RequirementType
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
