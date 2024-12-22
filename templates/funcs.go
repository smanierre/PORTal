package templates

import (
	"PORTal/types"
	"fmt"
)

func getDisplayNameFunc(serviceCode string) func(member types.Member) string {
	return func(m types.Member) string {
		return fmt.Sprintf("%s %s %s", m.GetRank(serviceCode), m.FirstName, m.LastName)
	}
}
