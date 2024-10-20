package api

import "PORTal/types"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ValidateLocalDataRequest struct {
	MemberId string `json:"member_id"`
}

type LoginResponse struct {
	Member         types.ApiMember       `json:"member"`
	Qualifications []types.Qualification `json:"qualifications"`
	Subordinates   []types.ApiMember     `json:"subordinates"`
}
