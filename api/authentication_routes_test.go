package api_test

import (
	"PORTal/api"
	"PORTal/types"
	"bytes"
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
	s := api.New(slog.Default(), m, false)

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
				var m types.Member
				err := json.NewDecoder(w.Body).Decode(&m)
				if err != nil {
					t.Errorf("Error decoding response into member struct: %s", err.Error())
				}
				if !reflect.DeepEqual(tt.expectedResponse, m.ToApiMember()) {
					t.Errorf("Expected response: %+v\nGot: %+v", tt.expectedResponse, m)
				}
				if w.Header().Get("Set-Cookie") == "" {
					t.Errorf("Expected Set-Cookie header, but didn't find one")
				}
			}
		})
	}
}

func TestValidateSession(t *testing.T) {
	validMemberID := uuid.NewString()
	m := newMockBackend()
	m.validateSessionOverride = func(sessionID, id, ipAddress string) error {
		switch sessionID {
		case "valid":
			break
		case "invalid":
			return errors.New("invalid session")
		default:
			t.Errorf("unexpected case")
			return errors.New("unexpected case")
		}
		switch id {
		case validMemberID:
			return nil
		default:
			t.Errorf("Unexpected case")
			return errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), m, false)

	tc := []struct {
		name          string
		sessionCookie *http.Cookie
		memberID      types.IdJson
		statusCode    int
	}{
		{
			name: "Successful Validation",
			sessionCookie: &http.Cookie{
				Name:  "session-id",
				Value: "valid",
			},
			memberID:   types.IdJson{ID: validMemberID},
			statusCode: http.StatusOK,
		},
		{
			name: "Bad cookie name",
			sessionCookie: &http.Cookie{
				Name:  "sessionid",
				Value: "irrelevant",
			},
			memberID:   types.IdJson{ID: validMemberID},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid Session",
			sessionCookie: &http.Cookie{
				Name:  "session-id",
				Value: "invalid",
			},
			memberID:   types.IdJson{ID: validMemberID},
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := json.NewEncoder(b).Encode(tt.memberID)
			if err != nil {
				t.Fatalf("Error encoding memberID to string: %s", err.Error())
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/validateSession", b)
			r.AddCookie(tt.sessionCookie)
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
