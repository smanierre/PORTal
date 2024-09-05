package types

import (
	"fmt"
	"log/slog"
	"time"
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
	return slog.StringValue(fmt.Sprintf("ID: %s Member: %s %s %s Username: %s Supervisor ID: %s", m.ID, m.Rank, m.FirstName, m.LastName, m.Username, m.SupervisorID))
}

func (m Member) ToApiMember() ApiMember {
	return m.ApiMember
}

func (m Member) MergeIn(new Member) Member {
	if new.FirstName != "" {
		m.FirstName = new.FirstName
	}
	if new.LastName != "" {
		m.LastName = new.LastName
	}
	if new.Rank != "" {
		m.Rank = new.Rank
	}
	if new.SupervisorID != "" {
		m.SupervisorID = new.SupervisorID
	}
	if new.Username != "" {
		m.Username = new.Username
	}
	if new.Password != "" {
		m.Password = new.Password
	}
	return m
}

type ApiMember struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Rank         Rank   `json:"rank"`
	SupervisorID string `json:"supervisor_id,omitempty"`
}

type Session struct {
	SessionID string
	UserAgent string
	Expires   time.Time
}
