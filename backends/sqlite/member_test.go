package sqlite_test

import (
	"PORTal/backends/sqlite"
	"PORTal/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func TestAddMember(t *testing.T) {
	id := uuid.NewString()
	supervisorId := uuid.NewString()
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

	tc := []struct {
		name          string
		member        types.Member
		expectedError error
	}{
		{
			name: "Successful no supervisor create",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             supervisorId,
					FirstName:      "Test",
					LastName:       "Member",
					Rank:           "TSgt",
					Qualifications: nil,
					SupervisorID:   "",
				},
				Password: "",
				Hash:     "",
			},
			expectedError: nil,
		},
		{
			name: "Successful create with supervisor",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             uuid.NewString(),
					FirstName:      "Test 2",
					LastName:       "Member 2",
					Rank:           "SMSgt",
					Qualifications: nil,
					SupervisorID:   supervisorId,
				},
				Password: "",
				Hash:     "",
			},
			expectedError: nil,
		},
		{
			name: "Supervisor doesn't exist",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             uuid.NewString(),
					FirstName:      "Test",
					LastName:       "Member",
					Rank:           "Amn",
					Qualifications: nil,
					SupervisorID:   uuid.NewString(),
				},
				Password: "",
				Hash:     "",
			},
			expectedError: types.ErrSupervisorNotFound,
		},
		{
			name:          "Empty Member",
			member:        types.Member{},
			expectedError: types.ErrMissingArgs,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.AddMember(tt.member)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestGetMember(t *testing.T) {
	testMember := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "Member",
			Rank:           "TSgt",
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "",
		Hash:     "",
	}
	testMemberWithSupervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "User",
			Rank:           "CMSgt",
			Qualifications: nil,
			SupervisorID:   testMember.ID,
		},
		Password: "",
		Hash:     "",
	}
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
	err = backend.AddMember(testMember)
	if err != nil {
		t.Fatalf("Error creating member for TestGetMember: %s", err.Error())
	}
	err = backend.AddMember(testMemberWithSupervisor)
	if err != nil {
		t.Fatalf("Error creating member for TestGetMember: %s", err.Error())
	}
	tc := []struct {
		name           string
		id             string
		expectedResult types.Member
		expectedError  error
	}{
		{
			name:           "Successful get",
			id:             testMember.ID,
			expectedResult: testMember,
			expectedError:  nil,
		},
		{
			name:           "Member not found",
			id:             uuid.NewString(),
			expectedResult: types.Member{},
			expectedError:  types.ErrMemberNotFound,
		},
		{
			name:           "Successful get with supervisor",
			id:             testMemberWithSupervisor.ID,
			expectedResult: testMemberWithSupervisor,
			expectedError:  nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			m, err := backend.GetMember(tt.id)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				if !reflect.DeepEqual(m, tt.expectedResult) {
					t.Errorf("Expected value: %+v\nGot: %+v", tt.expectedResult, m)
				}
			}
		})
	}
}

func TestGetAllMembers(t *testing.T) {
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
	testMember := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "Member",
			Rank:           "Amn",
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "",
		Hash:     "",
	}

	testMemberWithSupervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "User",
			Rank:           "CMSgt",
			Qualifications: nil,
			SupervisorID:   testMember.ID,
		},
		Password: "",
		Hash:     "",
	}

	tc := []struct {
		name           string
		expectedResult []types.Member
		expectedError  error
		setupFunc      func()
	}{
		{
			name:           "No result get",
			expectedResult: []types.Member{},
			expectedError:  nil,
			setupFunc:      func() {},
		},
		{
			name:           "Single item get without sup",
			expectedResult: []types.Member{testMember},
			expectedError:  nil,
			setupFunc: func() {
				err := backend.AddMember(testMember)
				if err != nil {
					t.Fatalf("Error adding member to database for TestGetAllMembers: %s", err.Error())
				}
			},
		},
		{
			name:           "Multi item get with sup",
			expectedResult: []types.Member{testMember, testMemberWithSupervisor},
			expectedError:  nil,
			setupFunc: func() {
				err := backend.AddMember(testMemberWithSupervisor)
				if err != nil {
					t.Fatalf("Error adding member to database for TestGetAllMembers: %s", err.Error())
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()
			members, err := backend.GetAllMembers()
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				for _, m := range tt.expectedResult {
					b := &bytes.Buffer{}
					json.NewEncoder(b).Encode(m)
					match := false
					for _, mm := range members {
						bb := &bytes.Buffer{}
						json.NewEncoder(bb).Encode(mm)
						if bb.String() == b.String() {
							match = true
						}
					}
					if !match {
						t.Errorf("Expected result: %+v\nGot: %+v", tt.expectedResult, members)
					}
				}
			}
		})
	}
}

func TestUpdateMember(t *testing.T) {
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
	originalMember := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Old",
			LastName:       "Old",
			Rank:           "Old",
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "old",
		Hash:     "old",
	}
	supervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Old",
			LastName:       "Old",
			Rank:           "Old",
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "old",
		Hash:     "old",
	}
	err = backend.AddMember(originalMember)
	if err != nil {
		t.Errorf("Error adding member to database for TestUpdateMember: %s", err.Error())
	}
	err = backend.AddMember(supervisor)
	if err != nil {
		t.Errorf("Error adding member to database for TestUpdateMember: %s", err.Error())
	}

	tc := []struct {
		name          string
		member        types.Member
		expectedError error
	}{
		{
			name: "Successful update without supervisor ID",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             originalMember.ID,
					FirstName:      "New",
					LastName:       "",
					Rank:           "",
					Qualifications: nil,
					SupervisorID:   "",
				},
				Password: "",
				Hash:     "",
			},
			expectedError: nil,
		},
		{
			name: "Successful update with supervisor ID",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             originalMember.ID,
					FirstName:      "New",
					LastName:       "New",
					Rank:           "New",
					Qualifications: nil,
					SupervisorID:   supervisor.ID,
				},
				Password: "",
				Hash:     "",
			},
			expectedError: nil,
		},
		{
			name: "Member not found",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             uuid.NewString(),
					FirstName:      "",
					LastName:       "",
					Rank:           "",
					Qualifications: nil,
					SupervisorID:   "",
				},
				Password: "",
				Hash:     "",
			},
			expectedError: types.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.UpdateMember(tt.member)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				m, err := backend.GetMember(tt.member.ID)
				if err != nil {
					t.Errorf("Error when getting member from database: %s", err.Error())
				}
				if !reflect.DeepEqual(tt.member, m) {
					t.Errorf("Expected member in database to be: %+v\nGot: %+v", tt.member, m)
				}
			}
		})
	}
}

func TestDeleteMember(t *testing.T) {
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
	noSupervisorMember := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "No",
			LastName:  "Supervisor",
			Rank:      types.E1,
		},
		Password: "",
		Hash:     "",
	}
	supervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Super",
			LastName:  "Visor",
			Rank:      types.E1,
		},
		Password: "",
		Hash:     "",
	}
	supervisorMember := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "Member With",
			LastName:     "Supervisor",
			Rank:         types.E1,
			SupervisorID: supervisor.ID,
		},
		Password: "",
		Hash:     "",
	}
	supervisorMember2 := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "2ND Member",
			LastName:     "With supervisor",
			Rank:         types.E9,
			SupervisorID: supervisor.ID,
		},
		Password: "",
		Hash:     "",
	}
	if err := backend.AddMember(noSupervisorMember); err != nil {
		t.Fatalf("Error adding member for TestDeleteMember: %s", err.Error())
	}
	if err := backend.AddMember(supervisor); err != nil {
		t.Fatalf("Error adding member for TestDeleteMember: %s", err.Error())
	}
	if err := backend.AddMember(supervisorMember); err != nil {
		t.Fatalf("Error adding member for TestDeleteMember: %s", err.Error())
	}
	if err := backend.AddMember(supervisorMember2); err != nil {
		t.Fatalf("Error adding member for TestDeleteMember: %s", err.Error())
	}

	tc := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "Successful delete without supervisor",
			id:            noSupervisorMember.ID,
			expectedError: nil,
		},
		{
			name:          "Successful delete with supervisor",
			id:            supervisorMember2.ID,
			expectedError: nil,
		},
		{
			name:          "Delete supervisor",
			id:            supervisor.ID,
			expectedError: nil,
		},
		{
			name:          "Member not found",
			id:            uuid.NewString(),
			expectedError: types.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.DeleteMember(tt.id)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				if _, err = backend.GetMember(tt.id); !errors.Is(err, types.ErrMemberNotFound) {
					t.Errorf("Expected error member not found, got: %s", err.Error())
				}
				if tt.id == supervisor.ID {
					m, err := backend.GetMember(supervisorMember.ID)
					if err != nil {
						t.Errorf("Expected member with supervisor to exist, but got error when retreiving: %s", err.Error())
					}
					if m.SupervisorID != "" {
						t.Errorf("Expected member with supervisor who got delete to have empty supervisor id, got: %s", err.Error())
					}
				}
			}
		})
	}
}
