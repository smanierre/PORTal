package server_test

import (
	"PORTal/backend"
	"PORTal/server"
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

func TestReferenceQualification(t *testing.T) {
	b := newMockBackend()
	b.addReferenceOverride = func(r types.Reference) (types.Reference, error) {
		switch r.Name {
		case "test":
			return r, nil
		case "error":
			return types.Reference{}, errors.New("generic error")
		default:
			t.Error("Unexpected case")
			return types.Reference{}, errors.New("unexpected case")
		}
	}

	s := server.New(slog.Default(), b, false, server.Config{JWTSecret: "test"})

	tc := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "Successful create",
			body:       `{"name":"test","volume": 3,"paragraph":"test paragraph"}`,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Backend error",
			body:       `{"name":"error","volume": 0, "paragraph":"paragraph"}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Bad request",
			body:       `{"name":"test", "volume": 0,"paragraph": "paragraph"`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/reference", strings.NewReader(tt.body))

			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetReference(t *testing.T) {
	goodId := uuid.NewString()
	notFoundId := uuid.NewString()
	badId := uuid.NewString()
	testReference := types.Reference{
		ID:        goodId,
		Name:      "test reference",
		Paragraph: "paragraph",
		Volume:    0,
	}
	b := newMockBackend()

	b.getReferenceOverride = func(id string) (types.Reference, error) {
		switch id {
		case goodId:
			return testReference, nil
		case notFoundId:
			return types.Reference{}, backend.ErrReferenceNotFound
		case badId:
			return types.Reference{}, errors.New("generic error")
		default:
			t.Errorf("unexpected case")
			return types.Reference{}, errors.New("unexpected case")
		}
	}
	s := server.New(slog.Default(), b, false, server.Config{JWTSecret: "test"})

	tc := []struct {
		name             string
		id               string
		statusCode       int
		expectedResponse types.Reference
	}{
		{
			name:             "Successful get",
			id:               goodId,
			statusCode:       http.StatusOK,
			expectedResponse: testReference,
		},
		{
			name:             "Qualification not found",
			id:               notFoundId,
			statusCode:       http.StatusNotFound,
			expectedResponse: types.Reference{},
		},
		{
			name:             "Backend error",
			id:               badId,
			statusCode:       http.StatusInternalServerError,
			expectedResponse: types.Reference{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/reference/%s", tt.id), nil)
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

func TestGetAllReferences(t *testing.T) {
	testReference := types.Reference{
		ID:        uuid.NewString(),
		Name:      "test reference",
		Volume:    0,
		Paragraph: "Paragraph 2",
	}
	testReference2 := types.Reference{
		ID:        uuid.NewString(),
		Name:      "test qual 2",
		Volume:    1,
		Paragraph: "10",
	}

	b := newMockBackend()
	shouldSucceed := false
	noItems := false
	b.getReferencesOverride = func() ([]types.Reference, error) {
		if shouldSucceed && !noItems {
			return []types.Reference{testReference, testReference2}, nil
		}
		if shouldSucceed && noItems {
			var res []types.Reference
			return res, nil
		}
		return nil, errors.New("generic error")
	}
	s := server.New(slog.Default(), b, false, server.Config{JWTSecret: "test"})

	tc := []struct {
		name             string
		statusCode       int
		expectedResponse []types.Reference
		shouldSucceed    bool
		noItems          bool
	}{
		{
			name:             "Successful get",
			statusCode:       http.StatusOK,
			expectedResponse: []types.Reference{testReference, testReference2},
			shouldSucceed:    true,
		},
		{
			name:             "Backend error",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: nil,
			shouldSucceed:    false,
		},
		{
			name:             "No References",
			statusCode:       http.StatusOK,
			expectedResponse: []types.Reference{},
			shouldSucceed:    true,
			noItems:          true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/references", nil)
			shouldSucceed = tt.shouldSucceed
			noItems = tt.noItems
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

//
//func TestUpdateQualification(t *testing.T) {
//	originalQualification := types.Qualification{
//		ID:                    "old",
//		Name:                  "old",
//		InitialRequirements:   nil,
//		RecurringRequirements: nil,
//		Notes:                 "old notes",
//		Expires:               false,
//		ExpirationInterval:        0,
//	}
//	b := newMockBackend()
//	b.updateQualificationOverride = func(q types.Qualification, forceUpdateExpiration bool) (types.Qualification, error) {
//		switch q.Name {
//		case "test":
//			return q, nil
//		case "bad":
//			return types.Qualification{}, errors.New("generic error")
//		case "not found":
//			return types.Qualification{}, backend.ErrQualificationNotFound
//		default:
//			t.Error("unexpected case")
//			return types.Qualification{}, errors.New("unexpected case")
//		}
//	}
//
//	b.getRequirementOverride = func(id string) (types.Requirement, error) {
//		switch id {
//		case "not found":
//			return types.Requirement{}, backend.ErrRequirementNotFound
//		case "found":
//			return types.Requirement{}, nil
//		default:
//			t.Error("unexpected case")
//			return types.Requirement{}, errors.New("unexpected case")
//		}
//	}
//
//	b.getQualificationOverride = func(id string) (types.Qualification, error) {
//		if id == "not found" {
//			return types.Qualification{}, backend.ErrQualificationNotFound
//		}
//		return originalQualification, nil
//	}
//
//	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})
//
//	tc := []struct {
//		name       string
//		body       string
//		statusCode int
//	}{
//		{
//			name:       "Successful update",
//			body:       `{"id": "old", "name":"test","notes":"test notes","initial_requirements":[{"id":"found"}],"recurring_requirements":[{"id":"found"}], "expires":true,"expiration_days":10000}`,
//			statusCode: http.StatusOK,
//		},
//		{
//			name:       "Backend error",
//			body:       `{"id":"old", "name":"bad","notes":"test notes","expires":true,"expiration_days":10000}`,
//			statusCode: http.StatusInternalServerError,
//		},
//		{
//			name:       "ID Update",
//			body:       `{"id":"new", "name":"test","notes":"test notes","expires":true,"expiration_days":10000}`,
//			statusCode: http.StatusBadRequest,
//		},
//		{
//			name:       "Initial Requirement doesn't exist",
//			body:       `{"id":"old", "name":"test","notes":"test notes","initial_requirements": [{"id":"not found"}], "expires":true,"expiration_days":10000}`,
//			statusCode: http.StatusBadRequest,
//		},
//		{
//			name:       "Recurring requirement doesn't exist",
//			body:       `{"id":"old", "name":"test","notes":"test notes","recurring_requirements": [{"id":"not found"}], "expires":true,"expiration_days":10000}`,
//			statusCode: http.StatusBadRequest,
//		},
//		{
//			name:       "Bad JSON request",
//			body:       `{"name":"test","notes":"test notes","expires":true,"expiration_days":10000`,
//			statusCode: http.StatusBadRequest,
//		},
//		{
//			name:       "Qualification not found",
//			body:       `{"id": "not found", "name":"not found","notes":"test notes","expires":true,"expiration_days":10000}`,
//			statusCode: http.StatusNotFound,
//		},
//	}
//
//	for _, tt := range tc {
//		t.Run(tt.name, func(t *testing.T) {
//			w := httptest.NewRecorder()
//			r := httptest.NewRequest(http.MethodPut, "/api/qualification/irrelevant", strings.NewReader(tt.body))
//			s.ServeHTTP(w, r)
//			if w.Code != tt.statusCode {
//				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
//			}
//		})
//	}
//}
//
//func TestDeleteQualification(t *testing.T) {
//	goodId := uuid.NewString()
//	badId := uuid.NewString()
//	notFoundId := uuid.NewString()
//	b := newMockBackend()
//	b.deleteQualificationOverride = func(id string) error {
//		switch id {
//		case goodId:
//			return nil
//		case badId:
//			return errors.New("generic error")
//		case notFoundId:
//			return backend.ErrQualificationNotFound
//		default:
//			t.Errorf("unexpected case")
//			return errors.New("unexpected case")
//		}
//	}
//
//	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})
//
//	tc := []struct {
//		name       string
//		id         string
//		statusCode int
//	}{
//		{
//			name:       "Successful delete",
//			id:         goodId,
//			statusCode: http.StatusOK,
//		},
//		{
//			name:       "Backend error",
//			id:         badId,
//			statusCode: http.StatusInternalServerError,
//		},
//		{
//			name:       "Qualification not found",
//			id:         notFoundId,
//			statusCode: http.StatusNotFound,
//		},
//		{
//			name:       "Invalid UUID",
//			id:         uuid.NewString()[1:],
//			statusCode: http.StatusBadRequest,
//		},
//	}
//
//	for _, tt := range tc {
//		t.Run(tt.name, func(t *testing.T) {
//			w := httptest.NewRecorder()
//			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/qualification/%s", tt.id), nil)
//			s.ServeHTTP(w, r)
//
//			if w.Code != tt.statusCode {
//				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
//			}
//		})
//	}
//}
