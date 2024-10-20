package api_test

import (
	"PORTal/api"
	"PORTal/testutils"
	"PORTal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	member := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Test",
			LastName:  "Member",
			Username:  "tmember",
			Rank:      types.E8,
		},
		Password: "valid",
	}
	m := newMockBackend()
	m.loginOverride = func(username, password string) (types.Member, error) {
		if username == member.Username && password == member.Password {
			return member, nil
		}
		return types.Member{}, errors.New("generic error")
	}
	s := api.New(slog.Default(), m, false, api.Config{JWTSecret: "test", JWTExpiration: 1 * time.Hour})

	tc := []struct {
		name             string
		username         string
		password         string
		expectedResponse types.ApiMember
		statusCode       int
	}{
		{
			name:             "Successful login",
			username:         member.Username,
			password:         member.Password,
			expectedResponse: member.ToApiMember(),
			statusCode:       http.StatusOK,
		},
		{
			name:             "Invalid credentials",
			username:         "invalid",
			password:         "invalid",
			expectedResponse: types.ApiMember{},
			statusCode:       http.StatusUnauthorized,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, tt.username, tt.password)))
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code: %d, got: %d", tt.statusCode, w.Code)
			}
			if w.Code == http.StatusOK {
				var res api.LoginResponse
				err := json.NewDecoder(w.Body).Decode(&res)
				if err != nil {
					t.Errorf("Error decoding response into member struct: %s", err.Error())
				}
				if !reflect.DeepEqual(tt.expectedResponse, res.Member) {
					t.Errorf("Expected response: %+v\nGot: %+v", tt.expectedResponse, m)
				}
				setCookieHeader := w.Header().Get("Set-Cookie")
				if setCookieHeader == "" {
					t.Errorf("Expected Set-Cookie header to be set, but it wasn't")
				}
			}
		})
	}
}

func TestLogout(t *testing.T) {
	m := newMockBackend()
	s := api.New(slog.Default(), m, false, api.Config{
		JWTSecret: "supersecret",
	})

	member := testutils.RandomMember(false)
	member.ID = uuid.NewString()

	token, err := api.CreateToken(member, time.Hour, []byte("supersecret"))
	if err != nil {
		t.Fatalf("Error creating token for TestLogout: %s", err.Error())
	}

	tc := []struct {
		Name  string
		Token string
	}{
		{
			Name:  "Successful logout",
			Token: token,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/logout", nil)
			r.AddCookie(&http.Cookie{
				Name:    "identity",
				Value:   tt.Token,
				Path:    "/api",
				Expires: time.Now().Add(time.Hour),
			})
			s.ServeHTTP(w, r)
			if w.Header().Get("Set-Cookie") == "" {
				t.Errorf("Expected Set-Cookie header to be set, but it wasn't")
			}
			setCookieHeaderParts := strings.Split(w.Header().Get("Set-Cookie"), "=")
			ti, err := time.Parse(time.RFC1123, setCookieHeaderParts[3])
			if err != nil {
				t.Fatalf("Error parsing cookie expiration time: %s", err.Error())
			}
			if ti.After(time.Now()) {
				t.Errorf("Expected cookie to be expired but it isn't")
			}
		})
	}
}

func TestCheckAdmin(t *testing.T) {
	m := newMockBackend()
	s := api.New(slog.Default(), m, false, api.Config{
		JWTSecret: "supersecret",
	})

	adminMember := testutils.RandomMember(true)
	adminMember.ID = uuid.NewString()

	normalMember := testutils.RandomMember(false)
	normalMember.ID = uuid.NewString()

	adminToken, err := api.CreateToken(adminMember, time.Hour, []byte("supersecret"))
	if err != nil {
		t.Fatalf("Error creating adminToken for TestCheckAdmin: %s", err.Error())
	}

	normalToken, err := api.CreateToken(normalMember, time.Hour, []byte("supersecret"))
	if err != nil {
		t.Fatalf("Error creating normalToken for TestCheckAdmin: %s", err.Error())
	}

	invalidSignatureToken, err := api.CreateToken(adminMember, time.Hour, []byte("differentsecret"))
	if err != nil {
		t.Fatalf("Error creating invalidSignatureToken for TestCheckAdmin: %s", err.Error())
	}

	expiredToken, err := api.CreateToken(adminMember, time.Millisecond, []byte("supersecret"))
	if err != nil {
		t.Fatalf("Error creating expiredToken for TestCheckAdmin: %s", err.Error())
	}

	tc := []struct {
		Name           string
		Token          string
		ExpectedStatus int
	}{
		{
			Name:           "Successful validation",
			Token:          adminToken,
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "Not admin",
			Token:          normalToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name:           "Invalid Signature",
			Token:          invalidSignatureToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Name:           "Expired Token",
			Token:          expiredToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/checkAdmin", nil)
			r.AddCookie(&http.Cookie{
				Name:    api.JWTCookieName,
				Value:   tt.Token,
				Path:    "/api",
				Domain:  "localhost",
				Expires: time.Now().Add(time.Hour),
			})
			s.ServeHTTP(w, r)
			if w.Code != tt.ExpectedStatus {
				t.Errorf("Expected status code: %d, got: %d", tt.ExpectedStatus, w.Code)
			}
		})
	}
}
