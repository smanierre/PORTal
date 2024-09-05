package types

import (
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

type Qualification struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	InitialRequirements   []Requirement `json:"initial_requirements"`
	RecurringRequirements []Requirement `json:"recurring_requirements,omitempty"`
	Notes                 string        `json:"notes,omitempty"`
	Expires               bool          `json:"expires"`
	ExpirationDays        int           `json:"expiration_days,omitempty"`
}

func (q Qualification) MergeIn(incoming Qualification, forceUpdateExpiration bool) Qualification {
	if incoming.Name != "" {
		q.Name = incoming.Name
	}
	if incoming.InitialRequirements != nil {
		q.InitialRequirements = incoming.InitialRequirements
	}
	if incoming.RecurringRequirements != nil {
		q.RecurringRequirements = incoming.RecurringRequirements
	}
	if incoming.Notes != "" {
		q.Notes = incoming.Notes
	}
	if incoming.Expires == false && incoming.ExpirationDays == 0 && forceUpdateExpiration {
		q.Expires = false
		q.ExpirationDays = 0
	} else if incoming.Expires == true {
		q.Expires = true
	}
	if incoming.ExpirationDays != 0 && (q.Expires == true || incoming.Expires == true) {
		q.ExpirationDays = incoming.ExpirationDays
	}
	return q
}

type Requirement struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Reference    Reference `json:"reference"`
	Description  string    `json:"description"`
	Notes        string    `json:"notes,omitempty"`
	DaysValidFor int       `json:"days_valid_for,omitempty"`
}

func (r Requirement) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Name: %s Description: %s Notes: %s DaysValidFor: %d", r.ID, r.Name, r.Description, r.Notes, r.DaysValidFor))
}

func (r Requirement) MergeIn(incoming Requirement) Requirement {
	if incoming.Name != "" {
		r.Name = incoming.Name
	}
	if incoming.Notes != "" {
		r.Notes = incoming.Notes
	}
	if incoming.Description != "" {
		r.Description = incoming.Description
	}
	if incoming.DaysValidFor != 0 {
		r.DaysValidFor = incoming.DaysValidFor
	}
	if incoming.Reference.ID != "" && incoming.Reference.ID != r.Reference.ID {
		r.Reference = incoming.Reference
	}
	return r
}

type MemberRequirement struct {
	MemberID      string `json:"member_id"`
	Requirement   `json:"requirement"`
	Completed     bool      `json:"completed"`
	CompletedDate time.Time `json:"completed_date,omitempty"`
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

func (r Reference) MergeIn(incoming Reference, overrideNoVolume bool) Reference {
	if incoming.Name != "" {
		r.Name = incoming.Name
	}
	if incoming.Volume == 0 && overrideNoVolume {
		r.Volume = 0
	} else if incoming.Volume > 0 {
		r.Volume = incoming.Volume
	}
	if incoming.Paragraph != "" {
		r.Paragraph = incoming.Paragraph
	}
	return r
}
