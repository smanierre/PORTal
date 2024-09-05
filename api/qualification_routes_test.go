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

func TestAddQualification(t *testing.T) {
	b := newMockBackend()
	b.addQualificationOverride = func(q types.Qualification) (types.Qualification, error) {
		switch q.Name {
		case "test":
			return q, nil
		case "error":
			return types.Qualification{}, errors.New("generic error")
		default:
			t.Error("Unexpected case")
			return types.Qualification{}, errors.New("unexpected case")
		}
	}

	s := api.New(slog.Default(), b, false)

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful create",
			body:       `{"name":"test","notes":"test notes","expires":true,"expiration_days":10000}`,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Backend error",
			body:       `{"name":"error","notes":"test notes","expires":true,"expiration_days":10000}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Missing name",
			body:       `{"notes":"test notes","expires":true,"expiration_days":10000}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Missing expiration days with qual expiration",
			body:       `{"name":"test", "notes":"test notes","expires":true}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "No Expiration",
			body:       `{"name":"test", "notes":"test notes"}`,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Malformed JSON",
			body:       `{"name":"test", "notes":"test notes","expires":true`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/qualification", strings.NewReader(tt.body))

			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetQualification(t *testing.T) {
	goodId := uuid.NewString()
	notFoundId := uuid.NewString()
	badId := uuid.NewString()
	testQualification := types.Qualification{
		ID:                    goodId,
		Name:                  "test qual",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "these are some test notes",
		Expires:               true,
		ExpirationDays:        347,
	}
	b := newMockBackend()

	b.getQualificationOverride = func(id string) (types.Qualification, error) {
		switch id {
		case goodId:
			return testQualification, nil
		case notFoundId:
			return types.Qualification{}, backend.ErrQualificationNotFound
		case badId:
			return types.Qualification{}, errors.New("generic error")
		default:
			t.Errorf("unexpected case")
			return types.Qualification{}, errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), b, false)

	tc := []struct {
		name             string
		id               string
		statusCode       int
		expectedResponse types.Qualification
	}{
		{
			name:             "Successful get",
			id:               goodId,
			statusCode:       http.StatusOK,
			expectedResponse: testQualification,
		},
		{
			name:             "Qualification not found",
			id:               notFoundId,
			statusCode:       http.StatusNotFound,
			expectedResponse: types.Qualification{},
		},
		{
			name:             "Backend error",
			id:               badId,
			statusCode:       http.StatusInternalServerError,
			expectedResponse: types.Qualification{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/qualification/%s", tt.id), nil)
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusOK {
				b := &bytes.Buffer{}
				_ = json.NewEncoder(b).Encode(tt.expectedResponse)
				if b.String() != w.Body.String() {
					t.Errorf("Expected response: %s\nGot: %s", b.String(), w.Body.String())
				}
			}
		})
	}
}

func TestGetAllQualifications(t *testing.T) {
	testQualification := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "test qual",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "these are some test notes",
		Expires:               true,
		ExpirationDays:        347,
	}
	testQualification2 := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "test qual 2",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "these are some other test notes",
		Expires:               false,
		ExpirationDays:        0,
	}

	b := newMockBackend()
	shouldSucceed := false
	b.getAllQualificationsOverride = func() ([]types.Qualification, error) {
		if shouldSucceed {
			return []types.Qualification{testQualification, testQualification2}, nil
		}
		return nil, errors.New("generic error")
	}
	s := api.New(slog.Default(), b, false)

	tc := []struct {
		name             string
		statusCode       int
		expectedResponse []types.Qualification
		shouldSucceed    bool
	}{
		{
			name:             "Successful get",
			statusCode:       http.StatusOK,
			expectedResponse: []types.Qualification{testQualification, testQualification2},
			shouldSucceed:    true,
		},
		{
			name:             "Backend error",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: nil,
			shouldSucceed:    false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/qualifications", nil)
			shouldSucceed = tt.shouldSucceed
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
			if shouldSucceed {
				b := &bytes.Buffer{}
				json.NewEncoder(b).Encode(tt.expectedResponse)
				if b.String() != w.Body.String() {
					t.Errorf("Expected response: %s\nGot: %s", b.String(), w.Body.String())
				}
			}
		})
	}
}

func TestUpdateQualification(t *testing.T) {
	originalQualification := types.Qualification{
		ID:                    "old",
		Name:                  "old",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "old notes",
		Expires:               false,
		ExpirationDays:        0,
	}
	b := newMockBackend()
	b.updateQualificationOverride = func(q types.Qualification, forceUpdateExpiration bool) (types.Qualification, error) {
		switch q.Name {
		case "test":
			return q, nil
		case "bad":
			return types.Qualification{}, errors.New("generic error")
		case "not found":
			return types.Qualification{}, backend.ErrQualificationNotFound
		default:
			t.Error("unexpected case")
			return types.Qualification{}, errors.New("unexpected case")
		}
	}

	b.getRequirementOverride = func(id string) (types.Requirement, error) {
		switch id {
		case "not found":
			return types.Requirement{}, backend.ErrRequirementNotFound
		case "found":
			return types.Requirement{}, nil
		default:
			t.Error("unexpected case")
			return types.Requirement{}, errors.New("unexpected case")
		}
	}

	b.getQualificationOverride = func(id string) (types.Qualification, error) {
		if id == "not found" {
			return types.Qualification{}, backend.ErrQualificationNotFound
		}
		return originalQualification, nil
	}

	s := api.New(slog.Default(), b, false)

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful update",
			body:       `{"id": "old", "name":"test","notes":"test notes","initial_requirements":[{"id":"found"}],"recurring_requirements":[{"id":"found"}], "expires":true,"expiration_days":10000}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Backend error",
			body:       `{"id":"old", "name":"bad","notes":"test notes","expires":true,"expiration_days":10000}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "ID Update",
			body:       `{"id":"new", "name":"test","notes":"test notes","expires":true,"expiration_days":10000}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Initial Requirement doesn't exist",
			body:       `{"id":"old", "name":"test","notes":"test notes","initial_requirements": [{"id":"not found"}], "expires":true,"expiration_days":10000}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Recurring requirement doesn't exist",
			body:       `{"id":"old", "name":"test","notes":"test notes","recurring_requirements": [{"id":"not found"}], "expires":true,"expiration_days":10000}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Bad JSON request",
			body:       `{"name":"test","notes":"test notes","expires":true,"expiration_days":10000`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Qualification not found",
			body:       `{"id": "not found", "name":"not found","notes":"test notes","expires":true,"expiration_days":10000}`,
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/api/qualification/irrelevant", strings.NewReader(tt.body))
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestDeleteQualification(t *testing.T) {
	goodId := uuid.NewString()
	badId := uuid.NewString()
	notFoundId := uuid.NewString()
	b := newMockBackend()
	b.deleteQualificationOverride = func(id string) error {
		switch id {
		case goodId:
			return nil
		case badId:
			return errors.New("generic error")
		case notFoundId:
			return backend.ErrQualificationNotFound
		default:
			t.Errorf("unexpected case")
			return errors.New("unexpected case")
		}
	}

	s := api.New(slog.Default(), b, false)

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
			name:       "Qualification not found",
			id:         notFoundId,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Invalid UUID",
			id:         uuid.NewString()[1:],
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/qualification/%s", tt.id), nil)
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
