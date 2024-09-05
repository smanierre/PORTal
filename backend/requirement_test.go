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
	"testing"
)

func TestAddGetRequirement(t *testing.T) {
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

	ref1, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestAddGetRequirement: %s", err.Error())
	}
	req1, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestAddGetRequirement: %s", err.Error())
	}

	// Test duplicate requirement name
	req2 := testutils.RandomRequirement(ref1)
	req2.Name = req1.Name
	_, err = b.AddRequirement(req2)
	if !errors.Is(err, backend.ErrDuplicateRequirement) {
		t.Errorf("Expected error: %s\nGot: %s", backend.ErrDuplicateRequirement.Error(), err)
	}

	// Test missing values
	_, err = b.AddRequirement(types.Requirement{})
	if !errors.Is(err, backend.ErrMissingArgs) {
		t.Errorf("Expected error: %s\nGot: %s", backend.ErrMissingArgs.Error(), err)
	}

	tc := []struct {
		Name                string
		RequirementID       string
		ExpectedRequirement types.Requirement
		ExpectedError       error
	}{
		{
			Name:                "Successful Add",
			RequirementID:       req1.ID,
			ExpectedRequirement: req1,
			ExpectedError:       nil,
		},
		{
			Name:                "Requirement not found",
			RequirementID:       uuid.NewString(),
			ExpectedRequirement: types.Requirement{},
			ExpectedError:       backend.ErrRequirementNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			req, err := b.GetRequirement(tt.RequirementID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil && !reflect.DeepEqual(req, tt.ExpectedRequirement) {
				t.Errorf("Expected: %+v\nGot: %+v", tt.ExpectedRequirement, req)
			}
		})
	}
}

func TestGetAllRequirements(t *testing.T) {
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

	ref1, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestGetAllRequirements: %s", err.Error())
	}
	req1 := testutils.RandomRequirement(ref1)

	ref2, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestGetAllRequirements: %s", err.Error())
	}
	req2 := testutils.RandomRequirement(ref2)

	type testStruct struct {
		Name          string
		Results       []types.Requirement
		ExpectedError error
		SetupFunc     func(t *testing.T, tt *testStruct)
	}

	tc := []testStruct{
		{
			Name:          "No requirements",
			Results:       []types.Requirement{},
			ExpectedError: nil,
			SetupFunc:     func(_ *testing.T, _ *testStruct) {},
		},
		{
			Name:          "One result",
			Results:       []types.Requirement{},
			ExpectedError: nil,
			SetupFunc: func(t *testing.T, tt *testStruct) {
				req1, err = b.AddRequirement(req1)
				if err != nil {
					t.Fatalf("Error adding requirement for TestGetAllRequirements: %s", err.Error())
				}
				tt.Results = append(tt.Results, req1)
			},
		},
		{
			Name:          "Two results",
			Results:       []types.Requirement{},
			ExpectedError: nil,
			SetupFunc: func(t *testing.T, tt *testStruct) {
				req2, err = b.AddRequirement(req2)
				if err != nil {
					t.Fatalf("Error adding requirement for TestGetAllRequirements: %s", err.Error())
				}
				tt.Results = append(tt.Results, req1, req2)
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SetupFunc(t, &tt)
			requirements, err := b.GetAllRequirements()
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				if !slices.EqualFunc(requirements, tt.Results, testutils.CompareRequirements) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.Results, requirements)
				}
			}
		})
	}
}

func TestUpdateRequirement(t *testing.T) {
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

	originalRef, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error when adding reference for TestUpdateRequirement: %s", err.Error())
	}
	originalReq, err := b.AddRequirement(testutils.RandomRequirement(originalRef))
	if err != nil {
		t.Fatalf("Error when adding requirement for TestUpdateRequirement: %s", err.Error())
	}
	newRef, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error when adding reference for TestUpdateRequirement: %s", err.Error())
	}

	uninsertedRef := testutils.RandomReference()
	uninsertedRef.ID = uuid.NewString()
	tc := []struct {
		Name          string
		Update        types.Requirement
		Expected      types.Requirement
		ExpectedError error
	}{
		{
			Name: "Successful full update",
			Update: types.Requirement{
				ID:           originalReq.ID,
				Name:         "New Name",
				Reference:    newRef,
				Description:  "New Description",
				Notes:        "New Notes",
				DaysValidFor: 180,
			},
			Expected: types.Requirement{
				ID:           originalReq.ID,
				Name:         "New Name",
				Reference:    newRef,
				Description:  "New Description",
				Notes:        "New Notes",
				DaysValidFor: 180,
			},
			ExpectedError: nil,
		},
		{
			Name: "Successful single field update",
			Update: types.Requirement{
				ID:           originalReq.ID,
				DaysValidFor: -1,
			},
			Expected: types.Requirement{
				ID:           originalReq.ID,
				Name:         "New Name",
				Reference:    newRef,
				Description:  "New Description",
				Notes:        "New Notes",
				DaysValidFor: -1,
			},
			ExpectedError: nil,
		},
		{
			Name: "Reference not found",
			Update: types.Requirement{
				ID:        originalReq.ID,
				Reference: uninsertedRef,
			},
			Expected:      types.Requirement{},
			ExpectedError: backend.ErrReferenceNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := b.UpdateRequirement(tt.Update)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				got, err := b.GetRequirement(tt.Expected.ID)
				if err != nil {
					t.Errorf("Error getting requirement for TestUpdateRequirement: %s", err.Error())
				}
				if !testutils.CompareRequirements(got, tt.Expected) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.Expected, got)
				}
			}
		})
	}
}

func TestDeleteRequirement(t *testing.T) {
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

	ref1, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestDeleteRequirement: %s", err.Error())
	}
	req1, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteRequirement: %s", err.Error())
	}

	req2, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteRequirement: %s", err.Error())
	}
	qual := testutils.RandomQualification()
	qual.InitialRequirements = []types.Requirement{req2}

	qual, err = b.AddQualification(qual)
	if err != nil {
		t.Fatalf("Error adding qualification for TestDeleteRequirement: %s", err.Error())
	}

	tc := []struct {
		Name          string
		RequirementID string
		ExpectedError error
	}{
		{
			Name:          "Successful delete",
			RequirementID: req1.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Requirement not found",
			RequirementID: uuid.NewString(),
			ExpectedError: backend.ErrRequirementNotFound,
		},
		{
			Name:          "Requirement assigned to qualification",
			RequirementID: req2.ID,
			ExpectedError: backend.ErrRequirementInUse,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err := b.DeleteRequirement(tt.RequirementID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				_, err = b.GetRequirement(tt.RequirementID)
				if !errors.Is(err, backend.ErrRequirementNotFound) {
					t.Errorf("Expected error: %s\nGot: %s", backend.ErrRequirementNotFound.Error(), err.Error())
				}
			}
		})
	}
}
