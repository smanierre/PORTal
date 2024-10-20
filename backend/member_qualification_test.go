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
	"slices"
	"sort"
	"testing"
)

func TestAddGetMemberQualification(t *testing.T) {
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
	b := backend.New(logger, provider, provider, provider, backend.Config{BcryptCost: bcrypt.MinCost}, nil)

	member1, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestAddGetMemberQualification: %s", err.Error())
	}
	qualification1, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for TestAddGetMemberQualification: %s", err.Error())
	}
	qualification2, err := b.AddQualification(testutils.RandomQualification())
	err = b.AssignMemberQualification(member1.ID, qualification1.ID)
	if err != nil {
		t.Fatalf("Error adding qualification for TestAddGetMemberQualification: %s", err.Error())
	}

	// Test assigning qualification again
	err = b.AssignMemberQualification(member1.ID, qualification1.ID)
	if !errors.Is(err, backend.ErrQualificationAlreadyAssigned) {
		t.Errorf("Expected error %s, got: %s", backend.ErrQualificationAlreadyAssigned.Error(), err.Error())
	}

	// Test add member not found
	err = b.AssignMemberQualification(uuid.NewString(), qualification1.ID)
	if !errors.Is(err, backend.ErrMemberNotFound) {
		t.Errorf("Expected error %s, got: %s", backend.ErrMemberNotFound.Error(), err.Error())
	}

	// Test add qualification not found
	err = b.AssignMemberQualification(member1.ID, uuid.NewString())
	if !errors.Is(err, backend.ErrQualificationNotFound) {
		t.Errorf("Expected error %s, got: %s", backend.ErrQualificationNotFound.Error(), err.Error())
	}

	nonExistentQual := testutils.RandomQualification()
	nonExistentQual.ID = uuid.NewString()
	tc := []struct {
		Name          string
		MemberID      string
		Qualification types.Qualification
		ExpectedError error
	}{
		{
			Name:          "Successful add",
			MemberID:      member1.ID,
			Qualification: qualification1,
			ExpectedError: nil,
		},
		{
			Name:          "Member does not exist",
			MemberID:      uuid.NewString(),
			Qualification: qualification1,
			ExpectedError: backend.ErrMemberQualificationNotFound,
		},
		{
			Name:          "Qualification not assigned",
			MemberID:      member1.ID,
			Qualification: qualification2,
			ExpectedError: backend.ErrMemberQualificationNotFound,
		},
		{
			Name:          "Qualification doesn't exist",
			MemberID:      member1.ID,
			Qualification: nonExistentQual,
			ExpectedError: backend.ErrMemberQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			qual, err := b.GetMemberQualification(tt.MemberID, tt.Qualification.ID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				if !testutils.CompareQuals(qual, tt.Qualification) {
					t.Errorf("Expected qualification: %+v\nGot: %+v", tt.Qualification, qual)
				}
			}
		})
	}
}

func TestGetMemberQualifications(t *testing.T) {
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
	b := backend.New(logger, provider, provider, provider, backend.Config{BcryptCost: bcrypt.MinCost}, nil)

	member1, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestGetMemberQualifications: %s", err.Error())
	}
	member2, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestGetMemberQualifications: %s", err.Error())
	}
	qualification1, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for TestGetMemberQualifications: %s", err.Error())
	}
	err = b.AssignMemberQualification(member2.ID, qualification1.ID)
	if err != nil {
		t.Fatalf("Error adding qualification to member for TestGetMemberQualifications: %s", err.Error())
	}
	member3, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestGetMemberQualifications: %s", err.Error())
	}
	qualification2, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for TestGetMemberQualifications: %s", err.Error())
	}
	err = b.AssignMemberQualification(member3.ID, qualification2.ID)
	if err != nil {
		t.Fatalf("Error adding qualification to member for TestGetMemberQualifications: %s", err.Error())
	}
	qualification3, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for TestGetMemberQualifications: %s", err.Error())
	}
	err = b.AssignMemberQualification(member3.ID, qualification3.ID)
	if err != nil {
		t.Fatalf("Error adding qualification to member for TestGetMemberQualifications: %s", err.Error())
	}

	tc := []struct {
		Name                   string
		MemberID               string
		ExpectedQualifications []types.Qualification
		ExpectedError          error
	}{
		{
			Name:                   "No qualifications",
			MemberID:               member1.ID,
			ExpectedQualifications: []types.Qualification{},
			ExpectedError:          nil,
		},
		{
			Name:                   "One qualification",
			MemberID:               member2.ID,
			ExpectedQualifications: []types.Qualification{qualification1},
			ExpectedError:          nil,
		},
		{
			Name:                   "Two qualifications",
			MemberID:               member3.ID,
			ExpectedQualifications: []types.Qualification{qualification2, qualification3},
			ExpectedError:          nil,
		},
		{
			Name:                   "Member not found",
			MemberID:               uuid.NewString(),
			ExpectedQualifications: []types.Qualification{},
			ExpectedError:          nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			qualifications, err := b.GetMemberQualifications(tt.MemberID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				sort.Slice(tt.ExpectedQualifications, func(i, j int) bool {
					return tt.ExpectedQualifications[i].ID > tt.ExpectedQualifications[j].ID
				})
				sort.Slice(qualifications, func(i, j int) bool {
					return qualifications[i].ID > qualifications[j].ID
				})
				if !slices.EqualFunc(qualifications, tt.ExpectedQualifications, testutils.CompareQuals) {
					t.Errorf("Expected qualifications: %+v\nGot: %+v", tt.ExpectedQualifications, qualifications)
				}
			}
		})
	}
}

func TestRemoveMemberQualification(t *testing.T) {
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
	b := backend.New(logger, provider, provider, provider, backend.Config{BcryptCost: bcrypt.MinCost}, nil)

	member1, err := b.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for TestRemoveMemberQualification: %s", err.Error())
	}
	qual1, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding Qualification for TestRemoveMemberQualification: %s", err.Error())
	}
	err = b.AssignMemberQualification(member1.ID, qual1.ID)
	if err != nil {
		t.Fatalf("Error assigning qualification to member for TestRemoveMemberQualification: %s", err.Error())
	}

	tc := []struct {
		Name            string
		MemberID        string
		QualificationID string
		ExpectedError   error
	}{
		{
			Name:            "Successful delete",
			MemberID:        member1.ID,
			QualificationID: qual1.ID,
			ExpectedError:   nil,
		},
		{
			Name:            "Member not found",
			MemberID:        uuid.NewString(),
			QualificationID: qual1.ID,
			ExpectedError:   backend.ErrMemberQualificationNotFound,
		},
		{
			Name:            "Qualification not found",
			MemberID:        member1.ID,
			QualificationID: uuid.NewString(),
			ExpectedError:   backend.ErrMemberQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err := b.RemoveMemberQualification(tt.MemberID, tt.QualificationID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				_, err := b.GetMemberQualification(tt.MemberID, tt.QualificationID)
				if !errors.Is(err, backend.ErrMemberQualificationNotFound) {
					t.Errorf("Expected member qualification to not be found but got error: %s", err)
				}
			}
		})
	}
}
