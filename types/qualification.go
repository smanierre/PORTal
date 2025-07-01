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

var Never = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)

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
	Initial         bool
	Reference       string
	Notes           string
	Type            RequirementType
	QualificationID string
	Grade           Grade
	DaysValidFor    int
}

func (r Requirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Name: %s Initial: %t Reference: %s Notes: %s Type: %s QualificationID: %s Grade %s DaysValidFor: %d",
		r.ID, r.Name, r.Initial, r.Reference, r.Notes, r.Type, r.QualificationID, r.Grade, r.DaysValidFor))
}

type MemberQualification struct {
	MemberID        string
	QualificationID string
	DateAssigned    time.Time
	AssignedByID    string
}

func (m MemberQualification) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("MemberID: %s QualificationID: %s, DateAssigned: %s AssignedBy: %s",
		m.MemberID, m.QualificationID, m.DateAssigned.String(), m.AssignedByID))
}

type InitialMemberRequirement struct {
	MemberID      string
	RequirementID string
	CompletedDate time.Time
	AssignedBy    string
	CompletedBy   string
}

func (i InitialMemberRequirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("MemberID: %s RequirementID: %s CompletedDate: %s CompletedBy: %s",
		i.MemberID, i.RequirementID, i.CompletedDate.String(), i.CompletedBy))
}

type RecurringMemberRequirement struct {
	ID                string
	MemberID          string
	RequirementID     string
	AssignedBy        string
	CompletionHistory []RecurringMemberRequirementCompletion
}

func (r RecurringMemberRequirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s MemberID: %s RequirementID: %s CompletionHistory: %v", r.ID, r.MemberID, r.RequirementID, r.CompletionHistory))
}

type RecurringMemberRequirementCompletion struct {
	ID             string
	CompletionDate time.Time
	CompletedBy    string
}
