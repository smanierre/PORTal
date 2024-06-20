package sqlite_test

import (
	"PORTal/backends/sqlite"
	"PORTal/types"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAddQualificationToMember(t *testing.T) {
	id := uuid.NewString()
	t.Cleanup(func() {
		err := os.Remove(fmt.Sprintf("test-%s.Db", id))
		if err != nil {
			t.Errorf("error cleaning up database: %s", err.Error())
		}
	})
	backend, err := sqlite.New(slog.Default(), fmt.Sprintf("test-%s.Db", id), 1)
	if err != nil {
		t.Fatalf("error creating new sqlite backend: %s", err.Error())
	}
	member := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "Member",
			Rank:           types.E9,
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "",
		Hash:     "",
	}
	qualification := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Test Qualification",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "Test notes",
		References:            nil,
		Expires:               false,
		ExpirationDays:        0,
	}
	if err := backend.AddMember(member); err != nil {
		t.Fatalf("Error adding member for TestAddQualificationToMember: %s", err.Error())
	}
	if err := backend.AddQualification(qualification); err != nil {
		t.Fatalf("Error adding qualification for TestAddQualificationToMember: %s", err.Error())
	}

	tc := []struct {
		name            string
		memberID        string
		qualificationID string
		expectedError   error
	}{
		{
			name:            "Successful add",
			memberID:        member.ID,
			qualificationID: qualification.ID,
			expectedError:   nil,
		},
		{
			name:            "Member doesn't exist",
			memberID:        uuid.NewString(),
			qualificationID: qualification.ID,
			expectedError:   types.ErrMemberNotFound,
		},
		{
			name:            "Qualification doesn't exist",
			memberID:        member.ID,
			qualificationID: uuid.NewString(),
			expectedError:   types.ErrQualificationNotFound,
		},
		{
			name:            "Add duplicate qualification",
			memberID:        member.ID,
			qualificationID: qualification.ID,
			expectedError:   types.ErrQualificationAlreadyAssigned,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.AddMemberQualification(tt.qualificationID, tt.memberID)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				quals, err := backend.GetMemberQualifications(tt.memberID)
				if err != nil {
					t.Errorf("Error getting qualifications for member in TestAddQualificationToMember: %s", err.Error())
				}
				for _, q := range quals {
					found := false
					if reflect.DeepEqual(q.Qualification, qualification) {
						found = true
					}
					if !found {
						t.Errorf("Expected to find qualification: %s, but didn't", qualification.Name)
					}
				}
			}
		})
	}
}

func TestGetQualificationsForMember(t *testing.T) {
	id := uuid.NewString()
	t.Cleanup(func() {
		err := os.Remove(fmt.Sprintf("test-%s.Db", id))
		if err != nil {
			t.Errorf("error cleaning up database: %s", err.Error())
		}
	})
	backend, err := sqlite.New(slog.Default(), fmt.Sprintf("test-%s.Db", id), 1)
	if err != nil {
		t.Fatalf("error creating new sqlite backend: %s", err.Error())
	}
	m := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Test",
			LastName:  "Member",
			Rank:      types.E1,
		},
	}
	if err = backend.AddMember(m); err != nil {
		t.Fatalf("Error adding member for TestGetQualificationsForMember: %s", err.Error())
	}
	qual := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Test",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "Test notes",
		References:            nil,
		Expires:               false,
		ExpirationDays:        0,
	}
	if err := backend.AddQualification(qual); err != nil {
		t.Fatalf("Error adding qualification for TestGetQualificationsForMember: %s", err.Error())
	}
	qual2 := types.Qualification{
		ID:    uuid.NewString(),
		Name:  "Test 2",
		Notes: "Test notes 2",
	}
	if err := backend.AddQualification(qual2); err != nil {
		t.Fatalf("Error adding qualification for TestGetQualificationsForMember: %s", err.Error())
	}
	mq1 := types.MemberQualification{
		MemberID:      m.ID,
		Qualification: qual,
		Active:        false,
		ActiveDate:    time.Time{},
	}
	mq2 := types.MemberQualification{
		MemberID:      m.ID,
		Qualification: qual2,
		Active:        false,
		ActiveDate:    time.Time{},
	}

	tc := []struct {
		name           string
		id             string
		expectedResult []types.MemberQualification
		expectedError  error
		setupFunc      func()
	}{
		{
			name:           "No qualifications",
			id:             m.ID,
			expectedResult: []types.MemberQualification{},
			expectedError:  nil,
			setupFunc:      func() {},
		},
		{
			name:           "Successful single qualification get",
			id:             m.ID,
			expectedResult: []types.MemberQualification{mq1},
			expectedError:  nil,
			setupFunc: func() {
				if err := backend.AddMemberQualification(qual.ID, m.ID); err != nil {
					t.Fatalf("Error adding qualification to member for TestGetQualificationsForMember: %s", err.Error())
				}
			},
		},
		{
			name:           "Successful multi qualification get",
			id:             m.ID,
			expectedResult: []types.MemberQualification{mq1, mq2},
			expectedError:  nil,
			setupFunc: func() {
				if err := backend.AddMemberQualification(qual2.ID, m.ID); err != nil {
					t.Fatalf("Error adding qualification to member for TestGetQualificationsForMember: %s", err.Error())
				}
			},
		},
		{
			name:           "User not found",
			id:             uuid.NewString(),
			expectedResult: []types.MemberQualification{},
			expectedError:  types.ErrMemberNotFound,
			setupFunc:      func() {},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()
			quals, err := backend.GetMemberQualifications(tt.id)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				for _, q := range tt.expectedResult {
					found := false
					for _, qu := range quals {
						if reflect.DeepEqual(q.Qualification, qu.Qualification) {
							found = true
						}
					}
					if !found {
						t.Errorf("expected to find qualification but didnt: %+v", q.Qualification)
					}
				}
			}
		})
	}
}

func TestUpdateMemberQualification(t *testing.T) {
	id := uuid.NewString()
	t.Cleanup(func() {
		err := os.Remove(fmt.Sprintf("test-%s.Db", id))
		if err != nil {
			t.Errorf("error cleaning up database: %s", err.Error())
		}
	})
	backend, err := sqlite.New(slog.Default(), fmt.Sprintf("test-%s.Db", id), 1)
	if err != nil {
		t.Fatalf("error creating new sqlite backend: %s", err.Error())
	}
	member := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Test",
			LastName:  "User",
			Rank:      types.E5,
		},
	}
	if err := backend.AddMember(member); err != nil {
		t.Fatalf("Error adding member for TestUpdateMemberQualification: %s", err.Error())
	}

	qualification := types.Qualification{
		ID:             uuid.NewString(),
		Name:           "Test qual",
		Notes:          "Test notes",
		Expires:        false,
		ExpirationDays: 0,
	}
	if err := backend.AddQualification(qualification); err != nil {
		t.Fatalf("Error adding member for TestUpdateMemberQualification: %s", err.Error())
	}

	if err := backend.AddMemberQualification(qualification.ID, member.ID); err != nil {
		t.Fatalf("Error adding qualification to member for TestUpdateMemberQualification: %s", err.Error())
	}
	tc := []struct {
		name          string
		update        types.MemberQualification
		expectedError error
	}{
		{
			name: "Successful update",
			update: types.MemberQualification{
				MemberID:      member.ID,
				Qualification: qualification,
				Active:        true,
				ActiveDate:    time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "Member doesn't exist",
			update: types.MemberQualification{
				MemberID:      uuid.NewString(),
				Qualification: qualification,
				Active:        true,
				ActiveDate:    time.Now().Add(1 * time.Hour),
			},
			expectedError: types.ErrMemberNotFound,
		},
		{
			name: "Qualification doesn't exist",
			update: types.MemberQualification{
				MemberID: member.ID,
				Qualification: types.Qualification{
					ID: uuid.NewString(),
				},
				Active:     false,
				ActiveDate: time.Time{},
			},
			expectedError: types.ErrQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.UpdateMemberQualification(tt.update)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				memberQuals, err := backend.GetMemberQualifications(tt.update.MemberID)
				if err != nil {
					t.Errorf("Error when getting qualifications for member: %s", err.Error())
				}
				found := false
				for _, mq := range memberQuals {
					if reflect.DeepEqual(mq.Qualification, tt.update.Qualification) &&
						mq.Active == tt.update.Active &&
						closeEnough(mq.ActiveDate, tt.update.ActiveDate) {
						found = true
					}
				}
				if !found {
					t.Errorf("Expected to find member qualification but got: %+v", memberQuals)
				}
			}
		})
	}
}

func TestDeleteMemberQualification(t *testing.T) {
	id := uuid.NewString()
	t.Cleanup(func() {
		err := os.Remove(fmt.Sprintf("test-%s.Db", id))
		if err != nil {
			t.Errorf("error cleaning up database: %s", err.Error())
		}
	})
	backend, err := sqlite.New(slog.Default(), fmt.Sprintf("test-%s.Db", id), 1)
	if err != nil {
		t.Fatalf("error creating new sqlite backend: %s", err.Error())
	}
	member := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Test",
			LastName:  "Member",
			Rank:      types.E7,
		},
	}
	if err := backend.AddMember(member); err != nil {
		t.Fatalf("Error adding member for TestDeleteMemberQualification: %s", err.Error())
	}
	qualification := types.Qualification{
		ID:   uuid.NewString(),
		Name: "Test Qual",
	}
	if err := backend.AddQualification(qualification); err != nil {
		t.Fatalf("Error adding qualification for TestDeleteMemberQualification: %s", err.Error())
	}
	if err := backend.AddMemberQualification(qualification.ID, member.ID); err != nil {
		t.Fatalf("Error adding qualification to member for TestDeleteMemberQualification: %s", err.Error())
	}
	qualification2 := types.Qualification{
		ID:   uuid.NewString(),
		Name: "Test Qual 2",
	}
	if err := backend.AddQualification(qualification2); err != nil {
		t.Fatalf("Error adding qualification for TestDeleteMemberQualification: %s", err.Error())
	}
	if err := backend.AddMemberQualification(qualification2.ID, member.ID); err != nil {
		t.Fatalf("Error adding qualification to member for TestDeleteMemberQualification: %s", err.Error())
	}
	notAssignedQualification := types.Qualification{
		ID:   uuid.NewString(),
		Name: "Not assigned",
	}
	if err := backend.AddQualification(notAssignedQualification); err != nil {
		t.Fatalf("Error adding qualification for TestDeleteMemberQualification: %s", err.Error())
	}

	tc := []struct {
		name            string
		memberID        string
		qualificationID string
		expectedError   error
	}{
		{
			name:            "Successful delete",
			memberID:        member.ID,
			qualificationID: qualification.ID,
			expectedError:   nil,
		},
		{
			name:            "Member qualification not found",
			memberID:        member.ID,
			qualificationID: notAssignedQualification.ID,
			expectedError:   types.ErrMemberQualificationNotFound,
		},
		{
			name:            "Qualification not found",
			memberID:        member.ID,
			qualificationID: uuid.NewString(),
			expectedError:   types.ErrQualificationNotFound,
		},
		{
			name:            "Member not found",
			memberID:        uuid.NewString(),
			qualificationID: qualification.ID,
			expectedError:   types.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.DeleteMemberQualification(tt.qualificationID, tt.memberID)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				memberQuals, err := backend.GetMemberQualifications(member.ID)
				if err != nil {
					t.Errorf("Error getting member qualifications: %s", err.Error())
				}
				for _, mq := range memberQuals {
					if reflect.DeepEqual(mq.Qualification, qualification) {
						t.Error("Expected to not find qualification on member, but did")
					}
				}
			}
		})
	}
}

func closeEnough(t1, t2 time.Time) bool {
	c := t1.Compare(t2)
	if c == 0 {
		return true
	}
	if c == -1 {
		if t2.Sub(t1) < 1*time.Millisecond {
			return true
		}
		return false
	} else {
		if t1.Sub(t2) < 1*time.Millisecond {
			return true
		}
		return false
	}
}
