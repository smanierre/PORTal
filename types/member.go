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
	ID           string
	FirstName    string
	LastName     string
	Username     string
	Grade        Grade
	SupervisorID string
	Admin        bool
	Password     string
	Hash         string
	Disabled     bool
}

func (m Member) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("ID: %s Member: %s %s %s Username: %s Supervisor ID: %s Admin: %t", m.ID, m.Grade, m.FirstName, m.LastName, m.Username, m.SupervisorID, m.Admin))
}

func (m Member) GetPotentialSupervisors(members []Member) []Member {
	var ps []Member
	for _, member := range members {
		if m.Grade <= member.Grade && member.ID != m.ID {
			ps = append(ps, member)
		}
	}
	return ps
}

func (m Member) GetID() string {
	return m.ID
}

func (m Member) Display() string {
	return fmt.Sprintf("%s %s", m.FirstName, m.LastName)
}

func GetRank(member Member, service string) string {
	switch service {
	case "f":
		return AfRankMap[member.Grade]
	case "a":
		return ArmyRankMap[member.Grade]
	case "m":
		return MarineRankMap[member.Grade]
	case "n":
		return NavyRankMap[member.Grade]
	default:
		return ""
	}
}

type Session struct {
	SessionID string
	UserAgent string
	IpAddress string
	Expires   time.Time
}
