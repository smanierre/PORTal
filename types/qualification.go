package types

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

var Never time.Time
var Always time.Time
var Day = time.Hour * 24

func init() {
	var err error
	Never, err = time.Parse(time.DateOnly, "9999-12-31")
	if err != nil {
		panic(fmt.Sprintf("Error initializing never value: %s", err.Error()))
	}
	Always, err = time.Parse(time.DateOnly, "9999-12-31")
	if err != nil {
		panic(fmt.Sprintf("Error initializing always value: %s", err.Error()))
	}
}

type JSONSafeSlice[T any] []T

func (s JSONSafeSlice[T]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]T(s))
}

type Qualification struct {
	ID                    string                     `json:"id"`
	Name                  string                     `json:"name"`
	InitialRequirements   JSONSafeSlice[Requirement] `json:"initial_requirements"`
	RecurringRequirements JSONSafeSlice[Requirement] `json:"recurring_requirements"`
	Notes                 string                     `json:"notes"`
	Expires               bool                       `json:"expires"`
	ExpirationDays        int                        `json:"expiration_days"`
}

func (q Qualification) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s, Name: %s, Notes: %s, Expires: %t, Expiration Days: %d, Initial Requirements: %v, Recurring Requirements: %v",
		q.ID, q.Name, q.Notes, q.Expires, q.ExpirationDays, q.InitialRequirements, q.RecurringRequirements))
}

type Requirement struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Reference    Reference `json:"reference"`
	Notes        string    `json:"notes"`
	DaysValidFor int       `json:"days_valid_for"`
}

func (r Requirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Name: %s Notes: %s DaysValidFor: %d", r.ID, r.Name, r.Notes, r.DaysValidFor))
}

type MemberRequirement struct {
	MemberID      string `json:"member_id"`
	Requirement   `json:"requirement"`
	Completed     bool      `json:"completed"`
	CompletedDate time.Time `json:"completed_date"`
}

type Reference struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Volume    int    `json:"volume"`
	Paragraph string `json:"paragraph"`
}

func (r Reference) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s, Name: %s, Volume: %d, Paragraph: %s", r.ID, r.Name, r.Volume, r.Paragraph))
}
