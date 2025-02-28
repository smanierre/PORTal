package server

import (
	"PORTal/types"
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

//go:embed assets
var assetsDir embed.FS

type ContextKey string

const (
	MemberContextKey  ContextKey = "member"
	SessionCookieName            = "session_id"
)

type Backend interface {
	AddMember(m types.Member) (types.Member, error)
	GetMember(identifier string) (types.Member, error)
	GetDisabledMember(identifier string) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	GetDisabledMembers() ([]types.Member, error)
	GetSubordinates(memberID string) ([]types.Member, error)
	UpdateMember(m types.Member) (types.Member, error)
	DeleteMember(id string) error
	DisableMember(id string) error
	EnableMember(id string) error

	AddQualification(q types.Qualification) (types.Qualification, error)
	GetQualification(id string) (types.Qualification, error)
	GetAllQualifications() ([]types.Qualification, error)
	UpdateQualification(q types.Qualification) (types.Qualification, error)
	DeleteQualification(id string) error

	AssignMemberQualification(memberID, qualID string) error
	GetMemberQualification(memberID string, qualificationID string) (types.Qualification, error)
	GetMemberQualifications(memberID string) ([]types.Qualification, error)
	RemoveMemberQualification(memberID, qualificationID string) error

	AddRequirement(r types.Requirement) (types.Requirement, error)
	GetRequirement(id string) (types.Requirement, error)
	GetAllRequirements() ([]types.Requirement, error)
	UpdateRequirement(r types.Requirement) (types.Requirement, error)
	DeleteRequirement(id string) error

	Login(username, password string) (types.Member, error)
	CreateSession(memberID, userAgent, ipAddress string) (string, time.Time)
	ValidateSession(sessionID, userAgent, ipAddress string) (types.Member, error)
	DeleteSession(sessionID string)
}

type Config struct {
	Domain       string `yaml:"domain"`
	Port         int    `yaml:"port"`
	Organization string `yaml:"organization"`
	Service      string `yaml:"service"`
}

type Server struct {
	logger  *slog.Logger
	backend Backend
	mux     *http.ServeMux
	dev     bool
	config  Config
}

func New(logger *slog.Logger, backend Backend, dev bool, config Config) Server {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Loading templates...")
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating new server")
	s := Server{
		logger:  logger,
		backend: backend,
		mux:     http.NewServeMux(),
		dev:     dev,
		config:  config,
	}

	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering routes...")
	// Static assets
	s.mux.Handle("GET /assets/", http.FileServerFS(assetsDir))

	// Index is the login page
	s.mux.Handle("GET /", http.HandlerFunc(s.LoginGetHandler))
	s.mux.Handle("POST /", http.HandlerFunc(s.LoginPostHandler))
	s.mux.Handle("GET /logout", http.HandlerFunc(s.LogoutHandler))

	// Dashboard
	s.mux.Handle("GET /dashboard", http.HandlerFunc(s.DashboardGetHandler))

	// Admin page
	s.mux.Handle("GET /admin", http.HandlerFunc(s.AdminGetHandler))

	// Admin Member Routes
	s.mux.Handle("GET /admin/members", http.HandlerFunc(s.AdminMembersPaneGetHandler))
	s.mux.Handle("GET /admin/members/disabled", http.HandlerFunc(s.AdminMembersPaneGetDisabledHandler))
	s.mux.Handle("GET /admin/members/{id}", http.HandlerFunc(s.AdminMemberEditorGetHandler))
	s.mux.Handle("POST /admin/members/{id}", http.HandlerFunc(s.AdminMemberUpdateHandler))
	s.mux.Handle("POST /admin/members/{id}/disable", http.HandlerFunc(s.AdminMemberDisableHandler))
	s.mux.Handle("POST /admin/members/{id}/enable", http.HandlerFunc(s.AdminMemberEnableHandler))
	s.mux.Handle("GET /admin/members/add", http.HandlerFunc(s.AdminNewMemberHandler))
	s.mux.Handle("POST /admin/members/add", http.HandlerFunc(s.AdminMemberAddHandler))
	s.mux.Handle("GET /admin/members/{id}/potentialSupervisors/{grade}", http.HandlerFunc(s.AdminGetPotentialSupervisorsHandler))

	// Admin Qualification Routes
	s.mux.Handle("GET /admin/qualifications", http.HandlerFunc(s.AdminQualificationsPaneGetHandler))
	s.mux.Handle("GET /admin/qualifications/{id}", http.HandlerFunc(s.AdminQualificationEditorGetHandler))
	s.mux.Handle("POST /admin/qualifications/{id}", http.HandlerFunc(s.AdminQualificationUpdateHandler))
	s.mux.Handle("GET /admin/qualifications/add", http.HandlerFunc(s.AdminNewQualificationHandler))
	s.mux.Handle("POST /admin/qualifications/add", http.HandlerFunc(s.AdminQualificationAddHandler))

	// Component Routes
	s.mux.Handle("GET /components/requirementItem", http.HandlerFunc(s.RequirementItemComponent))
	s.mux.Handle("GET /components/requirementEditor", http.HandlerFunc(s.RequirementEditorComponent))

	logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully registered routes")
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.assetMiddleware(
		s.skipLoginMiddleware(
			s.sessionRequiredMiddleware(
				s.adminRequiredMiddleware(
					s.mux)))).ServeHTTP(w, r)
}

func (s Server) makeCookie(name, value string, expiration time.Time) *http.Cookie {

	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   s.config.Domain,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func removeCookie(w http.ResponseWriter, name, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  domain,
		Expires: time.Now(),
		Path:    "/",
	})
}

func checkHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}

func getSessionId(r *http.Request) (string, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func getMemberFromContext(ctx context.Context) (types.Member, error) {
	mVal := ctx.Value(MemberContextKey)
	m, ok := mVal.(types.Member)
	if !ok {
		return types.Member{}, errors.New("Unable to get member from context")
	}
	return m, nil
}
