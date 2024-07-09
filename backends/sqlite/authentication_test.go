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
)

const (
	getSessionForMemberQuery = "SELECT * FROM session WHERE id=(SELECT session_id FROM member_session WHERE member_id=$1);"
)

func TestAddSession(t *testing.T) {
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
			Rank:           types.E1,
			Username:       "username",
			Qualifications: nil,
		},
		Password: "password",
	}
	if err := backend.AddMember(member); err != nil {
		t.Fatalf("Error adding member for TestAddSession: %s", err.Error())
	}

	tc := []struct {
		name          string
		ipAddress     string
		memberID      string
		expectedError error
	}{
		{
			name:          "Successful add",
			ipAddress:     "192.168.0.1",
			memberID:      member.ID,
			expectedError: nil,
		},
		{
			name:          "Member doesn't exist",
			memberID:      uuid.NewString(),
			ipAddress:     "192.168.0.1",
			expectedError: types.ErrMemberNotFound,
		},
		{
			name:          "Invalid IP address",
			memberID:      member.ID,
			ipAddress:     "10.0.0",
			expectedError: types.ErrInvalidIP,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			id, err := backend.AddSession(tt.ipAddress, tt.memberID)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				err = backend.ValidateSession(id, tt.memberID, tt.ipAddress)
				if err != nil {
					t.Errorf("Expected valid session, but got: %s", err.Error())
				}
				row := backend.Db.QueryRow(getSessionForMemberQuery, tt.memberID)
				s := types.Session{}
				err = row.Scan(&s.SessionID, &s.Expires, &s.IPAddress)
				if err != nil {
					t.Errorf("Error when scanning session into struct: %s", err.Error())
				}
			}
		})
	}
}

func TestValidateSession(t *testing.T) {
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
	validIP := "192.168.0.1"
	member := types.Member{
		ApiMember: types.ApiMember{
			ID:        uuid.NewString(),
			FirstName: "Test",
			LastName:  "Member",
			Rank:      types.E9,
			Username:  "username",
		},
		Password: "password",
	}
	if err = backend.AddMember(member); err != nil {
		t.Fatalf("Error adding member for TestValidateSession: %s", err.Error())
	}
	validSessionID, err := backend.AddSession(validIP, member.ID)
	if err != nil {
		t.Fatalf("Error adding session for TestValidateSession: %s", err.Error())
	}

	tc := []struct {
		name          string
		sessionID     string
		ipAddress     string
		memberID      string
		expectedError error
	}{
		{
			name:          "Successful validate",
			sessionID:     validSessionID,
			ipAddress:     validIP,
			memberID:      member.ID,
			expectedError: nil,
		},
		{
			name:          "Invalid member ID",
			sessionID:     validSessionID,
			ipAddress:     validIP,
			memberID:      uuid.NewString(),
			expectedError: types.ErrSessionValidationFailed,
		},
		{
			name:          "Session ID not found",
			sessionID:     uuid.NewString(),
			ipAddress:     validIP,
			memberID:      member.ID,
			expectedError: types.ErrSessionValidationFailed,
		},
		{
			name:          "Invalid IP address",
			sessionID:     validSessionID,
			ipAddress:     "10.0.10.10",
			memberID:      member.ID,
			expectedError: types.ErrSessionValidationFailed,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.ValidateSession(tt.sessionID, tt.memberID, tt.ipAddress)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestLogin(t *testing.T) {
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
			ID:             uuid.NewString(),
			FirstName:      "Test",
			LastName:       "User",
			Username:       "testuser",
			Rank:           types.E8,
			Qualifications: nil,
			SupervisorID:   "",
		},
		Password: "password",
		Hash:     "",
	}
	if err := backend.AddMember(m); err != nil {
		t.Fatalf("Error adding member for TestLogin: %s", err.Error())
	}

	tc := []struct {
		name           string
		username       string
		password       string
		expectedMember types.Member
		expectedError  error
	}{
		{
			name:           "Successful login",
			username:       m.Username,
			password:       m.Password,
			expectedMember: m,
			expectedError:  nil,
		},
		{
			name:           "Bad password",
			username:       m.Username,
			password:       "random",
			expectedMember: types.Member{},
			expectedError:  types.ErrPasswordAuthenticationFailed,
		},
		{
			name:           "Member not found",
			username:       "random",
			password:       m.Password,
			expectedMember: types.Member{},
			expectedError:  types.ErrMemberNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			member, err := backend.Login(tt.username, tt.password)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil && !reflect.DeepEqual(tt.expectedMember.ApiMember, member.ApiMember) {
				t.Errorf("Expected result: %+v\nGot: %+v", tt.expectedMember.ApiMember, member.ApiMember)
			}
		})
	}
}
