package api

import (
	"PORTal/backend"
	"PORTal/types"
	"encoding/json"
	"errors"
	"net/http"
)

func imrToApiImr(imr types.InitialMemberRequirement) InitialMemberRequirement {
	i := InitialMemberRequirement{
		AssignedBy:    imr.AssignedBy,
		CompletedBy:   &imr.CompletedBy,
		CompletedDate: &imr.CompletedDate,
		MemberId:      imr.MemberID,
		RequirementId: imr.RequirementID,
	}
	if imr.CompletedDate.IsZero() {
		i.CompletedDate = nil
	}
	if imr.CompletedBy == "" {
		i.CompletedBy = nil
	}
	return i
}

func imrsToApiImrs(imrs []types.InitialMemberRequirement) []InitialMemberRequirement {
	i := []InitialMemberRequirement{}
	for _, imr := range imrs {
		i = append(i, imrToApiImr(imr))
	}
	return i
}

func rmrToApiRmr(rmr types.RecurringMemberRequirement) RecurringMemberRequirement {
	return RecurringMemberRequirement{
		AssignedBy:        rmr.AssignedBy,
		CompletionHistory: completionsToApiCompletions(rmr.CompletionHistory),
		Id:                rmr.ID,
		MemberId:          rmr.MemberID,
		RequirementId:     rmr.RequirementID,
	}
}

func rmrsToApiRmrs(rmrs []types.RecurringMemberRequirement) []RecurringMemberRequirement {
	r := []RecurringMemberRequirement{}
	for _, rmr := range rmrs {
		r = append(r, rmrToApiRmr(rmr))
	}
	return r
}

func completionToApiCompletion(c types.RecurringMemberRequirementCompletion) RecurringMemberRequirementCompletion {
	return RecurringMemberRequirementCompletion{
		CompletionDate: c.CompletionDate,
		Id:             c.ID,
		CompletedBy:    c.CompletedBy,
	}
}

func completionsToApiCompletions(cs []types.RecurringMemberRequirementCompletion) []RecurringMemberRequirementCompletion {
	r := []RecurringMemberRequirementCompletion{}
	for _, c := range cs {
		r = append(r, completionToApiCompletion(c))
	}
	return r
}

func (a api) GetApiMemberRequirementInitial(w http.ResponseWriter, r *http.Request, p GetApiMemberRequirementInitialParams) {
	ir, err := a.memberQualificationStore.GetInitialMemberRequirement(p.MemberId, p.RequirementId)
	if errors.Is(err, backend.ErrMemberNotFound) || errors.Is(err, backend.ErrRequirementNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(imrToApiImr(ir))
}

func (a api) GetApiMemberRequirementsInitial(w http.ResponseWriter, r *http.Request, p GetApiMemberRequirementsInitialParams) {
	irs, err := a.memberQualificationStore.GetInitialMemberRequirements(p.MemberId, p.QualificationId)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(imrsToApiImrs(irs))
}

func (a api) GetApiMemberRequirementsRecurring(w http.ResponseWriter, r *http.Request, p GetApiMemberRequirementsRecurringParams) {
	rrs, err := a.memberQualificationStore.GetRecurringMemberRequirements(p.MemberId, p.QualificationId)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(rmrsToApiRmrs(rrs))
}
