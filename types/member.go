package types

import (
	"fmt"
	"log/slog"
	"time"
)

type Grade string

func (g Grade) CanSupervise(grade Grade) bool {
	// There's probably a better way, but this won't really need to be changed
	switch g {
	case E9:
		return true
	case E8:
		return grade != E9
	case E7:
		return grade != E8 && grade != E9
	case E6:
		return grade != E7 && grade != E8 && grade != E9
	case E5:
		return grade != E6 && grade != E7 && grade != E8 && grade != E9
	case E4:
		return grade != E5 && grade != E6 && grade != E7 && grade != E8 && grade != E9
	case E3:
		return grade != E4 && grade != E5 && grade != E6 && grade != E7 && grade != E8 && grade != E9
	case E2:
		return grade != E3 && grade != E4 && grade != E5 && grade != E6 && grade != E7 && grade != E8 && grade != E9
	case E1:
		return grade == E1
	default:
		return false
	}
}

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

func FilterMembers(members []Member, test func(m Member) bool) []Member {
	var result []Member
	for _, member := range members {
		if test(member) {
			result = append(result, member)
		}
	}
	return result
}

type Session struct {
	SessionID string
	UserAgent string
	IpAddress string
	Expires   time.Time
}
