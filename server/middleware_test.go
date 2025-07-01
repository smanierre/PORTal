package server

import (
	"PORTal/backend"
	"PORTal/testutils"
	"PORTal/types"
	"bytes"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSkipLoginMiddleware_ShouldSkip(t *testing.T) {
	tc := []struct {
		Name   string
		Path   string
		Method string
	}{
		{
			Name:   "Non-Login Get",
			Path:   "/dashboard",
			Method: http.MethodGet,
		},
		{
			Name:   "Non-Login Post",
			Path:   "/dashboard",
			Method: http.MethodPost,
		},
		{
			Name:   "Login Post",
			Path:   "/",
			Method: http.MethodPost,
		},
	}

	// Setup handler for validation, since we don't expect the middleware to run, we can use the same validation funcs for all
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := MemberFromContext(r.Context())
		if err == nil {
			t.Error("Expected an error getting member from context, but got none")
		}
	})
	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		t.Errorf("This function should not be called!")
		return types.Member{}, nil
	}
	middleware := skipLoginMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), sessionStore)

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.Method, tt.Path, nil)
			r.AddCookie(&http.Cookie{
				Name:  SessionCookieName,
				Value: uuid.NewString(),
			})
			middleware.ServeHTTP(w, r)
		})
	}
}

func TestSkipLoginMiddleware_NoSession(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		t.Errorf("This function should not be called!")
		return types.Member{}, nil
	}
	// We expect the middleware to run, but not do anything since there isn't a session cookie
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := MemberFromContext(r.Context())
		if err == nil {
			t.Error("Expected an error getting member from context, but got none")
		}
	})
	middleware := skipLoginMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), sessionStore)
	middleware.ServeHTTP(w, r)
}

func TestSkipLoginMiddleware_InvalidSession(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: uuid.NewString(),
	})
	sessionStore := &testutils.MockSessionStore{}
	called := false
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		called = true
		return types.Member{}, backend.ErrSessionValidationFailed
	}
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := MemberFromContext(r.Context())
		if err == nil {
			t.Error("Expected an error getting member from context, but got none")
		}
		if !called {
			t.Error("Expected call to ValidateSession, but it was not called")
		}
	})
	middleware := skipLoginMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), sessionStore)
	middleware.ServeHTTP(w, r)
}

func TestSkipLoginMiddleware_ValidSession(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	sID := uuid.NewString()
	m := testutils.RandomMember(false)
	r.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: sID,
	})

	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		if sessionID == sID {
			return m, nil
		} else {
			t.Errorf("Expected sessionID: %s, got: %s", sID, sessionID)
			return types.Member{}, nil
		}
	}
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mem, err := MemberFromContext(r.Context())
		if err != nil {
			t.Errorf("Expected no error from MemberFromContext, but got: %s\n", err.Error())
		}
		if m != mem {
			t.Errorf("Expected member to be %v, got: %v\n", m, mem)
		}
	})

	middleware := skipLoginMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), sessionStore)
	middleware.ServeHTTP(w, r)
}

func TestSessionRequiredMiddleware_ShouldSkip(t *testing.T) {
	tc := []struct {
		Name   string
		Path   string
		Method string
	}{
		{
			Name:   "Login Get",
			Path:   "/",
			Method: http.MethodGet,
		},
		{
			Name:   "Login Post",
			Path:   "/",
			Method: http.MethodPost,
		},
	}
	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		t.Errorf("This function should not be called!")
		return types.Member{}, nil
	}
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := MemberFromContext(r.Context())
		if err == nil {
			t.Error("Expected an error getting member from context, but got none")
		}
	})
	middleware := sessionRequiredMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), "103d LRS", sessionStore)
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.Method, tt.Path, nil)
			middleware.ServeHTTP(w, r)
		})
	}
}

func TestSessionRequiredMiddleware_InvalidSessionNormal(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	r.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: uuid.NewString(),
	})

	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		return types.Member{}, backend.ErrSessionValidationFailed
	}
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This won't be called because the user will be redirected to the login
		t.Errorf("This handler shouldn't be called!")
	})
	middleware := sessionRequiredMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), "103d LRS", sessionStore)
	middleware.ServeHTTP(w, r)
	// Expect to be redirected to the login page
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got: %d", http.StatusFound, w.Code)
	}
	if w.Header().Get("Location") != "/" {
		t.Errorf("Expected location /, got %s", w.Header().Get("Location"))
	}
	// Expect set cookie header that will expire current session cookie
	if w.Header().Get("Set-Cookie") == "" {
		t.Error("Expected Set-Cookie header, but didn't find one")
	}
}

func TestSessionRequiredMiddleware_InvalidSessionHTMX(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	r.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: uuid.NewString(),
	})
	r.Header.Set("Hx-Request", "true")

	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		return types.Member{}, backend.ErrSessionValidationFailed
	}
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This won't be called because the user will be redirected to the login
		t.Errorf("This handler shouldn't be called!")
	})
	middleware := sessionRequiredMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), "103d LRS", sessionStore)
	middleware.ServeHTTP(w, r)
	// Expect to be rendered login page
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got: %d", http.StatusOK, w.Code)
	}
	if w.Header().Get("Hx-Push-Url") != "/" {
		t.Errorf("Expected header Hx-Push-Url /, got %s", w.Header().Get("Location"))
	}
	// Expect set cookie header that will expire current session cookie
	if w.Header().Get("Set-Cookie") == "" {
		t.Error("Expected Set-Cookie header, but didn't find one")
	}
}

func TestSessionRequiredMiddleware_ValidSession(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	sID := uuid.NewString()
	m := testutils.RandomMember(false)
	r.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: sID,
	})
	sessionStore := &testutils.MockSessionStore{}
	sessionStore.ValidateSessionOverride = func(sessionID, userAgent, ipAddress string) (types.Member, error) {
		if sessionID == sID {
			return m, nil
		} else {
			t.Errorf("Expected sessionID: %s, got: %s", sID, sessionID)
			return types.Member{}, nil
		}
	}
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mem, err := MemberFromContext(r.Context())
		if err != nil {
			t.Errorf("Expected no error from MemberFromContext, but got: %s\n", err.Error())
		}
		if m != mem {
			t.Errorf("Expected member to be %v, got: %v\n", m, mem)
		}
	})
	middleware := sessionRequiredMiddleware(validationHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)), "103d LRS", sessionStore)
	middleware.ServeHTTP(w, r)
}

func TestAdminRequiredMiddleware_ShouldSkip(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	// The successful skip, and successful validation both do nothing, so we can use a dummy
	// handler and verify by checking the logs
	validationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	b := &bytes.Buffer{}
	middleware := adminRequiredMiddleware(validationHandler, slog.New(slog.NewTextHandler(b, nil)))
	middleware.ServeHTTP(w, r)
	if strings.Contains(b.String(), "running middleware") {
		t.Error("Expected middleware to be skipped")
	}
}

func TestAdminRequiredMiddleware_NotAdmin(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	m := testutils.RandomMember(false)
	r = r.WithContext(context.WithValue(r.Context(), MemberContextKey, m))

	// This will only result in a redirect, so no verification needed in the handler
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := adminRequiredMiddleware(dummyHandler, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	middleware.ServeHTTP(w, r)
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got: %d", http.StatusFound, w.Code)
	}
	if w.Header().Get("Location") != "/dashboard" {
		t.Errorf("Expected location /dashboard, got %s", w.Header().Get("Location"))
	}
}

func TestAdminRequiredMiddleware_Admin(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	m := testutils.RandomMember(true)
	r = r.WithContext(context.WithValue(r.Context(), MemberContextKey, m))

	b := &bytes.Buffer{}
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := adminRequiredMiddleware(dummyHandler, slog.New(slog.NewTextHandler(b, nil)))
	middleware.ServeHTTP(w, r)
	if !strings.Contains(b.String(), "Member is admin") {
		t.Error("Expected string \"Member is admin\" in logs, but didn't find it")
	}
}
