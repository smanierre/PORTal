package core

import "net/http"

type DashboardData struct{}

func (c Core) Dashboard(r *http.Request) (DashboardData, error) {
	return DashboardData{}, nil
}
