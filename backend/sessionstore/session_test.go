package sessionstore_test

import (
	"PORTal/backend"
	"PORTal/backend/memberstore"
	"PORTal/backend/sessionstore"
	"PORTal/providers/sqlite/memberprovider"
	"PORTal/providers/sqlite/sessionprovider"
	"PORTal/testutils"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
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

func TestCreateAndValidateSession(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sessionprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for TestCreateAndValidateSession: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)
	b := sessionstore.New(provider, memberStore, nil, 1*time.Hour, logger)

	m, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}

	tc := []struct {
		Name            string
		ExpectedError   error
		MemberId        string
		UserAgent       string
		VerifyUserAgent string
		IpAddress       string
		VerifyIpAddress string
	}{
		{
			Name:            "Successful Validate",
			ExpectedError:   nil,
			MemberId:        m.ID,
			UserAgent:       "test agent",
			VerifyUserAgent: "test agent",
			IpAddress:       "localhost",
			VerifyIpAddress: "localhost",
		},
		{
			Name:            "Mismatched IP",
			ExpectedError:   backend.ErrSessionValidationFailed,
			MemberId:        m.ID,
			UserAgent:       "test agent",
			VerifyUserAgent: "test agent",
			IpAddress:       "localhost",
			VerifyIpAddress: "127.0.0.1",
		},
		{
			Name:            "Mismatched userAgent",
			ExpectedError:   backend.ErrSessionValidationFailed,
			MemberId:        m.ID,
			UserAgent:       "testAgent",
			VerifyUserAgent: "different agent",
			IpAddress:       "localhost",
			VerifyIpAddress: "localhost",
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			s, _ := b.CreateSession(tt.MemberId, tt.UserAgent, tt.IpAddress)
			member, err := b.ValidateSession(s, tt.VerifyUserAgent, tt.VerifyIpAddress)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError, err)
			}
			if tt.ExpectedError == nil {
				if member.ID != tt.MemberId {
					t.Errorf("Expected member with ID: %s\nGot: %s", tt.MemberId, member.ID)
				}
			}
		})
	}
}

func TestExpiredSession(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sessionprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for TestCreateAndValidateSession: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)

	// Create session store that will set session expiration way in the past
	b := sessionstore.New(provider, memberStore, expireClock{}, 1*time.Hour, logger)

	m, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestExpiredSession: %s", err.Error())
	}

	sessionID, _ := b.CreateSession(m.ID, "test", "127.0.0.1")
	// Create session store using current time to force expiration
	b = sessionstore.New(provider, memberStore, nil, 1*time.Hour, logger)
	_, err = b.ValidateSession(sessionID, "test", "127.0.0.1")

	if !errors.Is(err, backend.ErrSessionValidationFailed) {
		t.Errorf("Expected expected error: %s\n got: %s", backend.ErrSessionValidationFailed.Error(), err)
	}
}

func TestDeleteSession(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := sessionprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for TestCreateAndValidateSession: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)
	b := sessionstore.New(provider, memberStore, nil, 1*time.Hour, logger)

	m, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestExpiredSession: %s", err.Error())
	}
	sessionID, _ := b.CreateSession(m.ID, "test", "127.0.0.1")
	sessionID2, _ := b.CreateSession(m.ID, "test", "127.0.0.1")

	tc := []struct {
		Name          string
		SessionID     string
		ExpectedError error
		SetupFunc     func(t *testing.T)
	}{
		{
			Name:      "Successful Delete",
			SessionID: sessionID,
			SetupFunc: nil,
		},
		{
			Name:      "Member deleted",
			SessionID: sessionID2,
			SetupFunc: func(t *testing.T) {
				err := memberStore.DeleteMember(m.ID)
				if err != nil {
					t.Fatalf("Error deleting member for TestDeleteSession: %s", err.Error())
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SetupFunc != nil {
				tt.SetupFunc(t)
			}
			// There may be a better way to handle this, but that can be worried about when more tests are added
			if tt.Name != "Member Deleted" {
				b.DeleteSession(tt.SessionID)
			}
			_, err = provider.GetSession(tt.SessionID)
			if !errors.Is(err, backend.ErrSessionNotFound) {
				t.Errorf("Expected error %s\nGot: %s", backend.ErrSessionNotFound, err)
			}
		})
	}
}
