package types

import (
	"fmt"
	"log/slog"
)

type Rank string

const (
	E1 Rank = "AB"
	E2 Rank = "Amn"
	E3 Rank = "A1C"
	E4 Rank = "SrA"
	E5 Rank = "SSgt"
	E6 Rank = "TSgt"
	E7 Rank = "MSgt"
	E8 Rank = "SMSgt"
	E9 Rank = "CMSgt"
)

type Member struct {
	ApiMember
	Password string `json:"password,omitempty"`
	Hash     string
}

func (m Member) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Member: %s %s %s Supervisor ID: %s Qualifications: %v", m.ID, m.Rank, m.FirstName, m.LastName, m.SupervisorID, m.Qualifications))
}

func (m Member) ToApiMember() ApiMember {
	return m.ApiMember
}

func (m Member) CheckForMissingArgs() error {
	errors := []string{}
	if m.ID == "" {
		errors = append(errors, "ID")
	}
	if m.FirstName == "" {
		errors = append(errors, "FirstName")
	}
	if m.LastName == "" {
		errors = append(errors, "LastName")
	}
	if m.Rank == "" {
		errors = append(errors, "Rank")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingArgs, errors)
	}
	return nil
}

type ApiMember struct {
	ID             string                `json:"id"`
	FirstName      string                `json:"first_name"`
	LastName       string                `json:"last_name"`
	Rank           Rank                  `json:"rank"`
	Qualifications []MemberQualification `json:"qualifications,omitempty"`
	SupervisorID   string                `json:"supervisor_id,omitempty"`
}
