package api

import (
	"PORTal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func (s Server) addRequirement(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	var req types.Requirement
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Invalid requirement JSON sent from client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err = s.backend.AddRequirement(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Server) getRequirement(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Invalid UUID supplied by client", slog.String("id", id))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requirement, err := s.backend.GetRequirement(id)
	if errors.Is(err, types.ErrRequirementNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(requirement)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing requirement to client", slog.String("error", err.Error()))
	}
}

func (s Server) getAllRequirements(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	reqs, err := s.backend.GetAllRequirements()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(reqs)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing slice of requirements to client", slog.String("error", err.Error()))
	}
}

func (s Server) updateRequirement(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	var req types.Requirement
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Invalid requirement JSON received from client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	originalRequirement, err := s.backend.GetRequirement(req.ID)
	if errors.Is(err, types.ErrRequirementNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = verifyRequirementUpdate(req, originalRequirement)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Errors when verifying requirement update", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.backend.UpdateRequirement(req)
	if errors.Is(err, types.ErrRequirementNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Server) deleteRequirement(w http.ResponseWriter, r *http.Request) {
	//l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	id := r.PathValue("id")
	err := s.backend.DeleteRequirement(id)
	if errors.Is(err, types.ErrRequirementNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func verifyRequirementUpdate(new, old types.Requirement) error {
	errs := []string{}
	if new.ID != old.ID {
		errs = append(errs, "attempting to update requirement ID")
	}
	if new.DaysValidFor == 0 {
		errs = append(errs, "requirement valid for days is equal to 0")
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors with requirement update: %s", errs)
	}
	return nil
}
