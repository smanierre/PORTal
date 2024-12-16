package backend_test

import (
	"PORTal/backend"
	"PORTal/providers/sqlite"
	"PORTal/testutils"
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	b := backend.New(logger, provider, provider, provider, provider, backend.Config{BcryptCost: bcrypt.MinCost, SessionTimeout: 4}, nil)

	m, err := b.AddMember(testutils.RandomMember(false))
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
