package templates

import (
	"PORTal/templates/components/webcomponents"
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

func MembersToListItems(m []types.Member) []webcomponents.SearchableListItem {
	items := make([]webcomponents.SearchableListItem, len(m))
	for i, member := range m {
		items[i] = member
	}
	return items
}

func QualificationsToListItems(q []types.Qualification) []webcomponents.SearchableListItem {
	items := make([]webcomponents.SearchableListItem, len(q))
	for i, qualification := range q {
		items[i] = qualification
	}
	return items
}

func GetOrg() string {
	return cfg.Organization
}
