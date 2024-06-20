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

func TestAddQualification(t *testing.T) {
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
		qualification types.Qualification
		expectedError error
	}{
		{
			name: "Successful add",
			qualification: types.Qualification{
				ID:                    uuid.NewString(),
				Name:                  "Test qual",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "test notes",
				References:            nil,
				Expires:               true,
				ExpirationDays:        365,
			},
			expectedError: nil,
		},
		{
			name:          "Empty Qualification",
			qualification: types.Qualification{},
			expectedError: types.ErrMissingArgs,
		},
		{
			name: "Missing ID",
			qualification: types.Qualification{
				ID:                    "",
				Name:                  "Test",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "",
				References:            nil,
				Expires:               true,
				ExpirationDays:        -1,
			},
			expectedError: types.ErrMissingArgs,
		},
		{
			name: "Missing Name",
			qualification: types.Qualification{
				ID:                    uuid.NewString(),
				Name:                  "",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "",
				References:            nil,
				Expires:               false,
				ExpirationDays:        0,
			},
			expectedError: types.ErrMissingArgs,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.AddQualification(tt.qualification)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				qual, err := backend.GetQualification(tt.qualification.ID)
				if err != nil {
					t.Errorf("Error getting qualification to validate TestAddQualification: %s", err.Error())
				}
				if !reflect.DeepEqual(qual, tt.qualification) {
					t.Errorf("Expected qualification: %+v\nGot: %+v", tt.qualification, qual)
				}
			}
		})
	}
}

func TestGetQualification(t *testing.T) {
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

	testQual := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Test Qual",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "test notes",
		References:            nil,
		Expires:               true,
		ExpirationDays:        -1,
	}
	err = backend.AddQualification(testQual)
	if err != nil {
		t.Fatalf("Error creating qualification for TestGetQualification: %s", err.Error())
	}

	tc := []struct {
		name           string
		id             string
		expectedResult types.Qualification
		expectedError  error
	}{
		{
			name:           "Successful get",
			id:             testQual.ID,
			expectedResult: testQual,
			expectedError:  nil,
		},
		{
			name:           "Not found",
			id:             uuid.NewString(),
			expectedResult: types.Qualification{},
			expectedError:  types.ErrQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			qual, err := backend.GetQualification(tt.id)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil && !reflect.DeepEqual(qual, tt.expectedResult) {
				t.Errorf("Expected result: %+v\nGot: %+v", tt.expectedResult, qual)
			}
		})
	}
}

func TestGetAllQualifications(t *testing.T) {
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
	qual1 := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Qual 1",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "test notes",
		References:            nil,
		Expires:               false,
		ExpirationDays:        0,
	}
	qual2 := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Qual 2",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "Even more notes",
		References:            nil,
		Expires:               true,
		ExpirationDays:        99999,
	}

	tc := []struct {
		name           string
		expectedResult []types.Qualification
		expectedError  error
		setupFunc      func()
	}{
		{
			name:           "Successful empty get",
			expectedResult: []types.Qualification{},
			expectedError:  nil,
			setupFunc:      func() {},
		},
		{
			name:           "Successful single item get",
			expectedResult: []types.Qualification{qual1},
			expectedError:  nil,
			setupFunc: func() {
				if err := backend.AddQualification(qual1); err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQualifictions: %s", err.Error())
				}
			},
		},
		{
			name:           "Successful multi item get",
			expectedResult: []types.Qualification{qual1, qual2},
			expectedError:  nil,
			setupFunc: func() {
				if err := backend.AddQualification(qual2); err != nil {
					t.Fatalf("Error adding qualification for TestGetAllQualifictions: %s", err.Error())
				}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()
			quals, err := backend.GetAllQualifications()
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				for _, q := range tt.expectedResult {
					found := false
					for _, qu := range quals {
						if reflect.DeepEqual(q, qu) {
							found = true
						}
					}
					if !found {
						t.Errorf("Expected to find qualification but didnt: %+v", q)
					}
				}
			}
		})
	}
}

func TestUpdateQualification(t *testing.T) {
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
	qual := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Old",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "Old notes",
		References:            nil,
		Expires:               false,
		ExpirationDays:        0,
	}
	if err := backend.AddQualification(qual); err != nil {
		t.Fatalf("Error adding qualification to database for TestUpdateQualification: %s", err.Error())
	}

	tc := []struct {
		name          string
		qualification types.Qualification
		expectedError error
	}{
		{
			name:          "No Update",
			qualification: qual,
			expectedError: nil,
		},
		{
			name: "Successful full update",
			qualification: types.Qualification{
				ID:                    qual.ID,
				Name:                  "New",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "New Notes",
				References:            nil,
				Expires:               true,
				ExpirationDays:        -1,
			},
			expectedError: nil,
		},
		{
			name:          "Not found",
			qualification: types.Qualification{ID: uuid.NewString()},
			expectedError: types.ErrQualificationNotFound,
		},
		{
			name: "Invalid ExipirationDays",
			qualification: types.Qualification{
				ID:                    qual.ID,
				Name:                  "New",
				InitialRequirements:   nil,
				RecurringRequirements: nil,
				Notes:                 "New Notes",
				References:            nil,
				Expires:               true,
				ExpirationDays:        0,
			},
			expectedError: types.ErrBadUpdate,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.UpdateQualification(tt.qualification)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				q, err := backend.GetQualification(tt.qualification.ID)
				if err != nil {
					t.Errorf("Error getting qualification for verification of TestUpdateQualification: %s", err.Error())
				}
				if !reflect.DeepEqual(tt.qualification, q) {
					t.Errorf("Expected: %+v\nGot: %+v", tt.qualification, q)
				}
			}
		})
	}
}

func TestDeleteQualification(t *testing.T) {
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
	qual := types.Qualification{
		ID:                    uuid.NewString(),
		Name:                  "Test",
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 "Test",
		References:            nil,
		Expires:               false,
		ExpirationDays:        0,
	}
	if err = backend.AddQualification(qual); err != nil {
		t.Fatalf("Error inserting qualification into database for TestDeleteQualification: %s", err.Error())
	}

	tc := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "Successful delete",
			id:            qual.ID,
			expectedError: nil,
		},
		{
			name:          "User not found",
			id:            uuid.NewString(),
			expectedError: types.ErrQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := backend.DeleteQualification(tt.id)
			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.expectedError.Error(), err.Error())
			}
			if tt.expectedError == nil {
				_, err = backend.GetQualification(tt.id)
				if err == nil {
					t.Error("Expected to not find user, but did")
				} else if !errors.Is(err, types.ErrQualificationNotFound) {
					t.Errorf("Expected to not find member but got error: %s", err.Error())
				}
			}
		})
	}
}
