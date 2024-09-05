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

func TestAddAndGetQualification(t *testing.T) {
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
		t.Fatalf("Error adding reference for TestAddAndGetQualification: %s", err.Error())
	}
	ref2, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestAddAndGetQualification: %s", err.Error())
	}
	ref3, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestAddAndGetQualification: %s", err.Error())
	}
	ref4, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestAddAndGetQualification: %s", err.Error())
	}
	req1, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestAddQualification_Sqlite: %s", err.Error())
	}
	req2, err := b.AddRequirement(testutils.RandomRequirement(ref2))
	if err != nil {
		t.Fatalf("Error adding requirement for TestAddQualification_Sqlite: %s", err.Error())
	}
	req3, err := b.AddRequirement(testutils.RandomRequirement(ref3))
	if err != nil {
		t.Fatalf("Error adding requirement for TestAddQualification_Sqlite: %s", err.Error())
	}
	req4, err := b.AddRequirement(testutils.RandomRequirement(ref4))
	if err != nil {
		t.Fatalf("Error adding requirement for TestAddQualification_Sqlite: %s", err.Error())
	}

	tc := []struct {
		Name           string
		QualName       string
		Notes          string
		Expires        bool
		ExpirationDays int
		ExpectedError  error
		InitialReqs    []types.Requirement
		RecurringReqs  []types.Requirement
	}{
		{
			Name:           "Successful insert doesn't expire",
			QualName:       testutils.RandomString(),
			Notes:          testutils.RandomString(),
			Expires:        false,
			ExpirationDays: 0,
			ExpectedError:  nil,
		},
		{
			Name:           "Successful insert does expire",
			QualName:       testutils.RandomString(),
			Notes:          testutils.RandomString(),
			Expires:        true,
			ExpirationDays: 100,
			ExpectedError:  nil,
		},
		{
			Name:           "Expires with no expiration days",
			QualName:       testutils.RandomString(),
			Notes:          testutils.RandomString(),
			Expires:        true,
			ExpirationDays: 0,
			ExpectedError:  backend.ErrMissingArgs,
		},
		{
			Name:           "Expires with invalid expiration days",
			QualName:       testutils.RandomString(),
			Notes:          testutils.RandomString(),
			Expires:        true,
			ExpirationDays: -1,
			ExpectedError:  backend.ErrInvalidQualExpiration,
		},
		{
			Name:          "Initial Requirements",
			QualName:      testutils.RandomString(),
			Notes:         testutils.RandomString(),
			ExpectedError: nil,
			InitialReqs:   []types.Requirement{req1, req2},
		},
		{
			Name:           "Recurring Requirements",
			QualName:       testutils.RandomString(),
			Notes:          testutils.RandomString(),
			Expires:        false,
			ExpirationDays: 0,
			ExpectedError:  nil,
			RecurringReqs:  []types.Requirement{req3, req4},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			qual, err := b.AddQualification(types.Qualification{Name: tt.Name, Notes: tt.Notes, Expires: tt.Expires, ExpirationDays: tt.ExpirationDays,
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
		t.Fatalf("Error adding reference for TestGetAllQualifications: %s", err.Error())
	}
	req1, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestGetAllQualifications: %s", err.Error())
	}
	req2, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestGetAllQualifications: %s", err.Error())
	}
	qual1 := testutils.RandomQualification()
	qual1.RecurringRequirements = append(qual1.RecurringRequirements, req1)
	qual1.InitialRequirements = append(qual1.InitialRequirements, req2)

	qual2 := testutils.RandomQualification()

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
				qual, err := b.AddQualification(qual1)
				if err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQualifications: %s", err.Error())
				}
				qual1.ID = qual.ID
				tc.ExpectedQuals = append(tc.ExpectedQuals, qual)
			},
		},
		{
			Name:          "Two qualifications",
			ExpectedQuals: []types.Qualification{},
			ExpectedError: nil,
			Setup: func(t *testing.T, tc *testCase) {
				qual, err := b.AddQualification(qual2)
				if err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQUalifications_Sqlite")
				}
				tc.ExpectedQuals = append(tc.ExpectedQuals, qual1, qual)
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

func TestUpdateQualification(t *testing.T) {
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

	usedRef1, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestUpdateQualification: %s", err.Error())
	}
	usedRef2, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestUpdateQualification: %s", err.Error())
	}
	newRef1, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestUpdateQualification: %s", err.Error())
	}
	newRef2, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestUpdateQualification: %s", err.Error())
	}
	usedReq1, err := b.AddRequirement(testutils.RandomRequirement(usedRef1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestUpdateQualification: %s", err.Error())
	}
	usedReq2, err := b.AddRequirement(testutils.RandomRequirement(usedRef2))
	if err != nil {
		t.Fatalf("Error adding requirement for TestUpdateQualification: %s", err.Error())
	}
	newReq1, err := b.AddRequirement(testutils.RandomRequirement(newRef1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestUpdateQualification: %s", err.Error())
	}
	newReq2, err := b.AddRequirement(testutils.RandomRequirement(newRef2))
	if err != nil {
		t.Fatalf("Error adding requirement for TestUpdateQualification: %s", err.Error())
	}
	original := testutils.RandomQualification()
	original.InitialRequirements = []types.Requirement{usedReq1}
	original.RecurringRequirements = []types.Requirement{usedReq2}
	original, err = b.AddQualification(original)
	if err != nil {
		t.Fatalf("Error adding qualification for TestUpdateQualification_Sqlite: %s", err.Error())
	}

	tc := []struct {
		name                  string
		update                types.Qualification
		forceExpirationUpdate bool
		expectedResult        types.Qualification
		expectedError         error
	}{
		{
			name: "Successful full update",
			update: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name",
				InitialRequirements:   []types.Requirement{newReq1},
				RecurringRequirements: []types.Requirement{newReq2},
				Notes:                 "New notes",
				Expires:               true,
				ExpirationDays:        1000,
			},
			expectedResult: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name",
				InitialRequirements:   []types.Requirement{newReq1},
				RecurringRequirements: []types.Requirement{newReq2},
				Notes:                 "New notes",
				Expires:               true,
				ExpirationDays:        1000,
			},
			expectedError: nil,
		},
		{
			name: "Successful keeping requirements",
			update: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name 2",
				InitialRequirements:   []types.Requirement{newReq1, usedReq1},
				RecurringRequirements: []types.Requirement{newReq2, usedReq2},
				Notes:                 "New notes 2",
				Expires:               false,
				ExpirationDays:        0,
			},
			forceExpirationUpdate: true,
			expectedResult: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name 2",
				InitialRequirements:   []types.Requirement{newReq1, usedReq1},
				RecurringRequirements: []types.Requirement{newReq2, usedReq2},
				Notes:                 "New notes 2",
				Expires:               false,
				ExpirationDays:        0,
			},
			expectedError: nil,
		},
		{
			name: "Successful Update No Requirements",
			update: types.Qualification{
				ID:                    original.ID,
				InitialRequirements:   []types.Requirement{},
				RecurringRequirements: []types.Requirement{},
			},
			expectedResult: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name 2",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "New notes 2",
				Expires:               false,
				ExpirationDays:        0,
			},
			expectedError: nil,
		},
		{
			name: "Single field update",
			update: types.Qualification{
				ID:    original.ID,
				Notes: "New Notes 4",
			},
			expectedResult: types.Qualification{
				ID:                    original.ID,
				Name:                  "New name 2",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "New Notes 4",
				Expires:               false,
				ExpirationDays:        0,
			},
			expectedError: nil,
		},
		{
			name: "Initial Requirement not found",
			update: types.Qualification{
				ID: original.ID,
				InitialRequirements: []types.Requirement{
					types.Requirement{
						ID:           uuid.NewString(),
						Name:         "Random",
						Reference:    newRef1,
						Description:  "Random",
						Notes:        "Random",
						DaysValidFor: 100,
					},
				},
				RecurringRequirements: []types.Requirement{},
			},
			forceExpirationUpdate: false,
			expectedResult:        types.Qualification{},
			expectedError:         backend.ErrRequirementNotFound,
		},
		{
			name: "Recurring Requirement not found",
			update: types.Qualification{
				ID: original.ID,
				RecurringRequirements: []types.Requirement{
					types.Requirement{
						ID:           uuid.NewString(),
						Name:         "Random",
						Reference:    newRef1,
						Description:  "Random",
						Notes:        "Random",
						DaysValidFor: 100,
					},
				},
				InitialRequirements: []types.Requirement{},
			},
			forceExpirationUpdate: false,
			expectedResult:        types.Qualification{},
			expectedError:         backend.ErrRequirementNotFound,
		},
		{
			name: "Initial Requirement not found",
			update: types.Qualification{
				ID: original.ID,
				InitialRequirements: []types.Requirement{
					types.Requirement{
						ID:           uuid.NewString(),
						Name:         "Random",
						Reference:    newRef1,
						Description:  "Random",
						Notes:        "Random",
						DaysValidFor: 100,
					},
				},
				RecurringRequirements: []types.Requirement{},
			},
			forceExpirationUpdate: false,
			expectedResult:        types.Qualification{},
			expectedError:         backend.ErrRequirementNotFound,
		},
		{
			name: "Qualification not found",
			update: types.Qualification{
				ID:                    uuid.NewString(),
				RecurringRequirements: []types.Requirement{},
				InitialRequirements:   []types.Requirement{},
			},
			forceExpirationUpdate: false,
			expectedResult:        types.Qualification{},
			expectedError:         backend.ErrQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			_, err := b.UpdateQualification(tt.update, tt.forceExpirationUpdate)
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
		t.Fatalf("Error adding reference for TestDeleteQualification")
	}
	ref2, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestDeleteQualification")
	}
	req1, err := b.AddRequirement(testutils.RandomRequirement(ref1))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteQualification")
	}
	req2, err := b.AddRequirement(testutils.RandomRequirement(ref2))
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
