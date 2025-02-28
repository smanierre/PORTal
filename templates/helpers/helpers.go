package helpers

import (
	"PORTal/types"
	"encoding/json"
	"strconv"
	"time"
)

func DaysFromDuration(d time.Duration) string {
	days := d.Hours() / 24
	return strconv.Itoa(int(days))
}

func RequirementsToJSON(requirements []types.Requirement) string {
	if len(requirements) == 0 {
		return "[]"
	}
	b, err := json.Marshal(requirements)
	if err != nil {
		return "[]"
	}
	return string(b)
}
