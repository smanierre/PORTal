package sqlite_test

import (
	"PORTal/backends/sqlite"
	"PORTal/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func init() {
	sqlite.BcryptCost = bcrypt.MinCost
}

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
					Username:       "Username1",
					Qualifications: nil,
					SupervisorID:   "",
				},
				Password: "password",
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
					Username:       "Username2",
					Qualifications: nil,
					SupervisorID:   supervisorId,
				},
				Password: "password",
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
					Username:       "Username3",
					Qualifications: nil,
					SupervisorID:   uuid.NewString(),
				},
				Password: "password",
				Hash:     "",
			},
			expectedError: types.ErrSupervisorNotFound,
		},
		{
			name:          "Empty Member",
			member:        types.Member{},
			expectedError: types.ErrMissingArgs,
		},
		{
			name: "Username not provided",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             supervisorId,
					FirstName:      "Test",
					LastName:       "Member",
					Rank:           "TSgt",
					Qualifications: nil,
					SupervisorID:   "",
				},
				Password: "password",
				Hash:     "",
			},
			expectedError: types.ErrMissingArgs,
		},
		{
			name: "Username taken",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:             uuid.NewString(),
					FirstName:      "Test",
					LastName:       "Member",
					Rank:           "TSgt",
					Username:       "Username1",
					Qualifications: nil,
					SupervisorID:   "",
				},
				Password: "password",
				Hash:     "",
			},
			expectedError: types.ErrUsernameAlreadyExists,
		},
		{
			name: "No Password",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:        uuid.NewString(),
					FirstName: "Test",
					LastName:  "Member",
					Username:  "Username4",
					Rank:      types.E7,
				},
				Password: "",
			},
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
			Username:       "user1",
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "password",
		Hash:     "",
	}
	testMemberWithSupervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "User",
			Rank:           "CMSgt",
			Username:       "user2",
			Qualifications: nil,
			SupervisorID:   testMember.ID,
		},
		Password: "password",
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
				if !reflect.DeepEqual(m.ApiMember, tt.expectedResult.ApiMember) {
					t.Errorf("Expected value: %+v\nGot: %+v", tt.expectedResult.ApiMember, m.ApiMember)
				}
			}
		})
	}
}

func TestGetMemberByUsername(t *testing.T) {
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
			LastName:  "User",
			Username:  "tuser",
			Rank:      types.E8,
		},
		Password: "password",
	}
	supervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "super",
			LastName:     "visor",
			Username:     "supervisor",
			Rank:         types.E1,
			SupervisorID: "",
		},
		Password: "password",
	}
	withSup := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "with",
			LastName:     "sup",
			Username:     "withsup",
			Rank:         types.E9,
			SupervisorID: supervisor.ID,
		},
		Password: "password",
	}
	if err := backend.AddMember(m); err != nil {
		t.Fatalf("Error inserting member into database for TestGetMemberByUsername: %s", err.Error())
	}
	if err := backend.AddMember(supervisor); err != nil {
		t.Fatalf("Error inserting member into database for TestGetMemberByUsername: %s", err.Error())
	}
	if err := backend.AddMember(withSup); err != nil {
		t.Fatalf("Error inserting member into database for TestGetMemberByUsername: %s", err.Error())
	}

	tc := []struct {
		name             string
		username         string
		expectedResponse types.Member
		expectedError    error
	}{
		{
			name:             "Successful Get without supervisor",
			username:         m.Username,
			expectedResponse: m,
			expectedError:    nil,
		},
		{
			name:             "Successful get with supervisor",
			username:         withSup.Username,
			expectedResponse: withSup,
			expectedError:    nil,
		},
		{
			name:             "Member not found",
			username:         "random",
			expectedResponse: types.Member{},
			expectedError:    types.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			member, err := backend.GetMemberByUsername(tt.username)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				if !reflect.DeepEqual(member.ApiMember, tt.expectedResponse.ApiMember) {
					t.Errorf("Expected result: %+v\nGot: %+v", tt.expectedResponse.ApiMember, member.ApiMember)
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
			Username:       "user1",
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "password",
		Hash:     "",
	}

	testMemberWithSupervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "User",
			Rank:           "CMSgt",
			Username:       "user2",
			Qualifications: nil,
			SupervisorID:   testMember.ID,
		},
		Password: "password",
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
					json.NewEncoder(b).Encode(m.ApiMember)
					match := false
					for _, mm := range members {
						bb := &bytes.Buffer{}
						json.NewEncoder(bb).Encode(mm.ApiMember)
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
			Username:       "Old",
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
			Username:       "Old2",
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
		{
			name: "Update username",
			member: types.Member{
				ApiMember: types.ApiMember{
					ID:        originalMember.ID,
					FirstName: "New",
					LastName:  "New",
					Username:  "Updated",
					Rank:      "New",
				},
			},
			expectedError: nil,
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
				if !reflect.DeepEqual(originalMember.MergeIn(tt.member).ApiMember, m.ApiMember) {
					t.Errorf("Expected member in database to be: %+v\nGot: %+v", originalMember.MergeIn(tt.member).ApiMember, m.ApiMember)
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
			Username:  "user1",
		},
		Password: "password",
		Hash:     "",
	}
	supervisor := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Super",
			LastName:  "Visor",
			Rank:      types.E1,
			Username:  "user2",
		},
		Password: "password",
		Hash:     "",
	}
	supervisorMember := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "Member With",
			LastName:     "Supervisor",
			Rank:         types.E1,
			Username:     "user3",
			SupervisorID: supervisor.ID,
		},
		Password: "password",
		Hash:     "",
	}
	supervisorMember2 := types.Member{
		ApiMember: types.ApiMember{
			ID:           uuid.NewString(),
			FirstName:    "2ND Member",
			LastName:     "With supervisor",
			Rank:         types.E9,
			Username:     "user4",
			SupervisorID: supervisor.ID,
		},
		Password: "password",
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
