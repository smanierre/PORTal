package qualificationstore_test

import (
	"PORTal/backend"
	"PORTal/backend/qualificationstore"
	"PORTal/providers/sqlite/qualificationprovider"
	"PORTal/testutils"
	"PORTal/types"
	"bytes"
	"errors"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"slices"
	"sort"
	"testing"
)

// There are no requirements assigned because a qualification must be created before assigning qualifications to it.
// This is a limitation of the web UI, not because it can't be done.
func TestAddAndGetQualification(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := qualificationstore.New(provider, logger)

	tc := []struct {
		Name          string
		QualName      string
		Notes         string
		ExpectedError error
		InitialReqs   []types.Requirement
		RecurringReqs []types.Requirement
	}{
		{
			Name:          "No Requirements",
			QualName:      testutils.RandomString(),
			Notes:         testutils.RandomString(),
			ExpectedError: nil,
			InitialReqs:   []types.Requirement{},
			RecurringReqs: []types.Requirement{},
		},
		{
			Name:          "Missing Args",
			QualName:      "",
			Notes:         "",
			ExpectedError: backend.ErrMissingArgs,
			InitialReqs:   []types.Requirement{},
			RecurringReqs: []types.Requirement{},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			qual, err := b.AddQualification(types.Qualification{Name: tt.QualName, Notes: tt.Notes,
				InitialRequirements: tt.InitialReqs, RecurringRequirements: tt.RecurringReqs})
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				got, err := b.GetQualification(qual.ID)
				if err != nil {
					t.Errorf("Expected no error when getting inserted qual, but got: %s", err.Error())
				}
				if !testutils.CompareQuals(got, qual) {
					t.Errorf("Expected Qualification: %+v\nGot: %+v", qual, got)
				}
			}
		})
	}
}

func TestGetAllQualifications(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := qualificationstore.New(provider, logger)

	initialRequirement1, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding Initial requirement for TestGetAllQualifications: %s", err.Error())
	}
	initialRequirement2, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding Initial requirement for TestGetAllQualifications: %s", err.Error())
	}
	recurringRequirement1, err := b.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding recurring requirement for TestGetAllQualifications: %s", err.Error())
	}
	recurringRequirement2, err := b.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding recurring requirement for TestGetAllQualifications: %s", err.Error())
	}

	var expectedQuals []types.Qualification
	type testCase struct {
		Name          string
		ExpectedQuals []types.Qualification
		ExpectedError error
		Setup         func(t *testing.T, tc *testCase)
	}
	tc := []testCase{
		{
			Name:          "No Qualifications",
			ExpectedQuals: []types.Qualification{},
			ExpectedError: nil,
			Setup:         func(_ *testing.T, _ *testCase) {},
		},
		{
			Name:          "One Qualification",
			ExpectedQuals: []types.Qualification{},
			ExpectedError: nil,
			Setup: func(t *testing.T, tc *testCase) {
				qual1, err := b.AddQualification(testutils.RandomQualification())
				if err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQualifications: %s", err.Error())
				}
				tc.ExpectedQuals = append(expectedQuals, qual1)
				expectedQuals = append(expectedQuals, qual1)
			},
		},
		{
			Name:          "Two qualifications",
			ExpectedQuals: []types.Qualification{},
			ExpectedError: nil,
			Setup: func(t *testing.T, tc *testCase) {
				qual2, err := b.AddQualification(testutils.RandomQualification())
				if err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQualifications: %s", err.Error())
				}
				err = b.AssignRequirementToQualification(qual2.ID, initialRequirement1.ID, true)
				if err != nil {
					t.Fatalf("Error assigning Initial requirement to qualification for TestGetAllQualifications: %s", err.Error())
				}
				err = b.AssignRequirementToQualification(qual2.ID, recurringRequirement1.ID, false)
				if err != nil {
					t.Fatalf("Error assigning recurring requirement to qualification for TestGetAllQualifications: %s", err.Error())
				}
				qual2.InitialRequirements = append(qual2.InitialRequirements, initialRequirement1)
				qual2.RecurringRequirements = append(qual2.RecurringRequirements, recurringRequirement1)
				tc.ExpectedQuals = append(expectedQuals, qual2)
				expectedQuals = append(expectedQuals, qual2)
			},
		},
		{
			Name:          "Three qualifications",
			ExpectedQuals: []types.Qualification{},
			ExpectedError: nil,
			Setup: func(t *testing.T, tc *testCase) {
				qual3, err := b.AddQualification(testutils.RandomQualification())
				if err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQualifications: %s", err.Error())
				}
				err = b.AssignRequirementToQualification(qual3.ID, initialRequirement1.ID, true)
				if err != nil {
					t.Fatalf("Error assigning Initial requirement to qualification for TestGetAllQualifications: %s", err.Error())
				}
				err = b.AssignRequirementToQualification(qual3.ID, initialRequirement2.ID, true)
				if err != nil {
					t.Fatalf("Error assigning Initial requirement to qualification for TestGetAllQualifications: %s", err.Error())
				}
				qual3.InitialRequirements = append(qual3.InitialRequirements, initialRequirement1, initialRequirement2)
				err = b.AssignRequirementToQualification(qual3.ID, recurringRequirement1.ID, false)
				if err != nil {
					t.Fatalf("Error assigning recurring requirement to qualification for TestGetAllQualifications: %s", err.Error())
				}
				err = b.AssignRequirementToQualification(qual3.ID, recurringRequirement2.ID, false)
				if err != nil {
					t.Fatalf("Error assigning recurring requirement to qualification for TestGetAllQualifications: %s", err.Error())
				}
				qual3.RecurringRequirements = append(qual3.RecurringRequirements, recurringRequirement1, recurringRequirement2)
				tc.ExpectedQuals = append(expectedQuals, qual3)
				expectedQuals = append(expectedQuals, qual3)
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Setup(t, &tt)
			quals, err := b.GetAllQualifications()
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				// Sort slices to be the same then compare quals
				sort.Slice(quals, func(i, j int) bool {
					return quals[i].ID > quals[j].ID
				})
				sort.Slice(tt.ExpectedQuals, func(i, j int) bool {
					return tt.ExpectedQuals[i].ID > tt.ExpectedQuals[j].ID
				})
				if !slices.EqualFunc(quals, tt.ExpectedQuals, testutils.CompareQuals) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.ExpectedQuals, quals)
				}
			}
		})
	}
}

// Updating a qualification only handles things like Name, description, etc. Requirements are managed separately.
func TestUpdateQualification(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := qualificationstore.New(provider, logger)

	original := testutils.RandomQualification()
	original.InitialRequirements = []types.Requirement{}
	original.RecurringRequirements = []types.Requirement{}
	original, err = b.AddQualification(original)
	if err != nil {
		t.Fatalf("Error adding qualification for TestUpdateQualification_Sqlite: %s", err.Error())
	}

	tc := []struct {
		name           string
		update         types.Qualification
		expectedResult types.Qualification
		expectedError  error
	}{
		{
			name: "Successful full update",
			update: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name",
				InitialRequirements:   []types.Requirement{},
				RecurringRequirements: []types.Requirement{},
				Notes:                 "New notes",
			},
			expectedResult: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name",
				InitialRequirements:   []types.Requirement{},
				RecurringRequirements: []types.Requirement{},
				Notes:                 "New notes",
			},
			expectedError: nil,
		},
		{
			name: "Single field update",
			update: types.Qualification{
				ID:    original.ID,
				Name:  "New name",
				Notes: "New Notes 2",
			},
			expectedResult: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name",
				InitialRequirements:   []types.Requirement{},
				RecurringRequirements: []types.Requirement{},
				Notes:                 "New Notes 2",
			},
			expectedError: nil,
		},
		{
			name: "Qualification not found",
			update: types.Qualification{
				ID:                    uuid.NewString(),
				Name:                  "doesn't matter",
				Notes:                 "doesn't matter",
				RecurringRequirements: []types.Requirement{},
				InitialRequirements:   []types.Requirement{},
			},
			expectedResult: types.Qualification{},
			expectedError:  backend.ErrQualificationNotFound,
		},
		{
			name: "Missing fields",
			update: types.Qualification{
				ID: original.ID,
			},
			expectedResult: types.Qualification{},
			expectedError:  backend.ErrMissingArgs,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			_, err := b.UpdateQualification(tt.update)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.expectedError.Error(), err.Error())
			}

			if tt.expectedError == nil {
				updatedQual, err := b.GetQualification(tt.update.ID)
				if err != nil {
					t.Errorf("Error getting qualification from database: %s", err.Error())
				}
				if !testutils.CompareQuals(updatedQual, tt.expectedResult) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.expectedResult, updatedQual)
				}
			}
		})
	}
}

func TestDeleteQualification(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	provider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating provider for tests: %s", err.Error())
	}
	b := qualificationstore.New(provider, logger)

	req1, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteQualification")
	}
	req2, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteQualification")
	}
	qual := testutils.RandomQualification()
	qual.InitialRequirements = []types.Requirement{req1}
	qual.RecurringRequirements = []types.Requirement{req2}
	qual, err = b.AddQualification(qual)
	if err != nil {
		t.Fatalf("Error adding qualification for TestDeleteQualification")
	}

	tc := []struct {
		Name          string
		ID            string
		ExpectedError error
	}{
		{
			Name:          "Successful delete",
			ID:            qual.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Qualification Not Found",
			ID:            uuid.NewString(),
			ExpectedError: backend.ErrQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err := b.DeleteQualification(tt.ID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				_, err = b.GetQualification(tt.ID)
				if !errors.Is(err, backend.ErrQualificationNotFound) {
					t.Errorf("Expected qualification to not be found, but got: %s", err.Error())
				}
			}
		})
	}
}
