package api_test

import (
	"PORTal/api"
	"PORTal/backend"
	"PORTal/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddMember(t *testing.T) {
	b := newMockBackend()
	b.addMemberOverride = func(m types.Member) (types.Member, error) {
		if m.FirstName == "bad" {
			return types.Member{}, errors.New("error")
		} else if m.SupervisorID == "bad" {
			return types.Member{}, backend.ErrSupervisorNotFound
		}
		m.ID = uuid.NewString()
		return m, nil
	}

	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful create",
			body:       `{"first_name":"test","last_name":"member","rank":"TSgt","qualifications":null,"supervisor_id":"random"}`,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Missing first name",
			body:       `{"last_name":"member","rank":"TSgt","qualifications":null,"supervisor_id":"random"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Missing last name",
			body:       `{"first_name":"test","rank":"TSgt","qualifications":null,"supervisor_id":"random"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Missing rank",
			body:       `{"first_name":"test","last_name":"member","qualifications":null,"supervisor_id":"random"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Supervisor ID doesn't exist",
			body:       `{"first_name":"test","last_name":"member","rank":"TSgt","qualifications":null,"supervisor_id":"bad"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Malformed request",
			body:       `{"first_name":"test","last_name":"member","rank":"TSgt","qualifications":null,"supervisor_id":"random"`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Backend error",
			body:       "{\"first_name\":\"bad\",\"last_name\":\"member\",\"rank\":\"TSgt\",\"qualifications\":null,\"supervisor_id\":\"random\"}",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/member", strings.NewReader(tt.body))
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusCreated {
				var res types.Member
				err := json.NewDecoder(w.Body).Decode(&res)
				if err != nil {
					t.Errorf("Error when deserializing response into IdJson for TestAddMember: %s", err.Error())
				}
				if _, err := uuid.Parse(res.ID); err != nil {
					t.Errorf("Invalid UUID returned from server: %s", err.Error())
				}
			}
		})
	}
}

func TestGetMember(t *testing.T) {
	goodId := uuid.NewString()
	badId := uuid.NewString()
	notFoundId := uuid.NewString()
	testMember := types.Member{
		ApiMember: types.ApiMember{
			ID:           goodId,
			FirstName:    "test",
			LastName:     "member",
			Rank:         "TSgt",
			SupervisorID: "",
		},
		Password: "",
		Hash:     "",
	}
	b := newMockBackend()
	b.getMemberOverride = func(id string) (types.Member, error) {
		switch id {
		case goodId:
			return testMember, nil
		case badId:
			return types.Member{}, errors.New("generic error")
		case notFoundId:
			return types.Member{}, backend.ErrMemberNotFound
		default:
			return types.Member{}, errors.New("unexpected error")
		}

	}
	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name             string
		id               string
		expectedResponse types.Member
		statusCode       int
	}{
		{
			name:             "Successful get",
			id:               goodId,
			expectedResponse: testMember,
			statusCode:       http.StatusOK,
		},
		{
			name:             "Member not found",
			id:               notFoundId,
			expectedResponse: types.Member{},
			statusCode:       http.StatusNotFound,
		},
		{
			name:             "Backend error",
			id:               badId,
			expectedResponse: types.Member{},
			statusCode:       http.StatusInternalServerError,
		},
		{
			name:             "Malformed UUID",
			id:               uuid.NewString()[1:],
			expectedResponse: types.Member{},
			statusCode:       http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/member/%s", tt.id), nil)
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusOK {
				b := &bytes.Buffer{}
				json.NewEncoder(b).Encode(testMember.ToApiMember())
				if b.String() != w.Body.String() {
					t.Errorf("Expected response %s\nGot: %s", b.String(), w.Body.String())
				}
			}
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	testMember := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "test",
			LastName:     "member",
			Rank:         "TSgt",
			SupervisorID: "",
		},
		Password: "",
		Hash:     "",
	}
	testMember2 := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "test 2",
			LastName:     "member 2",
			Rank:         "MSgt",
			SupervisorID: uuid.NewString(),
		},
		Password: "",
		Hash:     "",
	}
	b := newMockBackend()
	shouldSucceed := false
	b.getAllMembersOverride = func() ([]types.Member, error) {
		if shouldSucceed {
			return []types.Member{testMember, testMember2}, nil
		} else {
			return nil, errors.New("generic error")
		}
	}
	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name            string
		expectedMembers []types.Member
		statusCode      int
		shouldSucceed   bool
	}{
		{
			name:            "Successful get",
			expectedMembers: []types.Member{testMember, testMember2},
			statusCode:      http.StatusOK,
			shouldSucceed:   true,
		},
		{
			name:            "Backend error",
			expectedMembers: nil,
			statusCode:      http.StatusInternalServerError,
			shouldSucceed:   false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/members", nil)
			shouldSucceed = tt.shouldSucceed
			s.ServeHTTP(w, r)
			if tt.statusCode != w.Code {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
			if shouldSucceed {
				var am []types.ApiMember
				b := &bytes.Buffer{}
				for _, m := range tt.expectedMembers {
					am = append(am, m.ToApiMember())
				}
				json.NewEncoder(b).Encode(am)
				if b.String() != w.Body.String() {
					t.Errorf("Expected response: %s\nGot: %s", b.String(), w.Body.String())
				}
			}
		})
	}
}

func TestUpdateMember(t *testing.T) {
	b := newMockBackend()
	b.updateMemberOverride = func(m types.Member) (types.Member, error) {
		if m.FirstName == "bad" {
			return types.Member{}, errors.New("generic error")
		} else if m.FirstName == "not found" {
			return types.Member{}, backend.ErrMemberNotFound
		} else if m.SupervisorID == "not found" {
			return types.Member{}, backend.ErrSupervisorNotFound
		}
		return m, nil
	}
	b.getMemberOverride = func(id string) (types.Member, error) {
		return types.Member{
			ApiMember: types.ApiMember{
				ID:           "old",
				FirstName:    "old",
				LastName:     "old",
				Rank:         "Old",
				SupervisorID: "Old",
			},
			Password: "Old",
			Hash:     "Old",
		}, nil
	}

	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful update",
			body:       `{"id":"old","first_name":"test new","last_name":"member new","rank":"SMSgt","qualifications":null,"supervisor_id":"random new"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Backend error",
			body:       `{"id":"old","first_name":"bad","last_name":"member new","rank":"SMSgt","qualifications":null,"supervisor_id":"random new"}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Bad JSON body",
			body:       `{"id":"old","first_name":"test new","last_name":"member new","rank":"SMSgt","qualifications":null,"supervisor_id":"random new"`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Member not found",
			body:       `{"id":"old","first_name":"not found","last_name":"member new","rank":"SMSgt","qualifications":null,"supervisor_id":"random new"}`,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Update ID",
			body:       `{"id":"new","first_name":"test new","last_name":"member new","rank":"SMSgt","qualifications":null,"supervisor_id":"random new"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Supervisor ID doesn't exist",
			body:       `{"id":"old","first_name":"test new","last_name":"member new","rank":"SMSgt","qualifications":null,"supervisor_id":"not found"}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/api/member/irrelevant", strings.NewReader(tt.body))
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code: %d, got: %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestDeleteMember(t *testing.T) {
	goodId := uuid.NewString()
	badId := uuid.NewString()
	b := newMockBackend()
	b.deleteMemberOverride = func(id string) error {
		switch id {
		case goodId:
			return nil
		case badId:
			return errors.New("generic error")
		default:
			t.Errorf("Unexpected case")
			return errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name       string
		id         string
		statusCode int
	}{
		{
			name:       "Successful delete",
			id:         goodId,
			statusCode: http.StatusOK,
		},
		{
			name:       "Backend error",
			id:         badId,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Invalid UUID",
			id:         goodId[1:],
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/member/%s", tt.id), nil)
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code: %d, got: %d", tt.statusCode, w.Code)
			}
		})
	}
}
