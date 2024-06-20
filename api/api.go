package api

import (
	"PORTal/types"
	"context"
	"log/slog"
	"net/http"
)

type Backend interface {
	AddMember(m types.Member) error
	GetMember(id string) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	UpdateMember(m types.Member) error
	DeleteMember(id string) error
	AddQualification(q types.Qualification) error
	GetQualification(id string) (types.Qualification, error)
	GetAllQualifications() ([]types.Qualification, error)
	UpdateQualification(q types.Qualification) error
	DeleteQualification(id string) error
	AddRequirement(r types.Requirement) error
	GetRequirement(id string) (types.Requirement, error)
	GetAllRequirements() ([]types.Requirement, error)
	UpdateRequirement(r types.Requirement) error
	DeleteRequirement(id string) error
	AddMemberQualification(qualID, memberID string) error
	GetMemberQualifications(memberID string) ([]types.MemberQualification, error)
	UpdateMemberQualification(mq types.MemberQualification) error
	DeleteMemberQualification(qualID, memberID string) error
}

func New(logger *slog.Logger, backend Backend) Server {
	logger = logger.With(slog.String("source", "api_server"))
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating new api server")
	s := Server{
		logger:  logger,
		backend: backend,
		mux:     http.NewServeMux(),
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering routes...")

	// Member CRUD routes
	s.mux.Handle("POST /api/member", http.HandlerFunc(s.addMember))
	s.mux.Handle("GET /api/member/{id}", http.HandlerFunc(s.getMember))
	s.mux.Handle("GET /api/members", http.HandlerFunc(s.getAllMembers))
	s.mux.Handle("PUT /api/member/{id}", http.HandlerFunc(s.updateMember))
	s.mux.Handle("DELETE /api/member/{id}", http.HandlerFunc(s.deleteMember))

	// Qualification CRUD routes
	s.mux.Handle("POST /api/qualification", http.HandlerFunc(s.addQualification))
	s.mux.Handle("GET /api/qualification/{id}", http.HandlerFunc(s.getQualification))
	s.mux.Handle("GET /api/qualifications", http.HandlerFunc(s.getAllQualifications))
	s.mux.Handle("PUT /api/qualification/{id}", http.HandlerFunc(s.updateQualification))
	s.mux.Handle("DELETE /api/qualification/{id}", http.HandlerFunc(s.deleteQualification))

	// Requirement CRUD routes
	s.mux.Handle("POST /api/requirement", http.HandlerFunc(s.addRequirement))
	s.mux.Handle("GET /api/requirement/{id}", http.HandlerFunc(s.getRequirement))
	s.mux.Handle("GET /api/requirements", http.HandlerFunc(s.getAllRequirements))
	s.mux.Handle("PUT /api/requirement/{id}", http.HandlerFunc(s.updateRequirement))
	s.mux.Handle("DELETE /api/requirement/{id}", http.HandlerFunc(s.deleteRequirement))

	// Member-Qualification routes
	s.mux.Handle("POST /api/member/{id}/qualification/{qualID}", http.HandlerFunc(s.addMemberQualification))
	s.mux.Handle("DELETE /api/member/{id}/qualification/{qualID}", http.HandlerFunc(s.deleteMemberQualification))

	logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully registered routes")
	return s
}

type Server struct {
	logger  *slog.Logger
	backend Backend
	mux     *http.ServeMux
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
