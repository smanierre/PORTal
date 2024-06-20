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
	"strings"
	"testing"
)

func TestAddRequirement(t *testing.T) {
	b := newMockBackend()
	b.addRequirementOverride = func(r types.Requirement) error {
		switch r.Name {
		case "test":
			return nil
		case "error":
			return errors.New("generic error")
		default:
			t.Error("unexpected case")
			return errors.New("unexpected case")
		}
	}

	s := api.New(slog.Default(), b)

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful create",
			body:       `{"name":"test", "notes":"test notes","description":"test description","days_valid_for": 365}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid JSON",
			body:       `{"name":"test", "notes":"test notes","description":"test description","days_valid_for": 365`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Backend error",
			body:       `{"name":"error", "notes":"test notes","description":"test description","days_valid_for": 365}`,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/requirement", strings.NewReader(tt.body))
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetRequirement(t *testing.T) {
	badId := uuid.NewString()
	notFoundId := uuid.NewString()

	testRequirement := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "test",
		Description:  "requirement",
		Notes:        "notes",
		DaysValidFor: 182,
	}
	b := newMockBackend()
	b.getRequirementOverride = func(id string) (types.Requirement, error) {
		switch id {
		case testRequirement.ID:
			return testRequirement, nil
		case notFoundId:
			return types.Requirement{}, types.ErrRequirementNotFound
		case badId:
			return types.Requirement{}, errors.New("generic error")
		default:
			t.Error("unexpected case")
			return types.Requirement{}, errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), b)

	tc := []struct {
		name             string
		id               string
		statusCode       int
		expectedResponse types.Requirement
	}{
		{
			name:             "Successful get",
			id:               testRequirement.ID,
			statusCode:       http.StatusOK,
			expectedResponse: testRequirement,
		},
		{
			name:             "Requirement not found",
			id:               notFoundId,
			statusCode:       http.StatusNotFound,
			expectedResponse: types.Requirement{},
		},
		{
			name:             "Backend error",
			id:               badId,
			statusCode:       http.StatusInternalServerError,
			expectedResponse: types.Requirement{},
		},
		{
			name:             "Invalid UUID",
			id:               uuid.NewString()[1:],
			statusCode:       http.StatusBadRequest,
			expectedResponse: types.Requirement{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/requirement/%s", tt.id), nil)
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusOK {
				b := &bytes.Buffer{}
				json.NewEncoder(b).Encode(tt.expectedResponse)
				if b.String() != w.Body.String() {
					t.Errorf("Expected response: %s\nGot: %s", b.String(), w.Body.String())
				}
			}
		})
	}
}

func TestGetAllRequirements(t *testing.T) {
	testRequirement := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test",
		Description:  "requirement",
		Notes:        "some test notes",
		DaysValidFor: -1,
	}
	testRequirement2 := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test 2",
		Description:  "requirement 2",
		Notes:        "some test notes 2",
		DaysValidFor: 400,
	}
	b := newMockBackend()
	shouldSucceed := false
	b.getAllRequirementsOverride = func() ([]types.Requirement, error) {
		if shouldSucceed {
			return []types.Requirement{testRequirement, testRequirement2}, nil
		} else {
			return nil, errors.New("generic error")
		}
	}
	s := api.New(slog.Default(), b)

	tc := []struct {
		name             string
		shouldSucceed    bool
		statusCode       int
		expectedResponse []types.Requirement
	}{
		{
			name:             "Successful get",
			shouldSucceed:    true,
			statusCode:       http.StatusOK,
			expectedResponse: []types.Requirement{testRequirement, testRequirement2},
		},
		{
			name:             "Backend error",
			shouldSucceed:    false,
			statusCode:       http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			shouldSucceed = tt.shouldSucceed
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/requirements", nil)
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusOK {
				b := &bytes.Buffer{}
				json.NewEncoder(b).Encode(tt.expectedResponse)
				if b.String() != w.Body.String() {
					t.Errorf("Expected response: %s\nGot: %s", b.String(), w.Body.String())
				}
			}
		})
	}
}

func TestUpdateRequirement(t *testing.T) {
	originalRequirement := types.Requirement{
		ID:           "old",
		Name:         "old",
		Description:  "old",
		Notes:        "old",
		DaysValidFor: 100,
	}
	b := newMockBackend()
	b.updateRequirementOverride = func(r types.Requirement) error {
		switch r.Name {
		case "error":
			return errors.New("generic error")
		default:
			return nil
		}
	}
	b.getRequirementOverride = func(id string) (types.Requirement, error) {
		switch id {
		case "not found":
			return types.Requirement{}, types.ErrRequirementNotFound
		default:
			return originalRequirement, nil
		}
	}

	s := api.New(slog.Default(), b)

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful update",
			body:       `{"id":"old", "name":"new", "notes":"new","description":"new","days_valid_for": 365}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid JSON",
			body:       `{"name":"new", "notes":"new","description":"new","days_valid_for": 365`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "ID Update",
			body:       `{"id":"new", "name":"new", "notes":"new","description":"new","days_valid_for": 365}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Requirement not found",
			body:       `{"id":"not found", "name":"new", "notes":"new","description":"new","days_valid_for": 365}`,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Backend error",
			body:       `{"id":"old", "name":"error", "notes":"new","description":"new","days_valid_for": 365}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Zero day expiration",
			body:       `{"id":"old", "name":"new", "notes":"new","description":"new","days_valid_for": 0}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/api/requirement/irrelevant", strings.NewReader(tt.body))
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestDeleteRequirement(t *testing.T) {
	goodId := uuid.NewString()
	badId := uuid.NewString()
	notFoundId := uuid.NewString()
	b := newMockBackend()
	b.deleteRequirementOverride = func(id string) error {
		switch id {
		case goodId:
			return nil
		case badId:
			return errors.New("generic error")
		case notFoundId:
			return types.ErrRequirementNotFound
		default:
			t.Error("unexpected case")
			return errors.New("unexpected case")
		}
	}

	s := api.New(slog.Default(), b)

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
			name:       "Not found",
			id:         notFoundId,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Backend error",
			id:         badId,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/requirement/%s", tt.id), nil)
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
