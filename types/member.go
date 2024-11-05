package types

import (
	"fmt"
	"log/slog"
	"time"
)

type Grade string

const (
	E1 Grade = "E1"
	E2 Grade = "E2"
	E3 Grade = "E3"
	E4 Grade = "E4"
	E5 Grade = "E5"
	E6 Grade = "E6"
	E7 Grade = "E7"
	E8 Grade = "E8"
	E9 Grade = "E9"
)

type Member struct {
	ApiMember
	Password string `json:"password,omitempty"`
	Hash     string
}

func (m Member) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Member: %s %s %s Username: %s Supervisor ID: %s Admin: %t", m.ID, m.Rank, m.FirstName, m.LastName, m.Username, m.SupervisorID, m.Admin))
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
	//TODO: Refactor this to check for explicitly blank supervisor ID
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
	Rank         Grade  `json:"rank"`
	SupervisorID string `json:"supervisor_id"`
	Admin        bool   `json:"admin"`
}

type Session struct {
	SessionID string
	UserAgent string
	Expires   time.Time
}
