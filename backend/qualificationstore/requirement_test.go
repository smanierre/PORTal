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
	"strings"
	"testing"
)

func TestAddRequirement(t *testing.T) {
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

	q, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for TestAddGetRequirement: %s", err.Error())
	}

	tc := []struct {
		Name             string
		RequirementName  string
		Notes            string
		Reference        string
		QualificationID  string
		Grade            types.Grade
		Type             types.RequirementType
		Initial          bool
		DaysValidFor     int
		ExpectedError    error
		VerificationFunc func(t *testing.T)
	}{
		{
			Name:            "Successful Add Qualification Type",
			RequirementName: "QualificationType",
			Notes:           "This is a test",
			Reference:       "This is a reference",
			QualificationID: q.ID,
			Type:            types.QualificationType,
			Initial:         true,
			ExpectedError:   nil,
		},
		{
			Name:            "Successful Add Grade Type",
			RequirementName: "Grade Type",
			Notes:           "notes",
			Reference:       "reference",
			QualificationID: "",
			Grade:           types.E7,
			Type:            types.GradeType,
			DaysValidFor:    0,
			Initial:         true,
			ExpectedError:   nil,
		},
		{
			Name:            "Successful Add Initial WBT Type",
			RequirementName: "WBT Initial Type",
			Notes:           "notes",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.WbtType,
			DaysValidFor:    0,
			Initial:         true,
			ExpectedError:   nil,
		},
		{
			Name:            "Successful Add Recurring WBT Type",
			RequirementName: "WBT Recurring Type",
			Notes:           "notes",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.WbtType,
			Initial:         false,
			DaysValidFor:    100,
			ExpectedError:   nil,
		},
		{
			Name:            "Successful Add Proficiency Type",
			RequirementName: "Proficiency Recurring",
			Notes:           "notes",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.ProficiencyType,
			DaysValidFor:    130,
			ExpectedError:   nil,
		},
		{
			Name:            "Invalid Initial Type",
			RequirementName: "Invalid Initial",
			Notes:           "",
			Reference:       "asdfasdf",
			QualificationID: "",
			Grade:           "",
			Type:            types.ProficiencyType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Invalid type for initial requirement") {
					t.Error("Expected error to contain \"Invalid type for initial requirement\", but it didn't.")
				}
			},
		},
		{
			Name:            "Invalid Recurring Type",
			RequirementName: "Invalid Recurring",
			Notes:           "",
			Reference:       "asdfasdf",
			QualificationID: "",
			Grade:           "",
			Type:            types.GradeType,
			Initial:         false,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Invalid type for recurring requirement") {
					t.Error("Expected error to contain \"Invalid type for recurring requirement\", but it didn't.")
				}
			},
		},
		{
			Name:            "Missing Initial Name and Reference",
			RequirementName: "",
			Notes:           "",
			Reference:       "",
			QualificationID: "",
			Grade:           "",
			Type:            types.WbtType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Name Missing Reference") {
					t.Error("Expected error to contain \"Missing Name Missing Reference\", but it didn't.")
				}
			},
		},
		{
			Name:            "Qualification Type Missing QualificationID",
			RequirementName: "Missing QualificationID",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.QualificationType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(*testing.T) {
				if !strings.Contains(buf.String(), "Missing QualificationID") {
					t.Error("Expected error to contain \"Missing QualificationID\", but it didn't.")
				}
			},
		},
		{
			Name:            "Qualification Type with Grade",
			RequirementName: "Qualification with grade",
			Notes:           "",
			Reference:       "reference",
			QualificationID: q.ID,
			Grade:           types.E7,
			Type:            types.QualificationType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Grade should be blank") {
					t.Error("Expected error to contain \"Grade should be blank\", but it didn't.")
				}
			},
		},
		{
			Name:            "Qualification Type with DaysValidFor",
			RequirementName: "Qualification with DaysValidFor",
			Notes:           "",
			Reference:       "reference",
			QualificationID: q.ID,
			Grade:           "",
			Type:            types.QualificationType,
			Initial:         true,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "DaysValidFor should be 0") {
					t.Error("Expected error to contain \"DaysValidFor should be 0\", but it didn't.")
				}
			},
		},
		{
			Name:            "Initial Wbt Type with QualificationID",
			RequirementName: "Wbt with Qualification",
			Notes:           "",
			Reference:       "reference",
			QualificationID: q.ID,
			Grade:           "",
			Type:            types.WbtType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "QualificationID should be blank") {
					t.Error("Expected error to contain \"QualificationID should be blank\", but it didn't.")
				}
			},
		},
		{
			Name:            "Initial Wbt Type with Grade",
			RequirementName: "Wbt with Grade",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           types.E8,
			Type:            types.WbtType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Grade should be blank") {
					t.Error("Expected error to contain \"Grade should be blank\", but it didn't.")
				}
			},
		},
		{
			Name:            "Initial Wbt Type with DaysValidFor",
			RequirementName: "Wbt with DaysValidFor",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.WbtType,
			Initial:         true,
			DaysValidFor:    199,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "DaysValidFor should be 0") {
					t.Error("Expected error to contain \"DaysValidFor should be 0\", but it didn't.")
				}
			},
		},
		{
			Name:            "Grade Type with QualificationID",
			RequirementName: "Grade with QualificationID",
			Notes:           "",
			Reference:       "reference",
			QualificationID: q.ID,
			Grade:           types.E8,
			Type:            types.GradeType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "QualificationID should be blank") {
					t.Error("Expected error to contain \"QualificationID should be blank\", but it didn't.")
				}
			},
		},
		{
			Name:            "Grade Type with DaysValidFor",
			RequirementName: "Grade with DaysValidFor",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           types.E7,
			Type:            types.GradeType,
			Initial:         true,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "DaysValidFor should be 0") {
					t.Error("Expected error to contain \"DaysValidFor should be 0\", but it didn't.")
				}
			},
		},
		{
			Name:            "Grade Type missing Grade",
			RequirementName: "Grade no Grade",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.GradeType,
			Initial:         true,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Grade") {
					t.Error("Expected error to contain \"Missing Grade\", but it didn't.")
				}
			},
		},
		{
			Name:            "Recurring Missing Name, Reference, DaysValidFor",
			RequirementName: "",
			Notes:           "",
			Reference:       "",
			QualificationID: "",
			Grade:           "",
			Type:            types.WbtType,
			Initial:         false,
			DaysValidFor:    0,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Name Missing Reference Missing DaysValidFor") {
					t.Error("Expected error to contain \"Missing Name Missing Reference Missing DaysValidFor\", but it didn't.")
				}
			},
		},
		{
			Name:            "Recurring Wbt with QualificationID",
			RequirementName: "Recurring Wbt with QualificationID",
			Notes:           "",
			Reference:       "reference",
			QualificationID: q.ID,
			Grade:           "",
			Type:            types.WbtType,
			Initial:         false,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "QualificationID should be blank") {
					t.Errorf("Expected error to contain \"QualificationID should be blank\", but it didn't")
				}
			},
		},
		{
			Name:            "Recurring Wbt with Grade",
			RequirementName: "Recurring Wbt with Grade",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           types.E7,
			Type:            types.WbtType,
			Initial:         false,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Grade should be blank") {
					t.Errorf("Expected error to contain \"Grade should be blank\", but it didn't")
				}
			},
		},
		{
			Name:            "Proficiency Type without notes",
			RequirementName: "Proficiency no notes",
			Notes:           "",
			Reference:       "reference",
			QualificationID: "",
			Grade:           "",
			Type:            types.ProficiencyType,
			Initial:         false,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Notes") {
					t.Errorf("Expected error to contain \"Missing Notes\", but it didn't")
				}
			},
		},
		{
			Name:            "Proficiency Type with QualificationID",
			RequirementName: "Proficiency with QualificationID",
			Notes:           "notes",
			Reference:       "reference",
			QualificationID: q.ID,
			Grade:           "",
			Type:            types.ProficiencyType,
			Initial:         false,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "QualificationID should be blank") {
					t.Errorf("Expected error to contain \"QualificationID should be blank\", but it didn't")
				}
			},
		},
		{
			Name:            "Proficiency Type with Grade",
			RequirementName: "Proficiency with Grade",
			Notes:           "notes",
			Reference:       "reference",
			QualificationID: "",
			Grade:           types.E8,
			Type:            types.ProficiencyType,
			Initial:         false,
			DaysValidFor:    100,
			ExpectedError:   backend.ErrValidation,
			VerificationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Grade should be blank") {
					t.Errorf("Expected error to contain \"Grade should be blank\", but it didn't")
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			buf.Reset()
			_, err := b.AddRequirement(types.Requirement{
				Name:            tt.RequirementName,
				Reference:       tt.Reference,
				Initial:         tt.Initial,
				QualificationID: tt.QualificationID,
				Grade:           tt.Grade,
				Notes:           tt.Notes,
				DaysValidFor:    tt.DaysValidFor,
				Type:            tt.Type,
			})
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.VerificationFunc != nil {
				tt.VerificationFunc(t)
			}
		})
	}
}

func TestGetRequirement(t *testing.T) {
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

	initial, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for TestGetRequirement: %s\n", err.Error())
	}
	recurring, err := b.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for TestGetRequirement: %s\n", err.Error())
	}

	tc := []struct {
		Name                string
		ID                  string
		ExpectedRequirement types.Requirement
		ExpectedError       error
	}{
		{
			Name:                "Successful Initial Get",
			ID:                  initial.ID,
			ExpectedRequirement: initial,
			ExpectedError:       nil,
		},
		{
			Name:                "Successful Recurring Get",
			ID:                  recurring.ID,
			ExpectedRequirement: recurring,
			ExpectedError:       nil,
		},
		{
			Name:                "Requirement Not Found",
			ID:                  uuid.NewString(),
			ExpectedRequirement: types.Requirement{},
			ExpectedError:       backend.ErrRequirementNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			req, err := b.GetRequirement(tt.ID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s\n", tt.ExpectedError.Error(), err.Error())
			}
			if err == nil && tt.ExpectedRequirement != req {
				t.Errorf("Expected requirement: %+v\nGot: %+v\n", tt.ExpectedRequirement, req)
			}
		})
	}
}

func TestGetAllRequirements(t *testing.T) {
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

	var expectedRequirements []types.Requirement
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
				req1, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
				if err != nil {
					t.Fatalf("Error adding requirement for TestGetAllRequirements: %s", err.Error())
				}
				expectedRequirements = append(expectedRequirements, req1)
				tt.Results = expectedRequirements
			},
		},
		{
			Name:          "Two results",
			Results:       []types.Requirement{},
			ExpectedError: nil,
			SetupFunc: func(t *testing.T, tt *testStruct) {
				req2, err := b.AddRequirement(testutils.RandomRecurringRequirement())
				if err != nil {
					t.Fatalf("Error adding requirement for TestGetAllRequirements: %s", err.Error())
				}
				expectedRequirements = append(expectedRequirements, req2)
				tt.Results = expectedRequirements
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

func TestStandardUpdateRequirement(t *testing.T) {
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

	q, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding Qualification for TestUpdateRequirement: %s\n", err.Error())
	}
	originalQualificationRequirement, err := b.AddRequirement(types.Requirement{
		Name:            "Original",
		Initial:         true,
		Reference:       "Original",
		Notes:           "Original",
		QualificationID: q.ID,
		Type:            types.QualificationType,
	})
	if err != nil {
		t.Fatalf("Error adding qualification requirement for TestUpdateRequirement: %s\n", err.Error())
	}
	originalGradeRequirement, err := b.AddRequirement(types.Requirement{
		Name:      "Original",
		Reference: "Original",
		Initial:   true,
		Notes:     "Original",
		Grade:     types.E1,
		Type:      types.GradeType,
	})
	if err != nil {
		t.Fatalf("Error adding qualification requirement for TestUpdateRequirement: %s\n", err.Error())
	}
	originalInitialWbtRequirement, err := b.AddRequirement(types.Requirement{
		Name:      "Original",
		Initial:   true,
		Reference: "Original",
		Notes:     "Original",
		Type:      types.WbtType,
	})
	if err != nil {
		t.Fatalf("Error adding qualification requirement for TestUpdateRequirement: %s\n", err.Error())
	}
	originalRecurringWbtRequirement, err := b.AddRequirement(types.Requirement{
		Name:         "Original",
		Reference:    "Original",
		Notes:        "Original",
		DaysValidFor: 100,
		Type:         types.WbtType,
	})
	if err != nil {
		t.Fatalf("Error adding qualification requirement for TestUpdateRequirement: %s\n", err.Error())
	}
	originalProficiencyRequirement, err := b.AddRequirement(types.Requirement{
		Name:         "Original",
		Reference:    "Original",
		Notes:        "Original",
		DaysValidFor: 100,
		Type:         types.ProficiencyType,
	})
	if err != nil {
		t.Fatalf("Error adding qualification requirement for TestUpdateRequirement: %s\n", err.Error())
	}
	tc := []struct {
		Name           string
		Update         types.Requirement
		ExpectedError  error
		ValidationFunc func(t *testing.T)
	}{
		{
			Name: "Successful update qualification",
			Update: types.Requirement{
				ID:              originalQualificationRequirement.ID,
				Initial:         true,
				Name:            "New Name",
				Reference:       "New Reference",
				Notes:           "New Notes",
				QualificationID: q.ID,
				Type:            types.QualificationType,
			},
			ExpectedError: nil,
		},
		{
			Name: "Successful update grade",
			Update: types.Requirement{
				ID:        originalGradeRequirement.ID,
				Initial:   true,
				Name:      "New Name",
				Reference: "New Reference",
				Grade:     types.E9,
				Notes:     "New Notes",
				Type:      types.GradeType,
			},
			ExpectedError: nil,
		},
		{
			Name: "Successful update wbt initial",
			Update: types.Requirement{
				ID:        originalInitialWbtRequirement.ID,
				Initial:   true,
				Name:      "New Name",
				Reference: "New Reference",
				Notes:     "New Notes",
				Type:      types.WbtType,
			},
			ExpectedError: nil,
		},
		{
			Name: "Successful update wbt recurring",
			Update: types.Requirement{
				ID:           originalRecurringWbtRequirement.ID,
				Name:         "New Name",
				Reference:    "New Reference",
				Notes:        "New Notes",
				DaysValidFor: 200,
				Type:         types.WbtType,
			},
			ExpectedError: nil,
		},
		{
			Name: "Successful update proficiency",
			Update: types.Requirement{
				ID:           originalProficiencyRequirement.ID,
				Name:         "New Name",
				Reference:    "New Reference",
				Notes:        "New Notes",
				DaysValidFor: 200,
				Type:         types.ProficiencyType,
			},
			ExpectedError: nil,
		},
		{
			Name: "Invalid Qualification Update",
			Update: types.Requirement{
				ID:              originalQualificationRequirement.ID,
				Initial:         true,
				Name:            "",
				Reference:       "",
				QualificationID: "",
				Type:            types.QualificationType,
			},
			ExpectedError: backend.ErrValidation,
			ValidationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Name") {
					t.Errorf("Expected error to contain string \"Missing Name\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing Reference") {
					t.Errorf("Expected error to contain string \"Missing Reference\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing QualificationID") {
					t.Errorf("Expected error to contain string \"Missing QualificationID\", but it didn't")
				}
			},
		},
		{
			Name: "Invalid Grade Update",
			Update: types.Requirement{
				ID:        originalGradeRequirement.ID,
				Initial:   true,
				Name:      "",
				Reference: "",
				Grade:     "",
				Type:      types.GradeType,
			},
			ExpectedError: backend.ErrValidation,
			ValidationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Name") {
					t.Errorf("Expected error to contain string \"Missing Name\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing Reference") {
					t.Errorf("Expected error to contain string \"Missing Reference\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing Grade") {
					t.Errorf("Expected error to contain string \"Missing Grade\", but it didn't")
				}
			},
		},
		{
			Name: "Invalid Initial Wbt Type",
			Update: types.Requirement{
				ID:        originalInitialWbtRequirement.ID,
				Initial:   true,
				Name:      "",
				Reference: "",
				Type:      types.WbtType,
			},
			ExpectedError: backend.ErrValidation,
			ValidationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Name") {
					t.Errorf("Expected error to contain string \"Missing Name\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing Reference") {
					t.Errorf("Expected error to contain string \"Missing Reference\", but it didn't")
				}
			},
		},
		{
			Name: "Invalid Recurring Wbt Type",
			Update: types.Requirement{
				ID:           originalRecurringWbtRequirement.ID,
				Name:         "",
				Reference:    "",
				DaysValidFor: 0,
				Type:         types.WbtType,
			},
			ExpectedError: backend.ErrValidation,
			ValidationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Missing Name") {
					t.Errorf("Expected error to contain string \"Missing Name\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing Reference") {
					t.Errorf("Expected error to contain string \"Missing Reference\", but it didn't")
				}
				if !strings.Contains(buf.String(), "Missing DaysValidFor") {
					t.Errorf("Expected error to contain string \"Missing DaysValidFor\", but it didn't")
				}
			},
		},
		{
			Name: "Requirement not found",
			Update: types.Requirement{
				ID:        uuid.NewString(),
				Initial:   true,
				Name:      "Test",
				Reference: "test",
				Type:      types.WbtType,
			},
			ExpectedError:  backend.ErrRequirementNotFound,
			ValidationFunc: nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			buf.Reset()
			_, err := b.UpdateRequirement(tt.Update)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				if err == nil {
					t.Errorf("Expected error: %s\nBut got none", tt.ExpectedError.Error())
				} else {
					t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
				}
			}
			if tt.ExpectedError == nil {
				got, err := b.GetRequirement(tt.Update.ID)
				if err != nil {
					t.Errorf("Error getting requirement for TestUpdateRequirement: %s", err.Error())
				}
				if !testutils.CompareRequirements(got, tt.Update) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.Update, got)
				}
				if tt.ValidationFunc != nil {
					tt.ValidationFunc(t)
				}
			}
		})
	}
}

func TestChangeRequirementType(t *testing.T) {
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
	q, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for TestChangeRequirementType: %s\n", err.Error())
	}
	originalInitialRequirement, err := b.AddRequirement(types.Requirement{
		Name:            "Original",
		Initial:         true,
		Reference:       "Original",
		QualificationID: q.ID,
		Type:            types.QualificationType,
	})
	if err != nil {
		t.Fatalf("Error adding requirement for TestChangeRequirementType: %s\n", err.Error())
	}

	originalRecurringRequirement, err := b.AddRequirement(types.Requirement{
		Name:         "Original",
		Reference:    "Original",
		Notes:        "Original",
		DaysValidFor: 100,
		Type:         types.ProficiencyType,
	})
	if err != nil {
		t.Fatalf("Error adding requirement for TestChangeRequirementType: %s\n", err.Error())
	}
	tc := []struct {
		Name           string
		Updates        types.Requirement
		ExpectedError  error
		ValidationFunc func(t *testing.T)
	}{
		{
			Name: "Valid Initial Change",
			Updates: types.Requirement{
				ID:        originalInitialRequirement.ID,
				Initial:   true,
				Name:      "New Name",
				Reference: "New Reference",
				Grade:     types.E8,
				Type:      types.GradeType,
			},
			ExpectedError:  nil,
			ValidationFunc: nil,
		},
		{
			Name: "Valid Recurring Change",
			Updates: types.Requirement{
				ID:           originalRecurringRequirement.ID,
				Name:         "New Name",
				Reference:    "New Reference",
				DaysValidFor: 900,
				Type:         types.WbtType,
			},
			ExpectedError:  nil,
			ValidationFunc: nil,
		},
		{
			Name: "Initial to Recurring change",
			Updates: types.Requirement{
				ID:           originalInitialRequirement.ID,
				Name:         "Switching Types",
				Reference:    "Reference",
				Notes:        "New Notes",
				DaysValidFor: 1000,
				Type:         types.ProficiencyType,
			},
			ExpectedError: backend.ErrValidation,
			ValidationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Changing from initial to recurring type not allowed") {
					t.Errorf("Expected logs to contain string \"Changing from initial to recurring type not allowed\", but didn't find it.")
				}
			},
		},
		{
			Name: "Recurring to Initial change",
			Updates: types.Requirement{
				ID:        originalRecurringRequirement.ID,
				Name:      "Switching Types",
				Initial:   true,
				Reference: "Reference",
				Type:      types.GradeType,
				Grade:     types.E8,
			},
			ExpectedError: backend.ErrValidation,
			ValidationFunc: func(t *testing.T) {
				if !strings.Contains(buf.String(), "Changing from recurring to initial type not allowed") {
					t.Errorf("Expected logs to contain string \"Changing from recurring to initial type not allowed\", but didn't find it.")
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := b.UpdateRequirement(tt.Updates)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s\n", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				if got != tt.Updates {
					t.Errorf("Expected requirement: %+v\nGot: %+v\n", tt.Updates, got)
				}
			}
			if tt.ValidationFunc != nil {
				tt.ValidationFunc(t)
			}
		})
	}
}

func TestDeleteRequirement(t *testing.T) {
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

	q, err := b.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding Qualification for TestDeleteRequirement")
	}

	initialRequirement, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteRequirement: %s", err.Error())
	}

	initialRequirement2, err := b.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteRequirement: %s", err.Error())
	}

	recurringRequirement, err := b.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteRequirement: %s", err.Error())
	}

	recurringRequirement2, err := b.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for TestDeleteRequirement: %s", err.Error())
	}

	err = b.AssignRequirementToQualification(q.ID, initialRequirement.ID, initialRequirement.Initial)
	if err != nil {
		t.Fatalf("Error assigning initial requirement to qualification for TestDeleteRequirement: %s", err.Error())
	}
	err = b.AssignRequirementToQualification(q.ID, recurringRequirement.ID, recurringRequirement.Initial)
	if err != nil {
		t.Fatalf("Error assigning recurring requirement to qualification for TestDeleteRequirement: %s", err.Error())
	}

	tc := []struct {
		Name             string
		RequirementID    string
		ExpectedError    error
		VerificationFunc func(t *testing.T)
	}{
		{
			Name:          "Successful initial delete",
			RequirementID: initialRequirement2.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Successful recurring delete",
			RequirementID: recurringRequirement2.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Requirement not found",
			RequirementID: uuid.NewString(),
			ExpectedError: backend.ErrRequirementNotFound,
		},
		{
			Name:          "Initial requirement assigned to qualification",
			RequirementID: initialRequirement.ID,
			ExpectedError: nil,
			VerificationFunc: func(t *testing.T) {
				qual, err := b.GetQualification(q.ID)
				if err != nil {
					t.Fatalf("Error getting qualification during validation: %s", err.Error())
				}
				if len(qual.InitialRequirements) != 0 {
					t.Errorf("Expected qualification to have no initial requirements but found: %d", len(qual.InitialRequirements))
				}
			},
		},
		{
			Name:          "Recurring requirement assigned to qualification",
			RequirementID: recurringRequirement.ID,
			ExpectedError: nil,
			VerificationFunc: func(t *testing.T) {
				qual, err := b.GetQualification(q.ID)
				if err != nil {
					t.Fatalf("Error getting qualification during validation: %s", err.Error())
				}
				if len(qual.RecurringRequirements) != 0 {
					t.Errorf("Expected qualification to have no recurring requirements but found: %d", len(qual.RecurringRequirements))
				}
			},
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
