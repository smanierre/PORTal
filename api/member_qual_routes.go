package api

import (
	"PORTal/backend"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func (s Server) assignMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	err := s.backend.AssignMemberQualification(memberID, qualID)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, backend.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, backend.ErrQualificationAlreadyAssigned) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Server) getMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	qual, err := s.backend.GetMemberQualification(memberID, qualID)
	if errors.Is(err, backend.ErrMemberNotFound) || errors.Is(err, backend.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(w).Encode(qual); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error serializing qualification to client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) getMemberQualifications(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	reqs, err := s.backend.GetMemberQualifications(memberID)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(reqs)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error serializing slice of MemberQualifications to client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) removeMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	err := s.backend.RemoveMemberQualification(memberID, qualID)
	if errors.Is(err, backend.ErrMemberNotFound) || errors.Is(err, backend.ErrMemberQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
