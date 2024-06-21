package api

import (
	"PORTal/types"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func (s Server) addMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	err := s.backend.AddMemberQualification(qualID, memberID)
	if errors.Is(err, types.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, types.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, types.ErrQualificationAlreadyAssigned) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Server) getMemberQualifications(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	reqs, err := s.backend.GetMemberQualifications(memberID)
	if errors.Is(err, types.ErrMemberNotFound) {
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

func (s Server) updateMemberQualification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	qualID := r.PathValue("qualID")
	var mq types.MemberQualification
	err := json.NewDecoder(r.Body).Decode(&mq)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error deserializing request into MemberQualification struct", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id != mq.MemberID {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Member ID doesn't match ID in provided member qualification", slog.String("path_id", id), slog.String("provided_id", mq.MemberID))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if qualID != mq.Qualification.ID {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Qualification ID doesn't match ID in provided member qualification", slog.String("path_id", qualID), slog.String("provided_id", mq.Qualification.ID))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.backend.UpdateMemberQualification(mq)
	if errors.Is(err, types.ErrMemberNotFound) || errors.Is(err, types.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Server) deleteMemberQualification(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("id")
	qualID := r.PathValue("qualID")
	err := s.backend.DeleteMemberQualification(qualID, memberID)
	if errors.Is(err, types.ErrMemberNotFound) || errors.Is(err, types.ErrMemberQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
