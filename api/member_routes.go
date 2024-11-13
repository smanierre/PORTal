package api

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (s Server) addMember(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	m := types.Member{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error deserializing body into member struct", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err = validateMember(m); err != nil {
		l.LogAttrs(r.Context(), slog.LevelInfo, "Incomplete create member request", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	insertedMember, err := s.backend.AddMember(m)
	if errors.Is(err, backend.ErrSupervisorNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(insertedMember)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing IdJson to client", slog.String("error", err.Error()))
	}
}

func (s Server) getLoggedInMember(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Checking for Identity cookie from initial page load")
	tokenCookie, err := r.Cookie(JWTCookieName)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error when retrieving Identity cookie", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := validateToken(tokenCookie.Value, s.jwtKeyFunc, s.logger)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	customClaims, ok := token.Claims.(*CustomClaims)
	if !ok {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error casting claims to CustomClaims")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	member, err := s.backend.GetMember(customClaims.Subject)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Identity validated, sending member back to client")
	err = json.NewEncoder(w).Encode(member.ToApiMember())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error encoding member to json response", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) getMember(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Invalid UUID passed to get member", slog.String("id", id))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := s.backend.GetMember(id)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(m.ToApiMember())
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error serializing ApiMember to client", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) getMemberSubordinates(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	members, err := s.backend.GetSubordinates(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	l.LogAttrs(r.Context(), slog.LevelInfo, "Converting members to ApiMember")
	var apiMembers []types.ApiMember
	for _, v := range members {
		apiMembers = append(apiMembers, v.ApiMember)
	}
	err = json.NewEncoder(w).Encode(apiMembers)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing slice of members to json", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s Server) getAllMembers(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	members, err := s.backend.GetAllMembers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var apiMembers []types.ApiMember
	for _, m := range members {
		apiMembers = append(apiMembers, m.ToApiMember())
	}
	err = json.NewEncoder(w).Encode(apiMembers)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error serializing list of ApiMember to client", slog.String("error", err.Error()))
	}
}

func (s Server) updateMember(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	m := types.Member{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		l.LogAttrs(r.Context(), slog.LevelError, "Error deserializing request body into member struct", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	existingMember, err := s.backend.GetMember(m.ID)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TODO: Implement authorization so that only the correct user is allowed to update an account
	if existingMember.ID != m.ID {
		l.LogAttrs(r.Context(), slog.LevelWarn, "User requesting to update ID", slog.Any("update_request", m))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	forceNoSupervisor := false
	if m.SupervisorID == "" {
		forceNoSupervisor = true
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Supervisor is blank in request, forcing removal")
	}
	member, err := s.backend.UpdateMember(m, forceNoSupervisor)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, backend.ErrSupervisorNotFound) {
		//w.WriteHeader(http.StatusBadRequest)
		//return
		//} else if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
	}
	err = json.NewEncoder(w).Encode(member.ApiMember)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error encoding member to client", slog.String("error", err.Error()))
	}
}

func (s Server) deleteMember(w http.ResponseWriter, r *http.Request) {
	l := s.logger.With(slog.String("path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		l.LogAttrs(r.Context(), slog.LevelWarn, "Invalid UUID provided", slog.String("id", id))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := s.backend.DeleteMember(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func validateMember(m types.Member) error {
	errs := []string{}
	if m.FirstName == "" {
		errs = append(errs, "FirstName")
	}
	if m.LastName == "" {
		errs = append(errs, "LastName")
	}
	if m.Grade == "" {
		errs = append(errs, "Rank")
	}
	if len(errs) > 0 {
		return fmt.Errorf("missing required values: %s", errs)
	}
	return nil
}
