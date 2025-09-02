package api

import (
	"PORTal/backend"
	"PORTal/types"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func apiMemberToMember(m Member) types.Member {
	return types.Member{
		ID:           m.Id,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		Username:     m.Username,
		Grade:        types.Grade(m.Grade),
		SupervisorID: *m.SupervisorId,
		Admin:        m.Admin,
		Disabled:     m.Disabled,
	}
}

func membersToApiMember(members []types.Member) []Member {
	newMems := []Member{}
	for _, m := range members {
		newMems = append(newMems, memberToApiMember(m))
	}
	return newMems
}

func memberToApiMember(m types.Member) Member {
	return Member{
		Admin:        m.Admin,
		Disabled:     m.Disabled,
		FirstName:    m.FirstName,
		Grade:        string(m.Grade),
		Id:           m.ID,
		LastName:     m.LastName,
		Username:     m.Username,
		SupervisorId: &m.SupervisorID,
	}
}

func (a api) GetApiMember(w http.ResponseWriter, r *http.Request, p GetApiMemberParams) {
	m, err := a.memberStore.GetMember(p.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(memberToApiMember(m))
}

// Requests for new member will include an initial password
type newMember struct {
	Member
	Password string `json:"password"`
}

func (a api) GetApiMembers(w http.ResponseWriter, r *http.Request, p GetApiMembersParams) {
	var mems []types.Member
	var err error
	if p.Filter == nil {
		mems, err = a.memberStore.GetAllMembers()
		disabledMems, nErr := a.memberStore.GetDisabledMembers()
		if nErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		mems = append(mems, disabledMems...)
	} else {
		switch *p.Filter {
		case Enabled:
			mems, err = a.memberStore.GetAllMembers()
		case Disabled:
			mems, err = a.memberStore.GetDisabledMembers()

		default:
			fmt.Println("we'll see")
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	apiMems := membersToApiMember(mems)
	_ = json.NewEncoder(w).Encode(apiMems)
}

func (a api) PostApiMember(w http.ResponseWriter, r *http.Request) {
	var m newMember
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "Error decoding member to JSON", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	tempMem := apiMemberToMember(m.Member)
	tempMem.Password = m.Password
	member, err := a.memberStore.AddMember(tempMem)
	if errors.Is(err, backend.ErrWeakPassword) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	if errors.Is(err, backend.ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(memberToApiMember(member))
}

func (a api) PutApiMember(w http.ResponseWriter, r *http.Request) {
	var m Member
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "Error decoding JSON", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	member, err := a.memberStore.UpdateMember(apiMemberToMember(m))
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(w)
	if errors.Is(err, backend.ErrWeakPassword) || errors.Is(err, backend.ErrPasswordTooLong) {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(Error{Message: err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(Error{Message: err.Error()})
		return
	}
	_ = encoder.Encode(memberToApiMember(member))
}

func (a api) PostApiDisableMember(w http.ResponseWriter, r *http.Request, p PostApiDisableMemberParams) {
	err := a.memberStore.DisableMember(p.Id)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a api) PostApiEnableMember(w http.ResponseWriter, r *http.Request, p PostApiEnableMemberParams) {
	err := a.memberStore.EnableMember(p.Id)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a api) GetApiPotentialSupervisors(w http.ResponseWriter, r *http.Request, p GetApiPotentialSupervisorsParams) {
	if p.Id == nil {
		*p.Id = ""
	}
	potentialSupervisors, err := a.memberStore.GetPotentialSupervisors(*p.Id, types.Grade(p.Grade))
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(w)
	if err != nil && strings.Contains(err.Error(), "invalid grade") {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(Error{Message: err.Error()})
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = encoder.Encode(membersToApiMember(potentialSupervisors))
}

func (a api) GetApiMemberSubordinates(w http.ResponseWriter, r *http.Request, p GetApiMemberSubordinatesParams) {
	subordinates := a.memberStore.GetSubordinates(p.MemberId)
	_ = json.NewEncoder(w).Encode(membersToApiMember(subordinates))
}
