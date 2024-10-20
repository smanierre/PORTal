package backend_test

import (
	"fmt"
	"time"
)

type expireClock struct{}

func (e expireClock) Now() time.Time {
	t, err := time.Parse(time.DateTime, "2000-01-01 01:00:00")
	if err != nil {
		panic(fmt.Sprintf("Error parsing time: %s", err.Error()))
	}
	return t
}
