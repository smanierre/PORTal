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

func TestAddGetReference(t *testing.T) {
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
		t.Fatalf("Error adding reference for TestAddGetReference: %s", err.Error())
	}

	// Test add duplicate reference name
	_, err = b.AddReference(types.Reference{
		Name:      ref1.Name,
		Volume:    0,
		Paragraph: "asdf",
	})
	if !errors.Is(err, backend.ErrDuplicateReference) {
		t.Errorf("Expected error: %s\nGot: %s", backend.ErrDuplicateReference, err)
	}

	tc := []struct {
		Name              string
		ReferenceID       string
		ExpectedReference types.Reference
		ExpectedError     error
	}{
		{
			Name:              "Successful get",
			ReferenceID:       ref1.ID,
			ExpectedReference: ref1,
			ExpectedError:     nil,
		},
		{
			Name:              "Reference not found",
			ReferenceID:       uuid.NewString(),
			ExpectedReference: types.Reference{},
			ExpectedError:     backend.ErrReferenceNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			ref, err := b.GetReference(tt.ReferenceID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil && !reflect.DeepEqual(ref, tt.ExpectedReference) {
				t.Errorf("Expected: %+v\nGot: %+v", tt.ExpectedReference, ref)
			}
		})
	}
}

func TestGetReferences(t *testing.T) {
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

	ref1 := testutils.RandomReference()
	ref2 := testutils.RandomReference()

	type testStruct struct {
		Name               string
		ExpectedReferences []types.Reference
		SetupFunc          func(t *testing.T, tt *testStruct)
		ExpectedError      error
	}
	tc := []testStruct{
		{
			Name:               "No results",
			ExpectedReferences: []types.Reference{},
			SetupFunc:          func(_ *testing.T, _ *testStruct) {},
			ExpectedError:      nil,
		},
		{
			Name:               "One result",
			ExpectedReferences: []types.Reference{},
			SetupFunc: func(t *testing.T, tt *testStruct) {
				ref1, err = b.AddReference(ref1)
				if err != nil {
					t.Errorf("Error adding reference for TestGetReferences: %s", err.Error())
				}
				tt.ExpectedReferences = append(tt.ExpectedReferences, ref1)
			},
			ExpectedError: nil,
		},
		{
			Name:               "Two results",
			ExpectedReferences: []types.Reference{},
			SetupFunc: func(t *testing.T, tt *testStruct) {
				ref2, err = b.AddReference(ref2)
				if err != nil {
					t.Errorf("Error adding reference for TestGetReferences: %s", err.Error())
				}
				tt.ExpectedReferences = append(tt.ExpectedReferences, ref1, ref2)
			},
			ExpectedError: nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			tt.SetupFunc(t, &tt)
			refs, err := b.GetReferences()
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil && !slices.Equal(tt.ExpectedReferences, refs) {
				t.Errorf("Expected: %+v\nGot: %+v", tt.ExpectedReferences, refs)
			}
		})
	}
}

func TestUpdateReferences(t *testing.T) {
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
		t.Fatalf("Error adding reference for TestUpdateReferences: %s", err.Error())
	}

	tc := []struct {
		Name             string
		NewRef           types.Reference
		OverrideNoVolume bool
		ExpectedError    error
	}{
		{
			Name: "Successful full update",
			NewRef: types.Reference{
				ID:        ref1.ID,
				Name:      "New Name",
				Volume:    904985,
				Paragraph: "New paragraph",
			},
			ExpectedError: nil,
		},
		{
			Name: "Single field update",
			NewRef: types.Reference{
				ID:     ref1.ID,
				Volume: 0,
			},
			OverrideNoVolume: true,
			ExpectedError:    nil,
		},
		{
			Name: "Reference not found",
			NewRef: types.Reference{
				ID:        uuid.NewString(),
				Name:      "idk",
				Volume:    23,
				Paragraph: "Subsection 69",
			},
			OverrideNoVolume: false,
			ExpectedError:    backend.ErrReferenceNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			ref, err := b.UpdateReference(tt.NewRef, tt.OverrideNoVolume)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				got, err := b.GetReference(tt.NewRef.ID)
				if err != nil {
					t.Errorf("Error getting reference from database: %s", err.Error())
				}
				if !reflect.DeepEqual(got, ref) {
					t.Errorf("Expected reference: %+v\nGot: %+v", ref, got)
				}
			}
		})
	}
}

func TestDeleteReference(t *testing.T) {
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

	ref, err := b.AddReference(testutils.RandomReference())
	if err != nil {
		t.Fatalf("Error adding reference for TestDeleteReference: %s", err.Error())
	}

	tc := []struct {
		Name          string
		ReferenceID   string
		ExpectedError error
	}{
		{
			Name:          "Successful delete",
			ReferenceID:   ref.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Reference not found",
			ReferenceID:   uuid.NewString(),
			ExpectedError: backend.ErrReferenceNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err := b.DeleteReference(tt.ReferenceID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				_, err := b.GetReference(tt.ReferenceID)
				if !errors.Is(err, backend.ErrReferenceNotFound) {
					t.Errorf("Expected to get: %s\nGot: %s", backend.ErrReferenceNotFound.Error(), err)
				}
			}
		})
	}
}
