package api

import (
	"PORTal/backend"
	"PORTal/types"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func qualificationToApiQualification(q types.Qualification) Qualification {
	ir := requirementsToApiRequirements(q.InitialRequirements)
	rr := requirementsToApiRequirements(q.RecurringRequirements)
	return Qualification{
		Id:                    q.ID,
		InitialRequirements:   &ir,
		Name:                  q.Name,
		Notes:                 &q.Notes,
		RecurringRequirements: &rr,
	}
}

func qualificationsToApiQualifications(quals []types.Qualification) []Qualification {
	newQuals := []Qualification{}
	for _, q := range quals {
		newQuals = append(newQuals, qualificationToApiQualification(q))
	}
	return newQuals
}

func apiQualificationToQualification(q Qualification) types.Qualification {
	if q.Notes == nil {
		*q.Notes = ""
	}
	irs := []types.Requirement{}
	if q.InitialRequirements != nil {
		for _, v := range *q.InitialRequirements {
			irs = append(irs, apiRequirementToRequirement(v))
		}
	}
	rrs := []types.Requirement{}
	if q.RecurringRequirements != nil {
		for _, v := range *q.RecurringRequirements {
			rrs = append(rrs, apiRequirementToRequirement(v))
		}
	}
	return types.Qualification{
		ID:                    q.Id,
		Name:                  q.Name,
		InitialRequirements:   irs,
		RecurringRequirements: rrs,
		Notes:                 *q.Notes,
	}
}

func (a api) PostApiQualification(w http.ResponseWriter, r *http.Request) {
	var q Qualification
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "error decoding qualification from JSON", slog.String("err", err.Error()))
		r.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	qual, err := a.qualificationStore.AddQualification(apiQualificationToQualification(q))
	if errors.Is(err, backend.ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(qualificationToApiQualification(qual))
}

func (a api) PutApiQualification(w http.ResponseWriter, r *http.Request) {
	var q Qualification
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "error decoding qualification from JSON", slog.String("err", err.Error()))
		r.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}
	qual, err := a.qualificationStore.UpdateQualification(apiQualificationToQualification(q))
	if errors.Is(err, backend.ErrValidation) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{Message: err.Error()})
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(qualificationToApiQualification(qual))
}

func (a api) GetApiQualifications(w http.ResponseWriter, r *http.Request) {
	qualifications, err := a.qualificationStore.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(qualificationsToApiQualifications(qualifications))
}
