package templates

import (
	"PORTal/types"
	"fmt"
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
