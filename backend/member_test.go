package backend_test

import (
	"PORTal/backend"
	"PORTal/providers/sqlite"
	"PORTal/testutils"
	"PORTal/types"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"sort"
	"testing"
)

func TestAddAndGetMember(t *testing.T) {
	dbID := uuid.NewString()
	t.Cleanup(func() {
		os.Remove(fmt.Sprintf("%s.db", dbID))
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sqlite.New(logger, fmt.Sprintf("%s.db", dbID), 1.0)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := backend.New(logger, provider, provider, provider, provider, &backend.Options{BcryptCost: bcrypt.MinCost})

	supervisor, err := b.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestAddMember_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name          string
		FirstName     string
		LastName      string
		UserName      string
		Rank          types.Rank
		SupervisorID  string
		Password      string
		ExpectedError error
	}{
		{
			Name:          "Successful Add",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      "username",
			Rank:          types.E4,
			SupervisorID:  "",
			Password:      testutils.RandomString(),
			ExpectedError: nil,
		},
		{
			Name:          "Member with supervisor",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Rank:          types.E8,
			SupervisorID:  supervisor.ID,
			Password:      testutils.RandomString(),
			ExpectedError: nil,
		},
		{
			Name:          "Supervisor doesn't exist",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Rank:          types.E4,
			SupervisorID:  uuid.NewString(),
			Password:      testutils.RandomString(),
			ExpectedError: backend.ErrSupervisorNotFound,
		},
		{
			Name:          "Duplicate username",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      "username",
			Rank:          types.E8,
			SupervisorID:  "",
			Password:      testutils.RandomString(),
			ExpectedError: backend.ErrDuplicateUsername,
		},
		{
			Name:          "Missing fields",
			FirstName:     "",
			LastName:      "",
			UserName:      "",
			Rank:          "",
			SupervisorID:  "",
			Password:      "",
			ExpectedError: backend.ErrMissingArgs,
		},
		{
			Name:          "Weak password",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Rank:          types.E3,
			SupervisorID:  "",
			Password:      "test",
			ExpectedError: backend.ErrWeakPassword,
		},
		{
			Name:          "Password too long",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Rank:          types.E7,
			SupervisorID:  "",
			Password:      "toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong",
			ExpectedError: backend.ErrPasswordTooLong,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			member, err := b.AddMember(types.Member{
				ApiMember: types.ApiMember{
					FirstName:    tt.FirstName,
					LastName:     tt.LastName,
					Username:     tt.UserName,
					Rank:         tt.Rank,
					SupervisorID: tt.SupervisorID,
				},
				Password: tt.Password,
				Hash:     "",
			})
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				m, err := b.GetMember(member.ID)
				if err != nil {
					t.Errorf("Error getting member that should exist: %s", err.Error())
				}
				if !reflect.DeepEqual(m, member) {
					t.Errorf("Expected member: %+v\nGot: %+v", member, m)
				}
			}
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	dbID := uuid.NewString()
	t.Cleanup(func() {
		os.Remove(fmt.Sprintf("%s.db", dbID))
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sqlite.New(logger, fmt.Sprintf("%s.db", dbID), 1.0)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := backend.New(logger, provider, provider, provider, provider, &backend.Options{BcryptCost: bcrypt.MinCost})

	member1 := testutils.RandomMember()
	member2 := testutils.RandomMember()

	type testCase struct {
		Name            string
		ExpectedMembers []types.Member
		SetupFunc       func(*testing.T, *testCase)
		ExpectedError   error
	}

	tc := []testCase{
		{
			Name:            "No members found",
			ExpectedMembers: []types.Member{},
			SetupFunc:       func(_ *testing.T, _ *testCase) {},
			ExpectedError:   nil,
		},
		{
			Name:            "One member",
			ExpectedMembers: []types.Member{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				member1, err = b.AddMember(member1)
				if err != nil {
					t.Fatalf("Error adding member for TestGetAllMembers: %s", err.Error())
				}
				tc.ExpectedMembers = append(tc.ExpectedMembers, member1)
			},
			ExpectedError: nil,
		},
		{
			Name:            "Two members",
			ExpectedMembers: []types.Member{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				member2, err = b.AddMember(member2)
				if err != nil {
					t.Fatalf("Error adding member for TestGetAllMembers: %s", err.Error())
				}
				tc.ExpectedMembers = append(tc.ExpectedMembers, member1, member2)
			},
			ExpectedError: nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SetupFunc(t, &tt)
			members, err := b.GetAllMembers()
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				sort.Slice(members, func(i, j int) bool {
					return members[i].ID < members[j].ID
				})
				sort.Slice(tt.ExpectedMembers, func(i, j int) bool {
					return tt.ExpectedMembers[i].ID < tt.ExpectedMembers[j].ID
				})
				if !slices.Equal(members, tt.ExpectedMembers) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.ExpectedMembers, members)
				}
			}
		})
	}
}

func TestUpdateMember(t *testing.T) {
	dbID := uuid.NewString()
	t.Cleanup(func() {
		os.Remove(fmt.Sprintf("%s.db", dbID))
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sqlite.New(logger, fmt.Sprintf("%s.db", dbID), 1.0)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := backend.New(logger, provider, provider, provider, provider, &backend.Options{BcryptCost: bcrypt.MinCost})

	member, err := b.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestUpdateMember: %s", err.Error())
	}
	supervisor, err := b.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestUpdateMember: %s", err.Error())
	}

	tc := []struct {
		Name          string
		Updates       types.Member
		ExpectedError error
	}{
		{
			Name: "Successful full update",
			Updates: types.Member{
				ApiMember: types.ApiMember{
					ID:           member.ID,
					FirstName:    "Joe",
					LastName:     "Schmoe",
					Username:     "newuser",
					Rank:         types.E1,
					SupervisorID: supervisor.ID,
				},
				Password: "newpassword",
				Hash:     "",
			},
			ExpectedError: nil,
		},
		{
			Name: "New Weak Password",
			Updates: types.Member{
				ApiMember: types.ApiMember{
					ID: member.ID,
				},
				Password: "weak",
			},
			ExpectedError: backend.ErrWeakPassword,
		},
		{
			Name: "New too long password",
			Updates: types.Member{
				ApiMember: types.ApiMember{
					ID: member.ID,
				},
				Password: "toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong",
			},
			ExpectedError: backend.ErrPasswordTooLong,
		},
		{
			Name: "New supervisor doesn't exist",
			Updates: types.Member{
				ApiMember: types.ApiMember{
					ID:           member.ID,
					SupervisorID: uuid.NewString(),
				},
			},
			ExpectedError: backend.ErrSupervisorNotFound,
		},
		{
			Name: "Non-existing member",
			Updates: types.Member{
				ApiMember: types.ApiMember{
					ID: uuid.NewString(),
				},
			},
			ExpectedError: backend.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			mem, err := b.UpdateMember(tt.Updates)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				testutils.VerifyUpdatedUser(member, tt.Updates, mem, t)
			}
		})
	}
}

func TestDeleteMember_Sqlite(t *testing.T) {
	dbID := uuid.NewString()
	t.Cleanup(func() {
		os.Remove(fmt.Sprintf("%s.db", dbID))
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sqlite.New(logger, fmt.Sprintf("%s.db", dbID), 1.0)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := backend.New(logger, provider, provider, provider, provider, &backend.Options{BcryptCost: bcrypt.MinCost})

	m1, err := b.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestDeleteMember_Sqlite: %s", err.Error())
	}
	supervisor, err := b.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestDeleteMember_Sqlite: %s", err.Error())
	}
	m3 := testutils.RandomMember()
	m3.SupervisorID = supervisor.SupervisorID
	m3, err = b.AddMember(m3)
	if err != nil {
		t.Fatalf("Error adding member for TestDeleteMember_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name          string
		Identifier    string
		ExpectedError error
		Verification  func(t *testing.T)
	}{
		{
			Name:          "Successful delete by ID",
			Identifier:    m1.ID,
			ExpectedError: nil,
			Verification: func(t *testing.T) {
				_, err := b.GetMember(m1.ID)
				if !errors.Is(err, backend.ErrMemberNotFound) {
					t.Errorf("Expected to not find user but did")
				}
			},
		},
		{
			Name:          "Removed supervisor",
			Identifier:    supervisor.ID,
			ExpectedError: nil,
			Verification: func(t *testing.T) {
				m, err := b.GetMember(m3.ID)
				if err != nil {
					t.Errorf("Error getting member for TestDeleteMember_Sqlite: %s", err.Error())
					return
				}
				if m.SupervisorID != "" {
					t.Errorf("Expected supervisor ID to be empty after deleting supervisor, but it isn't")
				}
			},
		},
		{
			Name:          "Successful delete by username",
			Identifier:    m3.Username,
			ExpectedError: nil,
			Verification: func(t *testing.T) {
				_, err := b.GetMember(m3.Username)
				if !errors.Is(err, backend.ErrMemberNotFound) {
					t.Errorf("Expected to not find member, but did")
				}
			},
		},
		{
			Name:          "Member not found by ID",
			Identifier:    uuid.NewString(),
			ExpectedError: backend.ErrMemberNotFound,
			Verification:  func(_ *testing.T) {},
		},
		{
			Name:          "Member not found by Username",
			Identifier:    "notfound",
			ExpectedError: backend.ErrMemberNotFound,
			Verification:  func(_ *testing.T) {},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err = b.DeleteMember(tt.Identifier)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			tt.Verification(t)
		})
	}
}
