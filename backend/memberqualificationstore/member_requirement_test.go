package memberqualificationstore_test

import (
	"PORTal/backend"
	"PORTal/backend/memberqualificationstore"
	"PORTal/backend/memberstore"
	"PORTal/backend/qualificationstore"
	"PORTal/providers/sqlite/memberprovider"
	"PORTal/providers/sqlite/qualificationprovider"
	"PORTal/testutils"
	"PORTal/types"
	"bytes"
	"errors"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"testing"
)

func TestGetInitialMemberRequirement(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for tests: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)
	qualificationProvider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating qualification provider for tests: %s", err.Error())
	}
	qualificationStore := qualificationstore.New(qualificationProvider, logger)
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, nil)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	m1 := testutils.RandomMember(false)
	m1.SupervisorID = supervisor.ID
	m1, err = memberStore.AddMember(m1)
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for tests: %s", err.Error())
	}
	ir1, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, ir1.ID, true)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	q2, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for tests: %s", err.Error())
	}
	ir2, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(true, q2.ID))
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, ir2.ID, true)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	rr1, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, rr1.ID, false)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	rr2, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, rr2.ID, false)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification for tests: %s", err.Error())
	}

	tc := []struct {
		Name              string
		MemberID          string
		RequirementID     string
		MemberRequirement types.InitialMemberRequirement
		ExpectedError     error
	}{
		{
			Name:          "Successful Initial Get",
			MemberID:      m1.ID,
			RequirementID: ir1.ID,
			MemberRequirement: types.InitialMemberRequirement{
				MemberID:      m1.ID,
				RequirementID: ir1.ID,
				CompletedDate: types.Never,
				AssignedBy:    supervisor.ID,
				CompletedBy:   "",
			},
			ExpectedError: nil,
		},
		{
			Name:              "Member not found",
			MemberID:          uuid.NewString(),
			RequirementID:     ir1.ID,
			MemberRequirement: types.InitialMemberRequirement{},
			ExpectedError:     backend.ErrInitialMemberRequirementNotFound,
		},
		{
			Name:              "Requirement not found",
			MemberID:          m1.ID,
			RequirementID:     uuid.NewString(),
			MemberRequirement: types.InitialMemberRequirement{},
			ExpectedError:     backend.ErrInitialMemberRequirementNotFound,
		},
	}
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			mr, err := mqStore.GetInitialMemberRequirement(tt.MemberID, tt.RequirementID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && err == nil {
				t.Errorf("No error return when expected error was: %s\n", tt.ExpectedError.Error())
			}
			if err != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error was: %s\nGot: %s\n", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				if mr != tt.MemberRequirement {
					t.Errorf("Expected value: %+v\n but got %+v\n", tt.MemberRequirement, mr)
				}
			}
		})
	}
}

func TestCompleteInitialMemberRequirement(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for tests: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)
	qualificationProvider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating qualification provider for tests: %s", err.Error())
	}
	qualificationStore := qualificationstore.New(qualificationProvider, logger)
	constClock := constantTimeClock{}
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, constClock)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	m1 := testutils.RandomMember(false)
	m1.SupervisorID = supervisor.ID
	m1, err = memberStore.AddMember(m1)
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for tests: %s", err.Error())
	}
	ir1, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, ir1.ID, true)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	rr1, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, rr1.ID, false)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification for tests: %s", err.Error())
	}

	tc := []struct {
		Name          string
		MemberID      string
		RequirementID string
		CompletedBy   string
		ExpectedError error
	}{
		{
			Name:          "Successful Complete",
			MemberID:      m1.ID,
			RequirementID: ir1.ID,
			CompletedBy:   supervisor.ID,
			ExpectedError: nil,
		},
		{
			Name:          "Member Not found",
			MemberID:      uuid.NewString(),
			RequirementID: ir1.ID,
			CompletedBy:   supervisor.ID,
			ExpectedError: backend.ErrInitialMemberRequirementNotFound,
		},
		{
			Name:          "Requirement not found",
			MemberID:      m1.ID,
			RequirementID: uuid.NewString(),
			CompletedBy:   supervisor.ID,
			ExpectedError: backend.ErrInitialMemberRequirementNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err = mqStore.CompleteInitialMemberRequirement(tt.MemberID, tt.RequirementID, tt.CompletedBy)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && err == nil {
				t.Errorf("No error return when expected error was: %s\n", tt.ExpectedError.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error was: %s\nGot: %s\n", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				mr, err := mqStore.GetInitialMemberRequirement(tt.MemberID, tt.RequirementID)
				if err != nil {
					t.Fatalf("Error getting initial member requirement for verification: %s\n", err.Error())
				}
				if mr.CompletedBy != tt.CompletedBy {
					t.Errorf("Expected \"CompletedBy\" value: %+v\n but got %+v\n", tt.CompletedBy, mr.CompletedBy)
				}
				if mr.CompletedDate != constClock.Now() {
					t.Errorf("Expected \"CompletedDate\" value: %+v\n but got %+v\n", constClock.Now(), mr.CompletedDate)
				}
			}

		})
	}
}

func TestGetRecurringMemberRequirement(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for tests: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)
	qualificationProvider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating qualification provider for tests: %s", err.Error())
	}
	qualificationStore := qualificationstore.New(qualificationProvider, logger)
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, nil)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	m1 := testutils.RandomMember(false)
	m1.SupervisorID = supervisor.ID
	m1, err = memberStore.AddMember(m1)
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for tests: %s", err.Error())
	}
	ir1, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, ir1.ID, true)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	q2, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for tests: %s", err.Error())
	}
	ir2, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(true, q2.ID))
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, ir2.ID, true)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	rr1, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, rr1.ID, false)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	rr2, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, rr2.ID, false)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification for tests: %s", err.Error())
	}

	tc := []struct {
		Name              string
		MemberID          string
		RequirementID     string
		MemberRequirement types.RecurringMemberRequirement
		ExpectedError     error
	}{
		{
			Name:          "Successful Recurring Get",
			MemberID:      m1.ID,
			RequirementID: rr1.ID,
			MemberRequirement: types.RecurringMemberRequirement{
				ID:                "",
				MemberID:          m1.ID,
				RequirementID:     rr1.ID,
				AssignedBy:        supervisor.ID,
				CompletionHistory: nil,
			},
			ExpectedError: nil,
		},
		{
			Name:              "Member not found",
			MemberID:          uuid.NewString(),
			RequirementID:     ir1.ID,
			MemberRequirement: types.RecurringMemberRequirement{},
			ExpectedError:     backend.ErrRecurringMemberRequirementNotFound,
		},
		{
			Name:              "Requirement not found",
			MemberID:          m1.ID,
			RequirementID:     uuid.NewString(),
			MemberRequirement: types.RecurringMemberRequirement{},
			ExpectedError:     backend.ErrRecurringMemberRequirementNotFound,
		},
	}
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			mr, err := mqStore.GetRecurringMemberRequirement(tt.MemberID, tt.RequirementID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && err == nil {
				t.Errorf("No error return when expected error was: %s\n", tt.ExpectedError.Error())
			}
			if err != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error was: %s\nGot: %s\n", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				// We don't quite care about the ID
				mr.ID = tt.MemberRequirement.ID
				if !reflect.DeepEqual(mr, tt.MemberRequirement) {
					t.Errorf("Expected value: %+v\n but got %+v\n", tt.MemberRequirement, mr)
				}
			}
		})
	}
}

func TestCompleteRecurringMemberRequirement(t *testing.T) {
	dbString := testutils.GetDbString()
	t.Cleanup(func() {
		os.Remove(dbString)
	})
	buf := &bytes.Buffer{}
	mr := io.MultiWriter(os.Stdout, buf)
	logger := slog.New(slog.NewTextHandler(mr, nil))
	memberProvider, err := memberprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating member provider for tests: %s", err.Error())
	}
	memberStore := memberstore.New(memberProvider, 4, logger)
	qualificationProvider, err := qualificationprovider.New(dbString, logger)
	if err != nil {
		t.Fatalf("Error creating qualification provider for tests: %s", err.Error())
	}
	qualificationStore := qualificationstore.New(qualificationProvider, logger)
	constClock := constantTimeClock{}
	constClock2 := constantTimeClock{
		t: "2000-01-01 02:00:00",
	}
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, constClock)
	mqStore2 := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, constClock2)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	m1 := testutils.RandomMember(false)
	m1.SupervisorID = supervisor.ID
	m1, err = memberStore.AddMember(m1)
	if err != nil {
		t.Fatalf("Error adding member for tests: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification for tests: %s", err.Error())
	}
	ir1, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, ir1.ID, true)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	rr1, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement for tests: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, rr1.ID, false)
	if err != nil {
		t.Fatalf("Error adding requirement to qualification for tests: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification for tests: %s", err.Error())
	}

	tc := []struct {
		Name                string
		MemberID            string
		RequirementID       string
		CompletedBy         string
		SecondStore         bool
		ExpectedRequirement types.RecurringMemberRequirement
		ExpectedError       error
	}{
		{
			Name:          "Successful initial complete",
			MemberID:      m1.ID,
			RequirementID: rr1.ID,
			CompletedBy:   supervisor.ID,
			SecondStore:   false,
			ExpectedRequirement: types.RecurringMemberRequirement{
				ID:            "",
				MemberID:      m1.ID,
				RequirementID: rr1.ID,
				AssignedBy:    supervisor.ID,
				CompletionHistory: []types.RecurringMemberRequirementCompletion{
					types.RecurringMemberRequirementCompletion{
						ID:             "",
						CompletionDate: constClock.Now(),
						CompletedBy:    supervisor.ID,
					},
				},
			},
			ExpectedError: nil,
		},
		{
			Name:          "Successful Second completion",
			MemberID:      m1.ID,
			RequirementID: rr1.ID,
			CompletedBy:   supervisor.ID,
			SecondStore:   true,
			ExpectedRequirement: types.RecurringMemberRequirement{
				ID:            "",
				MemberID:      m1.ID,
				RequirementID: rr1.ID,
				AssignedBy:    supervisor.ID,
				CompletionHistory: []types.RecurringMemberRequirementCompletion{
					{
						ID:             "",
						CompletionDate: constClock.Now(),
						CompletedBy:    supervisor.ID,
					},
					{
						ID:             "",
						CompletionDate: constClock2.Now(),
						CompletedBy:    supervisor.ID,
					},
				},
			},
			ExpectedError: nil,
		},
		{
			Name:                "Invalid Completion Date",
			MemberID:            m1.ID,
			RequirementID:       rr1.ID,
			CompletedBy:         supervisor.ID,
			SecondStore:         false,
			ExpectedRequirement: types.RecurringMemberRequirement{},
			ExpectedError:       backend.ErrInvalidRequirementCompletionDate,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SecondStore {
				err = mqStore2.CompleteRecurringMemberRequirement(tt.MemberID, tt.RequirementID, tt.CompletedBy)
			} else {
				err = mqStore.CompleteRecurringMemberRequirement(tt.MemberID, tt.RequirementID, tt.CompletedBy)
			}
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && err == nil {
				t.Error("Expected error but didn't get one")
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error was: %s\nGot: %s\n", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				recurringRequirement, err := mqStore.GetRecurringMemberRequirement(tt.MemberID, tt.RequirementID)
				if err != nil {
					t.Fatalf("Error getting recurring member requirement for test: %s", err.Error())
				}
				// We don't quite care about the ID of the requirement here
				if tt.ExpectedRequirement.RequirementID != recurringRequirement.RequirementID {
					t.Errorf("Expected RequirementID: %s\nGot: %s\n", tt.ExpectedRequirement.RequirementID, recurringRequirement.RequirementID)
				}
				if tt.ExpectedRequirement.MemberID != recurringRequirement.MemberID {
					t.Errorf("Expected MemberID: %s\nGot: %s\n", tt.ExpectedRequirement.MemberID, recurringRequirement.MemberID)
				}
				if tt.ExpectedRequirement.AssignedBy != recurringRequirement.AssignedBy {
					t.Errorf("Expected AssignedBy: %s\nGot: %s\n", tt.ExpectedRequirement.AssignedBy, recurringRequirement.AssignedBy)
				}
				if !slices.EqualFunc(tt.ExpectedRequirement.CompletionHistory, recurringRequirement.CompletionHistory, func(comp1, comp2 types.RecurringMemberRequirementCompletion) bool {
					return comp1.CompletedBy == comp2.CompletedBy &&
						comp1.CompletionDate == comp2.CompletionDate
				}) {
					t.Errorf("Expected CompletionHistory: %+v\nGot: %+v\n", tt.ExpectedRequirement.CompletionHistory, recurringRequirement.CompletionHistory)
				}
			}
		})
	}
}
