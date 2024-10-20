package api_test

import (
	"PORTal/api"
	"PORTal/backend"
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
	"testing"
)

func TestAddMemberQualification(t *testing.T) {
	b := newMockBackend()
	b.assignMemberQualificationOverride = func(memberID, qualID string) error {
		if qualID == "good" && memberID == "good" {
			return nil
		}
		if memberID == "notfound" {
			return backend.ErrMemberNotFound
		}
		if qualID == "notfound" {
			return backend.ErrQualificationNotFound
		}
		if qualID == "alreadyAssigned" {
			return backend.ErrQualificationAlreadyAssigned
		}
		if memberID == "bad" {
			return errors.New("generic error")
		}
		return errors.New("unexpected case")
	}
	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name            string
		memberId        string
		qualificationID string
		statusCode      int
	}{
		{
			name:            "Successful add",
			memberId:        "good",
			qualificationID: "good",
			statusCode:      http.StatusOK,
		},
		{
			name:            "Member not found",
			memberId:        "notfound",
			qualificationID: "good",
			statusCode:      http.StatusNotFound,
		},
		{
			name:            "Qualification not found",
			memberId:        "good",
			qualificationID: "notfound",
			statusCode:      http.StatusNotFound,
		},
		{
			name:            "Qualification already assigned",
			memberId:        "good",
			qualificationID: "alreadyAssigned",
			statusCode:      http.StatusBadRequest,
		},
		{
			name:            "Backend error",
			memberId:        "bad",
			qualificationID: "good",
			statusCode:      http.StatusInternalServerError,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/member/%s/qualification/%s", tt.memberId, tt.qualificationID), nil)
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetMemberQualification(t *testing.T) {

}

func TestGetMemberQualifications(t *testing.T) {
	qual1 := testutils.RandomQualification()
	qual2 := testutils.RandomQualification()
	b := newMockBackend()
	b.getMemberQualificationsOverride = func(memberID string) ([]types.Qualification, error) {
		switch memberID {
		case "none":
			return []types.Qualification{}, nil
		case "one":
			return []types.Qualification{qual1}, nil
		case "two":
			return []types.Qualification{qual1, qual2}, nil
		case "bad":
			return nil, errors.New("generic error")
		case "notfound":
			return nil, backend.ErrMemberNotFound
		default:
			return nil, errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name             string
		memberID         string
		expectedResponse []types.Qualification
		statusCode       int
	}{
		{
			name:             "Empty response",
			memberID:         "none",
			expectedResponse: []types.Qualification{},
			statusCode:       http.StatusOK,
		},
		{
			name:             "Single item response",
			memberID:         "one",
			expectedResponse: []types.Qualification{qual1},
			statusCode:       http.StatusOK,
		},
		{
			name:             "Multi item response",
			memberID:         "two",
			expectedResponse: []types.Qualification{qual1, qual2},
			statusCode:       http.StatusOK,
		},
		{
			name:             "Backend error",
			memberID:         "bad",
			expectedResponse: nil,
			statusCode:       http.StatusInternalServerError,
		},
		{
			name:             "Member not found",
			memberID:         "notfound",
			expectedResponse: nil,
			statusCode:       http.StatusNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/member/%s/qualifications", tt.memberID), nil)
			s.ServeHTTP(w, r)

			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusOK {
				var res []types.Qualification
				err := json.NewDecoder(w.Body).Decode(&res)
				if err != nil {
					t.Errorf("Error deserializing response from server: %s", err.Error())
				}
				for _, mq := range tt.expectedResponse {
					found := false
					for _, mq2 := range res {
						if reflect.DeepEqual(mq, mq2) {
							found = true
						}
					}
					if !found {
						t.Errorf("Didn't find MemberQualifiation in response: %+v", mq)
					}
				}
			}
		})
	}
}

func TestDeleteMemberQualification(t *testing.T) {
	goodMemberID := uuid.NewString()
	goodQualID := uuid.NewString()
	b := newMockBackend()
	b.removeMemberQualificationOverride = func(memberID, qualID string) error {
		if qualID == goodQualID && memberID == goodMemberID {
			return nil
		}
		if memberID == "notfound" {
			return backend.ErrMemberNotFound
		}
		if qualID == "notfound" {
			return backend.ErrMemberQualificationNotFound
		}
		return errors.New("unexpected case")
	}

	s := api.New(slog.Default(), b, false, api.Config{JWTSecret: "test"})

	tc := []struct {
		name       string
		memberID   string
		qualID     string
		statusCode int
	}{
		{
			name:       "Successful delete",
			memberID:   goodMemberID,
			qualID:     goodQualID,
			statusCode: http.StatusOK,
		},
		{
			name:       "Member not found",
			memberID:   "notfound",
			qualID:     "good",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Qualification not found",
			memberID:   "good",
			qualID:     "notfound",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Backend error",
			memberID:   "bad",
			qualID:     "good",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/member/%s/qualification/%s", tt.memberID, tt.qualID), nil)
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
