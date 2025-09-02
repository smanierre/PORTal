package api

import (
	"PORTal/backend"
	"PORTal/types"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func requirementToApiRequirement(r types.Requirement) Requirement {
	g := RequirementGrade(r.Grade)
	return Requirement{
		DaysValidFor:    &r.DaysValidFor,
		Grade:           &g,
		Id:              r.ID,
		Initial:         r.Initial,
		Name:            r.Name,
		Notes:           &r.Notes,
		QualificationId: &r.QualificationID,
		Reference:       r.Reference,
		Type:            string(r.Type),
	}
}

func requirementsToApiRequirements(reqs []types.Requirement) []Requirement {
	newReqs := []Requirement{}
	for _, r := range reqs {
		newReqs = append(newReqs, requirementToApiRequirement(r))
	}
	return newReqs
}

func apiRequirementToRequirement(req Requirement) types.Requirement {
	if req.Notes == nil {
		*req.Notes = ""
	}
	if req.QualificationId == nil {
		*req.QualificationId = ""
	}
	if req.Grade == nil {
		*req.Grade = ""
	}
	if req.DaysValidFor == nil {
		*req.DaysValidFor = 0
	}
	return types.Requirement{
		ID:              req.Id,
		Name:            req.Name,
		Initial:         req.Initial,
		Reference:       req.Reference,
		Notes:           *req.Notes,
		Type:            types.RequirementType(req.Type),
		QualificationID: *req.QualificationId,
		Grade:           types.Grade(*req.Grade),
		DaysValidFor:    *req.DaysValidFor,
	}
}

func (a api) PostApiRequirement(w http.ResponseWriter, r *http.Request) {
	var req Requirement
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "error decoding requirement from JSON", slog.String("err", err.Error()))
		r.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	newReq, err := a.qualificationStore.AddRequirement(apiRequirementToRequirement(req))
	if errors.Is(err, backend.ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(requirementToApiRequirement(newReq))
}

func (a api) PutApiRequirement(w http.ResponseWriter, r *http.Request) {
	var req Requirement
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "error decoding requirement from JSON", slog.String("err", err.Error()))
		r.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	updatedReq, err := a.qualificationStore.UpdateRequirement(apiRequirementToRequirement(req))
	if errors.Is(err, backend.ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(requirementToApiRequirement(updatedReq))
}
