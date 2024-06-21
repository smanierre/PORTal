package api_test

import (
	"PORTal/api"
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

func TestAddMemberQualification(t *testing.T) {
	b := newMockBackend()
	b.addMemberQualificationOverride = func(qualID, memberID string) error {
		if qualID == "good" && memberID == "good" {
			return nil
		}
		if memberID == "notfound" {
			return types.ErrMemberNotFound
		}
		if qualID == "notfound" {
			return types.ErrQualificationNotFound
		}
		if qualID == "alreadyAssigned" {
			return types.ErrQualificationAlreadyAssigned
		}
		if memberID == "bad" {
			return errors.New("generic error")
		}
		return errors.New("unexpected case")
	}
	s := api.New(slog.Default(), b)

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

func TestGetMemberQualifications(t *testing.T) {
	qual := types.Qualification{
		ID:   uuid.NewString(),
		Name: "Test Qual",
	}
	mq1 := types.MemberQualification{
		MemberID:      "first",
		Qualification: qual,
		Active:        false,
		ActiveDate:    time.Time{},
	}
	mq2 := types.MemberQualification{
		MemberID:      "second",
		Qualification: qual,
		Active:        false,
		ActiveDate:    time.Time{},
	}
	b := newMockBackend()
	b.getMemberQualificationsOverride = func(memberID string) ([]types.MemberQualification, error) {
		switch memberID {
		case "none":
			return []types.MemberQualification{}, nil
		case "one":
			return []types.MemberQualification{mq1}, nil
		case "two":
			return []types.MemberQualification{mq1, mq2}, nil
		case "bad":
			return nil, errors.New("generic error")
		case "notfound":
			return nil, types.ErrMemberNotFound
		default:
			return nil, errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), b)

	tc := []struct {
		name             string
		memberID         string
		expectedResponse []types.MemberQualification
		statusCode       int
	}{
		{
			name:             "Empty response",
			memberID:         "none",
			expectedResponse: []types.MemberQualification{},
			statusCode:       http.StatusOK,
		},
		{
			name:             "Single item response",
			memberID:         "one",
			expectedResponse: []types.MemberQualification{mq1},
			statusCode:       http.StatusOK,
		},
		{
			name:             "Multi item response",
			memberID:         "two",
			expectedResponse: []types.MemberQualification{mq1, mq2},
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
				var res []types.MemberQualification
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

func TestUpdateMemberQualification(t *testing.T) {
	b := newMockBackend()
	b.updateMemberQualificationOverride = func(mq types.MemberQualification) error {
		switch mq.MemberID {
		case "good":
			return nil
		case "bad":
			return errors.New("generic error")
		case "notfound":
			return types.ErrMemberNotFound
		}
		switch mq.Qualification.ID {
		case "notfound":
			return types.ErrQualificationNotFound
		default:
			return errors.New("unexpected case")
		}
	}
	s := api.New(slog.Default(), b)

	tc := []struct {
		name            string
		body            string
		memberID        string
		qualificationID string
		statusCode      int
	}{
		{
			name:            "Successful update",
			memberID:        "good",
			qualificationID: "good",
			body:            `{"member_id":"good","qualification":{"id":"good","name":"test"}}`,
			statusCode:      http.StatusOK,
		},
		{
			name:            "Bad JSON body",
			body:            `{"member_id":"good","qualification":{"id":"good","name":"test"}`,
			memberID:        "good",
			qualificationID: "good",
			statusCode:      http.StatusBadRequest,
		},
		{
			name:            "Backend error",
			body:            `{"member_id":"bad","qualification":{"id":"good","name":"test"}}`,
			memberID:        "bad",
			qualificationID: "good",
			statusCode:      http.StatusInternalServerError,
		},
		{
			name:            "Member not found",
			body:            `{"member_id":"notfound","qualification":{"id":"good","name":"test"}}`,
			memberID:        "notfound",
			qualificationID: "good",
			statusCode:      http.StatusNotFound,
		},
		{
			name:            "Qualification not found",
			body:            `{"member_id":"neither","qualification":{"id":"notfound","name":"test"}}`,
			memberID:        "neither",
			qualificationID: "notfound",
			statusCode:      http.StatusNotFound,
		},
		{
			name:            "Mismatched member IDs",
			body:            `{"member_id":"bad","qualification":{"id":"good","name":"test"}}`,
			memberID:        "good",
			qualificationID: "good",
			statusCode:      http.StatusBadRequest,
		},
		{
			name:            "Mismatched qualification IDs",
			body:            `{"member_id":"bad","qualification":{"id":"good","name":"test"}}`,
			memberID:        "bad",
			qualificationID: "bad",
			statusCode:      http.StatusBadRequest,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/member/%s/qualification/%s", tt.memberID, tt.qualificationID), strings.NewReader(tt.body))
			s.ServeHTTP(w, r)
			if w.Code != tt.statusCode {
				t.Errorf("Expected response code: %d, got: %d", tt.statusCode, w.Code)
			}

		})
	}
}

func TestDeleteMemberQualification(t *testing.T) {
	goodMemberID := uuid.NewString()
	goodQualID := uuid.NewString()
	b := newMockBackend()
	b.deleteMemberQualificationOverride = func(qualID, memberID string) error {
		if qualID == goodQualID && memberID == goodMemberID {
			return nil
		}
		if memberID == "notfound" {
			return types.ErrMemberNotFound
		}
		if qualID == "notfound" {
			return types.ErrMemberQualificationNotFound
		}
		return errors.New("unexpected case")
	}

	s := api.New(slog.Default(), b)

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
