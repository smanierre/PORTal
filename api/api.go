package api

import (
	"PORTal/templates"
	"PORTal/types"
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Backend interface {
	AddMember(m types.Member) (types.Member, error)
	GetMember(identifier string) (types.Member, error)
	GetAllMembers() ([]types.Member, error)
	GetSubordinates(memberID string) ([]types.Member, error)
	UpdateMember(m types.Member, forceNoSupervisor bool) (types.Member, error)
	DeleteMember(id string) error

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
}

type Config struct {
	Domain       string `yaml:"domain"`
	Port         int    `yaml:"port"`
	Organization string `yaml:"organization"`
}

func New(logger *slog.Logger, backend Backend, dev bool, config Config) Server {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating new api server")
	s := Server{
		logger:  logger,
		backend: backend,
		mux:     http.NewServeMux(),
		dev:     dev,
		config:  config,
	}
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Loading templates...")
	rootData := templates.RootData{
		Organization: config.Organization,
	}
	var t *templates.TemplateRepo
	if s.dev {
		t = templates.New(os.DirFS("templates"), rootData)
	} else {
		t = templates.New(templates.TemplateDir, rootData)
	}

	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering routes...")

	s.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Render(w, "root", nil)
	}))
	//// Member CRUD routes
	//s.mux.Handle("POST /api/member", http.HandlerFunc(s.addMember))
	//s.mux.Handle("GET /api/member", http.HandlerFunc(s.getLoggedInMember))
	//s.mux.Handle("GET /api/member/{id}", http.HandlerFunc(s.getMember))
	//s.mux.Handle("GET /api/member/{id}/subordinates", http.HandlerFunc(s.getMemberSubordinates))
	//s.mux.Handle("GET /api/members", http.HandlerFunc(s.getAllMembers))
	//s.mux.Handle("PUT /api/member/{id}", http.HandlerFunc(s.updateMember))
	//s.mux.Handle("DELETE /api/member/{id}", http.HandlerFunc(s.deleteMember))
	//
	//// Qualification CRUD routes
	//s.mux.Handle("POST /api/qualification", http.HandlerFunc(s.addQualification))
	//s.mux.Handle("GET /api/qualification/{id}", http.HandlerFunc(s.getQualification))
	//s.mux.Handle("GET /api/qualifications", http.HandlerFunc(s.getAllQualifications))
	//s.mux.Handle("PUT /api/qualification/{id}", http.HandlerFunc(s.updateQualification))
	//s.mux.Handle("DELETE /api/qualification/{id}", http.HandlerFunc(s.deleteQualification))
	//
	//// Requirement CRUD routes
	//s.mux.Handle("POST /api/requirement", http.HandlerFunc(s.addRequirement))
	//s.mux.Handle("GET /api/requirement/{id}", http.HandlerFunc(s.getRequirement))
	//s.mux.Handle("GET /api/requirements", http.HandlerFunc(s.getAllRequirements))
	//s.mux.Handle("PUT /api/requirement/{id}", http.HandlerFunc(s.updateRequirement))
	//s.mux.Handle("DELETE /api/requirement/{id}", http.HandlerFunc(s.deleteRequirement))
	//
	//// Reference CRUD routes
	//s.mux.Handle("POST /api/reference", http.HandlerFunc(s.addReference))
	//s.mux.Handle("GET /api/reference/{id}", http.HandlerFunc(s.getReference))
	//s.mux.Handle("GET /api/references", http.HandlerFunc(s.getReferences))
	//
	//// Member-Qualification routes
	//s.mux.Handle("POST /api/member/{id}/qualification/{qualID}", http.HandlerFunc(s.assignMemberQualification))
	//s.mux.Handle("GET /api/member/{id}/qualifications", http.HandlerFunc(s.getMemberQualifications))
	//s.mux.Handle("GET /api/member/{id}/qualification/{qualID}", http.HandlerFunc(s.getMemberQualification))
	//s.mux.Handle("DELETE /api/member/{id}/qualification/{qualID}", http.HandlerFunc(s.removeMemberQualification))
	//
	//// Authentication routes
	//s.mux.Handle("POST /api/login", http.HandlerFunc(s.login))
	//s.mux.Handle("GET /api/logout", http.HandlerFunc(s.logout))
	//s.mux.Handle("GET /api/checkAdmin", http.HandlerFunc(s.checkAdmin))

	logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully registered routes")
	//if dev {
	//	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering frontend from build folder")
	//	s.mux.Handle("GET /", http.HandlerFunc(s.frontendHandler("ui/dist/")))
	//} else {
	//	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering frontend from /app/dist/")
	//	s.mux.Handle("GET /", http.HandlerFunc(s.frontendHandler("/app/dist/")))
	//}
	return s
}

type Server struct {
	logger  *slog.Logger
	backend Backend
	mux     *http.ServeMux
	dev     bool
	config  Config
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.dev {
		s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Development mode, setting CORS to http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	preflightMiddleware(s.mux).ServeHTTP(w, r)
}

func (s Server) makeCookie(name, value string) *http.Cookie {
	if s.dev {
		s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Returning dev cookie")
		return createDevCookie(name, value)
	} else {
		return &http.Cookie{}
		//return createCookie(name, value, s.config.Domain, time.Now().Add(s.config.JWTExpiration*time.Hour))
	}
}

func createCookie(name, value, domain string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/api",
		Domain:   domain,
		Expires:  expires,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}

func createDevCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/api",
		Domain:   "localhost",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s Server) removeCookie(w http.ResponseWriter, name string) {
	if s.dev {
		removeCookieDev(w, name)
	} else {
		removeCookie(w, name, s.config.Domain)
	}
}

func removeCookie(w http.ResponseWriter, name, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  domain,
		Expires: time.Now(),
		Path:    "/api",
	})
}

func removeCookieDev(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  "localhost",
		Expires: time.Now(),
		Path:    "/api",
	})
}
