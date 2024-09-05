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
	"testing"
	"time"
)

type expireClock struct{}

func (e expireClock) Now() time.Time {
	t, err := time.Parse(time.DateTime, "2000-01-01 01:00:00")
	if err != nil {
		panic(fmt.Sprintf("Error parsing time: %s", err.Error()))
	}
	return t
}
func TestAddValidateSession(t *testing.T) {
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
		t.Fatalf("Error adding member for TestAddSession: %s", err.Error())
	}
	member1Session, err := b.AddSession(member.ID, testutils.RandomString())
	member2, err := b.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestAddSession: %s", err.Error())
	}

	expiredBackend := backend.New(logger, provider, provider, provider, provider, &backend.Options{
		BcryptCost: bcrypt.MinCost,
		Clock:      expireClock{},
	})

	member3, err := expiredBackend.AddMember(testutils.RandomMember())
	if err != nil {
		t.Fatalf("Error adding member for TestAddValidateSession: %s", err.Error())
	}
	member3Session, err := expiredBackend.AddSession(member3.ID, testutils.RandomString())
	if err != nil {
		t.Fatalf("Error adding session for TestAddValidateSession: %s", err.Error())
	}

	// Test member doesn't exist
	_, err = b.AddSession(uuid.NewString(), testutils.RandomString())
	if !errors.Is(err, backend.ErrMemberNotFound) {
		t.Errorf("Expected error: %s\nGot: %s", backend.ErrMemberNotFound.Error(), err)
	}

	// Test missing args
	_, err = b.AddSession("", "")
	if !errors.Is(err, backend.ErrMissingArgs) {
		t.Errorf("Expected error: %s\nGot: %s", backend.ErrMissingArgs, err)
	}

	tc := []struct {
		Name            string
		MemberID        string
		UserAgent       string
		SessionID       string
		ExpectedError   error
		BackendOverride *backend.Backend
	}{
		{
			Name:          "Successful add",
			MemberID:      member.ID,
			UserAgent:     member1Session.UserAgent,
			SessionID:     member1Session.SessionID,
			ExpectedError: nil,
		},
		{
			Name:          "No session",
			MemberID:      member2.ID,
			UserAgent:     testutils.RandomString(),
			SessionID:     uuid.NewString(),
			ExpectedError: backend.ErrSessionValidationFailed,
		},
		{
			Name:          "Invalid user agent",
			MemberID:      member.ID,
			UserAgent:     testutils.RandomString(),
			SessionID:     member1Session.SessionID,
			ExpectedError: backend.ErrSessionValidationFailed,
		},
		{
			Name:            "Expired session",
			MemberID:        member3.ID,
			UserAgent:       member3Session.UserAgent,
			SessionID:       member3Session.SessionID,
			ExpectedError:   backend.ErrSessionValidationFailed,
			BackendOverride: &expiredBackend,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err := b.ValidateSession(tt.SessionID, tt.MemberID, tt.UserAgent)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
		})
	}
}

func TestLogin(t *testing.T) {
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

	member := testutils.RandomMember()
	password := member.Password // Need to remember password as it gets removed after hashing
	member, err = b.AddMember(member)
	if err != nil {
		t.Fatalf("Error adding member for TestLogin: %s", err.Error())
	}

	tc := []struct {
		Name          string
		Member        types.Member
		Password      string
		ExpectedError error
	}{
		{
			Name:          "Successful login",
			Member:        member,
			Password:      password,
			ExpectedError: nil,
		},
		{
			Name: "Invalid username",
			Member: types.Member{
				ApiMember: types.ApiMember{
					Username: testutils.RandomString(),
				},
			},
			Password:      password,
			ExpectedError: backend.ErrAuthenticationFailed,
		},
		{
			Name: "Invalid password",
			Member: types.Member{
				ApiMember: types.ApiMember{
					Username: member.Username,
				},
			},
			Password:      testutils.RandomString(),
			ExpectedError: backend.ErrAuthenticationFailed,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			member, err := b.Login(tt.Member.Username, tt.Password)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil && !reflect.DeepEqual(member, tt.Member) {
				t.Errorf("Expected member: %+v\nGot: %+v", tt.Member, member)
			}
		})
	}
}
