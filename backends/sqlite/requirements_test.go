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
	"slices"
	"testing"
)

func TestAddRequirement(t *testing.T) {
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

	tc := []struct {
		name          string
		requirement   types.Requirement
		expectedError error
	}{
		{
			name: "Successful add",
			requirement: types.Requirement{
				ID:           uuid.NewString(),
				Name:         "Test Requirement",
				Description:  "Test description",
				Notes:        "Test notes",
				DaysValidFor: 188,
			},
			expectedError: nil,
		},
		{
			name:          "Empty requirement",
			requirement:   types.Requirement{},
			expectedError: types.ErrMissingArgs,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.AddRequirement(tt.requirement)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				req, err := backend.GetRequirement(tt.requirement.ID)
				if err != nil {
					t.Errorf("Error getting requirement: %s", err.Error())
				}
				if !reflect.DeepEqual(req, tt.requirement) {
					t.Errorf("Expected requirement: %+v\nGot: %+v", tt.requirement, req)
				}
			}
		})
	}
}

func TestGetRequirement(t *testing.T) {
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
	req := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test Requirement",
		Description:  "Test description",
		Notes:        "these are some notes",
		DaysValidFor: 100,
	}
	if err := backend.AddRequirement(req); err != nil {
		t.Fatalf("Error adding requirement for TestGetRequirement: %s", err.Error())
	}
	tc := []struct {
		name           string
		requirementID  string
		expectedResult types.Requirement
		expectedError  error
	}{
		{
			name:           "Successful get",
			requirementID:  req.ID,
			expectedResult: req,
			expectedError:  nil,
		},
		{
			name:           "Requirement not found",
			requirementID:  uuid.NewString(),
			expectedResult: types.Requirement{},
			expectedError:  types.ErrRequirementNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			req, err := backend.GetRequirement(tt.requirementID)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				if !reflect.DeepEqual(tt.expectedResult, req) {
					t.Errorf("Expected result: %+v\nGot: %+v", tt.expectedResult, req)
				}
			}
		})
	}
}

func TestGetAllRequirements(t *testing.T) {
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
	req := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test Requirement",
		Description:  "Description",
		Notes:        "Notes",
		DaysValidFor: 100,
	}
	req2 := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test Requirement",
		Description:  "Description",
		Notes:        "Notes",
		DaysValidFor: 100,
	}

	tc := []struct {
		name           string
		expectedResult []types.Requirement
		expectedError  error
		setupFunc      func()
	}{
		{
			name:           "No results",
			expectedResult: []types.Requirement{},
			expectedError:  nil,
			setupFunc:      func() {},
		},
		{
			name:           "Single item return",
			expectedResult: []types.Requirement{req},
			expectedError:  nil,
			setupFunc: func() {
				if err := backend.AddRequirement(req); err != nil {
					t.Fatalf("Error adding requirement for TestGetAllRequirements: %s", err.Error())
				}
			},
		},
		{
			name:           "Multi item return",
			expectedResult: []types.Requirement{req, req2},
			expectedError:  nil,
			setupFunc: func() {
				if err := backend.AddRequirement(req2); err != nil {
					t.Fatalf("Error adding requirement for TestGetAllRequirements: %s", err.Error())
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()
			reqs, err := backend.GetAllRequirements()
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				if !slices.Equal(reqs, tt.expectedResult) {
					t.Errorf("Expected result: %+v\nGot: %+v", tt.expectedResult, reqs)
				}
			}
		})
	}
}

func TestUpdateRequirement(t *testing.T) {
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
	requirement := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test Requirement",
		Description:  "Test description",
		Notes:        "Test Notes",
		DaysValidFor: -1,
	}
	if err := backend.AddRequirement(requirement); err != nil {
		t.Fatalf("Error adding requirement for TestUpdateRequirement: %s", err.Error())
	}

	tc := []struct {
		name          string
		update        types.Requirement
		expectedError error
	}{
		{
			name: "Successful full update",
			update: types.Requirement{
				ID:           requirement.ID,
				Name:         "New Requirement",
				Description:  "New Description",
				Notes:        "New Notes",
				DaysValidFor: 365,
			},
			expectedError: nil,
		},
		{
			name: "Requirement not found",
			update: types.Requirement{
				ID:           uuid.NewString(),
				Name:         "Test Requirement",
				Description:  " Blah",
				Notes:        "Blah",
				DaysValidFor: -1,
			},
			expectedError: types.ErrRequirementNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.UpdateRequirement(tt.update)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				req, err := backend.GetRequirement(tt.update.ID)
				if err != nil {
					t.Errorf("Error getting requirement: %s", err.Error())
				}
				if req != tt.update {
					t.Errorf("Expected database requirement to be: %+v\nGot: %+v", tt.update, req)
				}
			}
		})
	}
}

func TestDeleteRequirement(t *testing.T) {
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
	requirement := types.Requirement{
		ID:           uuid.NewString(),
		Name:         "Test Requirement",
		Description:  "Test description",
		Notes:        "Test Notes",
		DaysValidFor: -1,
	}
	if err := backend.AddRequirement(requirement); err != nil {
		t.Fatalf("Error adding requirement for TestUpdateRequirement: %s", err.Error())
	}

	tc := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "Successful delete",
			id:            requirement.ID,
			expectedError: nil,
		},
		{
			name:          "Requirement not found",
			id:            uuid.NewString(),
			expectedError: types.ErrRequirementNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.DeleteRequirement(tt.id)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				if _, err := backend.GetRequirement(tt.id); !errors.Is(err, types.ErrRequirementNotFound) {
					t.Errorf("Expected to get requirement not found, but didn't")
				}
			}
		})
	}
}
