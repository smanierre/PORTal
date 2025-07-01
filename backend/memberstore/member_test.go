package memberstore_test

import (
	"PORTal/backend"
	"PORTal/backend/memberstore"
	"PORTal/backend/sessionstore"
	"PORTal/providers/sqlite/memberprovider"
	"PORTal/providers/sqlite/sessionprovider"
	"PORTal/testutils"
	"PORTal/types"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddAndGetMember(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	supervisor, err := b.AddMember(testutils.RandomMember(true))
	if err != nil {
		t.Fatalf("Error adding member for TestAddMember_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name          string
		FirstName     string
		LastName      string
		UserName      string
		Grade         types.Grade
		SupervisorID  string
		Password      string
		ExpectedError error
	}{
		{
			Name:          "Successful Add",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      "username",
			Grade:         types.E4,
			SupervisorID:  "",
			Password:      testutils.RandomString(),
			ExpectedError: nil,
		},
		{
			Name:          "Member with supervisor",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Grade:         types.E8,
			SupervisorID:  supervisor.ID,
			Password:      testutils.RandomString(),
			ExpectedError: nil,
		},
		{
			Name:          "Supervisor doesn't exist",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Grade:         types.E4,
			SupervisorID:  uuid.NewString(),
			Password:      testutils.RandomString(),
			ExpectedError: backend.ErrSupervisorNotFound,
		},
		{
			Name:          "Duplicate username",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      "username",
			Grade:         types.E8,
			SupervisorID:  "",
			Password:      testutils.RandomString(),
			ExpectedError: backend.ErrDuplicateUsername,
		},
		{
			Name:          "Missing fields",
			FirstName:     "",
			LastName:      "",
			UserName:      "",
			Grade:         "",
			SupervisorID:  "",
			Password:      "",
			ExpectedError: backend.ErrMissingArgs,
		},
		{
			Name:          "Weak password",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Grade:         types.E3,
			SupervisorID:  "",
			Password:      "test",
			ExpectedError: backend.ErrWeakPassword,
		},
		{
			Name:          "Password too long",
			FirstName:     testutils.RandomString(),
			LastName:      testutils.RandomString(),
			UserName:      testutils.RandomString(),
			Grade:         types.E7,
			SupervisorID:  "",
			Password:      "toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong",
			ExpectedError: backend.ErrPasswordTooLong,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			member, err := b.AddMember(types.Member{
				FirstName:    tt.FirstName,
				LastName:     tt.LastName,
				Username:     tt.UserName,
				Grade:        tt.Grade,
				SupervisorID: tt.SupervisorID,
				Password:     tt.Password,
				Hash:         "",
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
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	member1 := testutils.RandomMember(true)
	member2 := testutils.RandomMember(false)

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
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	member, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestUpdateMember: %s", err.Error())
	}
	supervisor, err := b.AddMember(testutils.RandomMember(true))
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
				ID:           member.ID,
				FirstName:    "Joe",
				LastName:     "Schmoe",
				Username:     "newuser",
				Grade:        types.E1,
				Admin:        true,
				SupervisorID: supervisor.ID,
				Password:     "newpassword",
				Hash:         "",
			},
			ExpectedError: nil,
		},
		{
			Name: "New Weak Password",
			Updates: types.Member{
				ID:       member.ID,
				Password: "weak",
			},
			ExpectedError: backend.ErrWeakPassword,
		},
		{
			Name: "New too long password",
			Updates: types.Member{
				ID:       member.ID,
				Password: "toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong",
			},
			ExpectedError: backend.ErrPasswordTooLong,
		},
		{
			Name: "New supervisor doesn't exist",
			Updates: types.Member{
				ID:           member.ID,
				SupervisorID: uuid.NewString(),
			},
			ExpectedError: backend.ErrSupervisorNotFound,
		},
		{
			Name: "Non-existing member",
			Updates: types.Member{
				ID: uuid.NewString(),
			},
			ExpectedError: backend.ErrMemberNotFound,
		},
		{
			Name: "No supervisor",
			Updates: types.Member{
				ID:           member.ID,
				FirstName:    "Joe",
				LastName:     "Schmoe",
				Username:     "newuser",
				Grade:        types.E1,
				Admin:        true,
				Hash:         member.Hash,
				SupervisorID: "",
			},
			ExpectedError: nil,
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
			// Because the members password gets updated throughout the tests, the no supervisor one fails on the hash comparison
			// even though it's still valid
			if tt.ExpectedError == nil && tt.Name != "No supervisor" {
				testutils.VerifyUpdatedUser(tt.Updates, mem, t)
			} else if tt.ExpectedError == nil {
				testutils.VerifyUpdatedUserNoHash(tt.Updates, mem, t)
			}
		})
	}
}

func TestDisableMember_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	// Setup for clean disable
	m1, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestDisableMember_Sqlite: %s", err.Error())
	}

	// Setup for repeat disable
	m2, err := b.AddMember(testutils.RandomMember(true))
	if err != nil {
		t.Fatalf("Error adding member for TestDisableMember_Sqlite: %s", err.Error())
	}
	err = b.DisableMember(m2.ID)
	if err != nil {
		t.Fatalf("Error disabling member for TestDisableMember_Sqlite: %s", err.Error())
	}

	// Setup for member with subordinates
	m3, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestDisableMember_Sqlite: %s", err.Error())
	}
	m4 := testutils.RandomMember(false)
	m4.SupervisorID = m3.ID
	m4, err = b.AddMember(m4)
	if err != nil {
		t.Fatalf("Error adding member for TestDisableMember_Sqlite: %s", err.Error())
	}
	m5 := testutils.RandomMember(false)
	m5.SupervisorID = m3.ID
	m5, err = b.AddMember(m5)
	if err != nil {
		t.Fatalf("Error adding member for TestDisableMember_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name             string
		MemberID         string
		ExpectedError    error
		VerificationFunc func(*testing.T)
	}{
		{
			Name:          "Successful disable",
			MemberID:      m1.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Member not found",
			MemberID:      uuid.NewString(),
			ExpectedError: backend.ErrMemberNotFound,
		},
		{
			Name:          "Member already disabled",
			MemberID:      m2.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Member with subordinates",
			MemberID:      m3.ID,
			ExpectedError: nil,
			VerificationFunc: func(t *testing.T) {
				sub1, err := b.GetMember(m4.ID)
				if err != nil {
					t.Errorf("Error getting member for \"Member with subordinates\" validation func: %s", err.Error())
				}
				if sub1.SupervisorID != "" {
					t.Error("Expected sub1's supervisor to be blank, but it wasn't")
				}
				sub2, err := b.GetMember(m5.ID)
				if err != nil {
					t.Errorf("Error getting member for \"Member with subordinates\" validation func: %s", err.Error())
				}
				if sub2.SupervisorID != "" {
					t.Error("Expected sub1's supervisor to be blank, but it wasn't")
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err := b.DisableMember(tt.MemberID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError == nil {
				_, err = b.GetMember(tt.MemberID)
				if !errors.Is(err, backend.ErrMemberDisabled) {
					t.Errorf("Expected member to be disabled, but it wasn't: %s", err)
				}
				if tt.VerificationFunc != nil {
					tt.VerificationFunc(t)
				}
			}
		})
	}
}

func TestEnableMember_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	// Setup successful enable
	m1, err := b.AddMember(testutils.RandomMember(true))
	if err != nil {
		t.Fatalf("Error adding member for TestEnableMember_Sqlite: %s", err.Error())
	}
	err = b.DisableMember(m1.ID)
	if err != nil {
		t.Fatalf("Error disabling member for TestEnableMember_Sqlite: %s", err.Error())
	}
	tc := []struct {
		Name             string
		Id               string
		ExpectedError    error
		VerificationFunc func(*testing.T)
	}{
		{
			Name:          "Successful enable",
			Id:            m1.ID,
			ExpectedError: nil,
			VerificationFunc: func(t *testing.T) {
				m, err := b.GetMember(m1.ID)
				if err != nil {
					t.Errorf("Error getting member for TestEnableMember_Sqlite: %s", err.Error())
				}
				if m.Disabled {
					t.Errorf("Expected member to be enabled, but it isn't")
				}
			},
		},
		{
			Name:             "Member not found",
			Id:               uuid.NewString(),
			ExpectedError:    backend.ErrMemberNotFound,
			VerificationFunc: nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err = b.EnableMember(tt.Id)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && err == nil {
				t.Errorf("Expected error: %s", tt.ExpectedError.Error())
			}
			if tt.VerificationFunc != nil {
				tt.VerificationFunc(t)
			}
		})
	}
}

func TestGetDisabledMembers_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	// Create members to disable later
	m1, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when creating member for TestGetDisabledMembers_Sqlite: %s", err.Error())
	}
	m2, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when creating member for TestGetDisabledMembers_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name             string
		ExpectedError    error
		ExpectedResponse []types.Member
		SetupFunc        func(t *testing.T)
	}{
		{
			Name:             "No members",
			ExpectedResponse: []types.Member{},
			ExpectedError:    nil,
			SetupFunc:        func(t *testing.T) {},
		},
		{
			Name:             "One member",
			ExpectedError:    nil,
			ExpectedResponse: []types.Member{m1},
			SetupFunc: func(t *testing.T) {
				err := b.DisableMember(m1.ID)
				if err != nil {
					t.Errorf("Error disabling member for TestGetDisabledMembers_Sqlite: %s", err.Error())
				}
			},
		},
		{
			Name:             "Two members",
			ExpectedError:    nil,
			ExpectedResponse: []types.Member{m1, m2},
			SetupFunc: func(t *testing.T) {
				err := b.DisableMember(m2.ID)
				if err != nil {
					t.Errorf("Error disabling member for TestGetDisabledMembers_Sqlite: %s", err.Error())
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SetupFunc(t)
			members, err := b.GetDisabledMembers()
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if !slices.EqualFunc(members, tt.ExpectedResponse, func(m1, m2 types.Member) bool {
				return m1 == m2
			}) {
				t.Errorf("Expected response: %v\nGot: %v", tt.ExpectedResponse, members)
			}
		})
	}
}

func TestLogin_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	member := testutils.RandomMember(false)
	password := member.Password
	member, err = b.AddMember(member)
	if err != nil {
		t.Fatalf("Error when creating member for TestLogin_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name          string
		Username      string
		Password      string
		ExpectedError error
	}{
		{
			Name:          "Successful Login",
			Username:      member.Username,
			Password:      password,
			ExpectedError: nil,
		},
		{
			Name:          "Wrong password",
			Username:      member.Username,
			Password:      "hopefullywrong",
			ExpectedError: backend.ErrAuthenticationFailed,
		},
		{
			Name:          "User not found",
			Username:      "notfound",
			Password:      password,
			ExpectedError: backend.ErrAuthenticationFailed,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			mem, err := b.Login(tt.Username, tt.Password)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && err == nil {
				t.Errorf("Expected error: %s", tt.ExpectedError.Error())
			}
			if err == nil {
				testutils.VerifyUpdatedUser(member, mem, t)
			}
		})
	}
}

func TestGetPotentialSupervisors_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	m1 := testutils.RandomMember(false)
	m1.Grade = types.E1
	m1, err = b.AddMember(m1)
	if err != nil {
		t.Errorf("Error when adding member for TestGetPotentialSupervisors_Sqlite: %s", err.Error())
	}

	m2 := testutils.RandomMember(false)
	m2.Grade = types.E3
	m2, err = b.AddMember(m2)
	if err != nil {
		t.Errorf("Error when adding member for TestGetPotentialSupervisors_Sqlite: %s", err.Error())
	}

	m3 := testutils.RandomMember(false)
	m3.Grade = types.E3
	m3, err = b.AddMember(m3)
	if err != nil {
		t.Errorf("Error when adding member for TestGetPotentialSupervisors_Sqlite: %s", err.Error())
	}

	m4 := testutils.RandomMember(false)
	m4.Grade = types.E8
	m4, err = b.AddMember(m4)
	if err != nil {
		t.Errorf("Error when adding member for TestGetPotentialSupervisors_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name                string
		Member              types.Member
		ExpectedSupervisors []types.Member
	}{
		{
			Name:                "Three Supervisors",
			Member:              m1,
			ExpectedSupervisors: []types.Member{m2, m3, m4},
		},
		{
			Name:                "Two supervisors",
			Member:              m2,
			ExpectedSupervisors: []types.Member{m3, m4},
		},
		{
			Name:                "No Supervisors",
			Member:              m4,
			ExpectedSupervisors: []types.Member{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			supervisors, err := b.GetPotentialSupervisors(tt.Member, tt.Member.Grade)
			if err != nil {
				t.Errorf("Error when getting potential supervisors: %s", err.Error())
			}
			if !slices.Equal(supervisors, tt.ExpectedSupervisors) {
				t.Errorf("Expected supervisors: %v\nGot: %v", tt.ExpectedSupervisors, supervisors)
			}
		})
	}
}

func TestGetDisabledMember_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	m, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when adding member for TestGetDisabledMember_Sqlite: %s", err.Error())
	}
	err = b.DisableMember(m.ID)
	if err != nil {
		t.Fatalf("Error disabling member for TestGetDisabledMember_Sqlite: %s", err.Error())
	}

	m2, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when adding member for TestGetDisabledMember_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name          string
		MemberID      string
		ExpectedError error
	}{
		{
			Name:          "Successful Get",
			MemberID:      m.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Non-Existent Member",
			MemberID:      uuid.NewString(),
			ExpectedError: backend.ErrMemberNotFound,
		},
		{
			Name:          "Member Enabled",
			MemberID:      m2.ID,
			ExpectedError: backend.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			_, err = b.GetDisabledMember(tt.MemberID)
			if !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %v\nGot: %v", tt.ExpectedError, err)
			}
		})
	}
}

func TestGetMemberFromSession_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	m, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when adding member for TestGetMemberFromSession_Sqlite: %s", err.Error())
	}
	sessionProvider, err := sessionprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating session provider for TestGetMemberFromSession_Sqlite: %s", err.Error())
	}
	sessionStore := sessionstore.New(sessionProvider, b, nil, 24*time.Hour, logger)
	sessionID, _ := sessionStore.CreateSession(m.ID, "testing", "127.0.0.1")

	tc := []struct {
		Name          string
		Member        types.Member
		SessionID     string
		ExpectedError error
	}{
		{
			Name:          "Successful Get",
			Member:        m,
			SessionID:     sessionID,
			ExpectedError: nil,
		},
		{
			Name:          "SessionID not found",
			Member:        types.Member{},
			SessionID:     uuid.NewString(),
			ExpectedError: backend.ErrSessionNotFound,
		},
	}
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			mem, err := b.GetMemberFromSession(tt.SessionID)
			if !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError, err)
			}
			if tt.ExpectedError == nil && mem != tt.Member {
				t.Errorf("Expected member: %v\nGot: %v", tt.Member, mem)
			}
		})
	}
}

func TestGetSubordinates_Sqlite(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := memberstore.New(provider, 4, logger)

	noSubs, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when adding member for TestGetSubordinates_Sqlite: %s", err.Error())
	}
	supervisor1, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when adding member for TestGetSubordinates_Sqlite: %s", err.Error())
	}

	supervisor2, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error when adding member for TestGetSubordinates_Sqlite: %s", err.Error())
	}

	sub1 := testutils.RandomMember(false)
	sub1.SupervisorID = supervisor1.ID
	sub1, err = b.AddMember(sub1)
	if err != nil {
		t.Fatalf("Error when adding member for TestGetSubordinates_Sqlite: %s", err.Error())
	}
	sub2 := testutils.RandomMember(false)
	sub2.SupervisorID = supervisor2.ID
	sub2, err = b.AddMember(sub2)
	if err != nil {
		t.Fatalf("Error when adding member for TestGetSubordinates_Sqlite: %s", err.Error())
	}
	sub3 := testutils.RandomMember(false)
	sub3.SupervisorID = supervisor2.ID
	sub3, err = b.AddMember(sub3)
	if err != nil {
		t.Fatalf("Error when adding member for TestGetSubordinates_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name                 string
		MemberID             string
		ExpectedSubordinates []types.Member
	}{
		{
			Name:                 "No Subordinates",
			MemberID:             noSubs.ID,
			ExpectedSubordinates: []types.Member{},
		},
		{
			Name:                 "One Subordinate",
			MemberID:             supervisor1.ID,
			ExpectedSubordinates: []types.Member{sub1},
		},
		{
			Name:                 "Two Subordinates",
			MemberID:             supervisor2.ID,
			ExpectedSubordinates: []types.Member{sub2, sub3},
		},
		{
			Name:                 "Supervisor Not Found",
			MemberID:             uuid.NewString(),
			ExpectedSubordinates: []types.Member{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			subs := b.GetSubordinates(tt.MemberID)
			if !slices.Equal(tt.ExpectedSubordinates, subs) {
				t.Errorf("Expected subordinates: %v\nGot: %v", tt.ExpectedSubordinates, subs)
			}
		})
	}
}
