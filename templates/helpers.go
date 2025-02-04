package templates

import (
	"PORTal/types"
	"fmt"
	"strconv"
	"time"
)

type config struct {
	Organization string
	Service      string
}

var cfg config

func Initialize(organization string, service string) {
	cfg = config{organization, service}
}

func DisplayName(m types.Member) string {
	return fmt.Sprintf("%s %s %s", types.GetRank(m, cfg.Service), m.FirstName, m.LastName)
}

func GetOrg() string {
	return cfg.Organization
}

func DaysFromDuration(d time.Duration) string {
	days := d.Hours() / 24
	return strconv.Itoa(int(days))
}
