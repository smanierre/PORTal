package server

import (
	"PORTal/templates"
	"PORTal/types"
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"os"
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
	UpdateMember(m types.Member, forceNoSupervisor, forceNoAdmin bool) (types.Member, error)
	DeleteMember(id string) error
	DisableMember(id string) error
	EnableMember(id string) error

	AddQualification(q types.Qualification) (types.Qualification, error)
	GetQualification(id string) (types.Qualification, error)
	GetAllQualifications() ([]types.Qualification, error)
	UpdateQualification(q types.Qualification, forceExpirationUpdate bool) (types.Qualification, error)
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

	AddReference(r types.Reference) (types.Reference, error)
	GetReference(id string) (types.Reference, error)
	GetReferences() ([]types.Reference, error)
	UpdateReference(reference types.Reference, overrideNoVolume bool) (types.Reference, error)
	DeleteReference(id string) error

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
	logger       *slog.Logger
	backend      Backend
	mux          *http.ServeMux
	dev          bool
	config       Config
	templateRepo *templates.TemplateRepo
}

func New(logger *slog.Logger, backend Backend, dev bool, config Config) Server {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Loading templates...")
	rootData := templates.RootData{
		Organization: config.Organization,
	}
	var t *templates.TemplateRepo
	if dev {
		t = templates.New(os.DirFS("templates"), rootData, config.Service)
	} else {
		t = templates.New(templates.TemplateDir, rootData, config.Service)
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating new server")
	s := Server{
		logger:       logger,
		backend:      backend,
		mux:          http.NewServeMux(),
		dev:          dev,
		config:       config,
		templateRepo: t,
	}

	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering routes...")
	// Static assets
	s.mux.Handle("GET /assets/", http.FileServerFS(assetsDir))

	// Index is the login page
	s.mux.Handle("GET /", http.HandlerFunc(s.LoginGetHandler))
	// Handle login requests
	s.mux.Handle("POST /", http.HandlerFunc(s.LoginPostHandler))
	// Logout request
	s.mux.Handle("GET /logout", http.HandlerFunc(s.LogoutHandler))

	// Dashboard
	s.mux.Handle("GET /dashboard", http.HandlerFunc(s.DashboardGetHandler))

	// Admin page
	s.mux.Handle("GET /admin", http.HandlerFunc(s.AdminGetHandler))
	s.mux.Handle("GET /admin/members", http.HandlerFunc(s.AdminMembersGetHandler))
	s.mux.Handle("GET /admin/members/disabled", http.HandlerFunc(s.AdminMembersGetDisabledHandler))
	s.mux.Handle("GET /admin/members/{id}", http.HandlerFunc(s.AdminMemberEditorGetHandler))
	s.mux.Handle("POST /admin/members/{id}", http.HandlerFunc(s.UpdateMember))
	s.mux.Handle("POST /admin/members/{id}/disable", http.HandlerFunc(s.AdminMemberDisableHandler))
	s.mux.Handle("POST /admin/members/{id}/enable", http.HandlerFunc(s.AdminMemberEnableHandler))
	s.mux.Handle("GET /admin/members/add", http.HandlerFunc(s.AdminGetNewMember))
	s.mux.Handle("POST /admin/members/add", http.HandlerFunc(s.AddMember))
	s.mux.Handle("GET /admin/members/{id}/potentialSupervisors/{grade}", http.HandlerFunc(s.AdminGetPotentialSupervisors))

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
