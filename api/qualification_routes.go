package api

import (
	"PORTal/backend"
	"PORTal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (s Server) addQualification(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	q := types.Qualification{}
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Error deserializing body into qualification struct", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = validateQualification(q); err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Incomplete qualification creation request", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := uuid.NewString()
	q.ID = id
	qual, err := s.backend.AddQualification(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(qual); err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing qualification to client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) getQualification(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	id := r.PathValue("id")
	q, err := s.backend.GetQualification(id)
	if errors.Is(err, backend.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(q)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing qualification to client", slog.String("error", err.Error()))
	}
}

func (s Server) getAllQualifications(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	quals, err := s.backend.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(quals)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing []types.Qualification to client", slog.String("error", err.Error()))
	}
}

func (s Server) updateQualification(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	var q types.Qualification
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error decoding qualification into struct", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//TODO: Implement authorization so that only the correct user is allowed to update an account
	existingQualification, err := s.backend.GetQualification(q.ID)
	if errors.Is(err, backend.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if existingQualification.ID != q.ID {
		l.LogAttrs(r.Context(), slog.LevelWarn, "User requesting to update qualification ID", slog.Any("update_request", q))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//TODO: fix this shit
	for _, req := range q.InitialRequirements {
		l.LogAttrs(r.Context(), slog.LevelInfo, "Verifying that all provided initial requirements exist...")
		_, err := s.backend.GetRequirement(req.ID)
		if errors.Is(err, backend.ErrRequirementNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	for _, req := range q.RecurringRequirements {
		l.LogAttrs(r.Context(), slog.LevelInfo, "Verifying that all provided recurring requirements exist...")
		_, err := s.backend.GetRequirement(req.ID)
		if errors.Is(err, backend.ErrRequirementNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	forceExpiration := q.Expires == false
	qualification, err := s.backend.UpdateQualification(q, forceExpiration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(qualification)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error serializing qualification to client", slog.String("error", err.Error()))
	}
}

func (s Server) deleteQualification(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Invalid UUID provided", slog.String("id", id))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := s.backend.DeleteQualification(id)
	if errors.Is(err, backend.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func validateQualification(q types.Qualification) error {
	errs := []string{}
	if q.Name == "" {
		errs = append(errs, "Name")
	}
	if q.Expires && q.ExpirationDays == 0 {
		errs = append(errs, "ExpirationDays")
	}
	if len(errs) > 0 {
		return fmt.Errorf("missing required values: %s", errs)
	}
	return nil
}
