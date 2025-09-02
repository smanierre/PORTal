package api

import (
	"encoding/json"
	"net/http"
)

type organizationResponse struct {
	Organization string `json:"organization"`
}

func (a api) GetApiOrganization(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(organizationResponse{Organization: a.organization})
}

func (a api) GetApiRanks(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(a.ranks)
}
