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
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

type constantTimeClock struct {
	t string
}

func (c constantTimeClock) Now() time.Time {
	if c.t == "" {
		c.t = "2000-01-01 01:00:00"
	}
	t, err := time.Parse(time.DateTime, c.t)
	if err != nil {
		panic(fmt.Sprintf("Error parsing time: %s", err.Error()))
	}
	return t
}

func TestAssignMemberQualification(t *testing.T) {
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
	supvervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	m1, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}

	tc := []struct {
		Name            string
		ExpectedError   error
		MemberId        string
		QualificationId string
		AssignedByID    string
	}{
		{
			Name:            "Member Not found",
			ExpectedError:   backend.ErrMemberNotFound,
			MemberId:        uuid.NewString(),
			QualificationId: q1.ID,
		},
		{
			Name:            "Assigned by not found",
			ExpectedError:   backend.ErrMemberNotFound,
			MemberId:        m1.ID,
			QualificationId: q1.ID,
			AssignedByID:    uuid.NewString(),
		},
		{
			Name:            "Qualification not found",
			ExpectedError:   backend.ErrQualificationNotFound,
			MemberId:        m1.ID,
			QualificationId: uuid.NewString(),
			AssignedByID:    supvervisor.ID,
		},
		{
			Name:            "Swapped arguments",
			ExpectedError:   backend.ErrMemberNotFound,
			MemberId:        q1.ID,
			QualificationId: m1.ID,
		},
		{
			Name:            "Successful assign",
			ExpectedError:   nil,
			MemberId:        m1.ID,
			QualificationId: q1.ID,
			AssignedByID:    supvervisor.ID,
		},
		{
			Name:            "Duplicate Assignment",
			ExpectedError:   backend.ErrQualificationAlreadyAssigned,
			MemberId:        m1.ID,
			QualificationId: q1.ID,
			AssignedByID:    supvervisor.ID,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err = mqStore.AssignMemberQualification(tt.MemberId, tt.QualificationId, tt.AssignedByID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError, err)
			}
		})
	}
}

func TestAssignMemberQualificationRequirements(t *testing.T) {
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

	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	initialReq1, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, initialReq1.ID, true)
	if err != nil {
		t.Fatalf("Error assigning requirement: %s", err.Error())
	}
	initialReq2, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, initialReq2.ID, true)
	if err != nil {
		t.Fatalf("Error assigning requirement: %s", err.Error())
	}
	recurringReq1, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, recurringReq1.ID, false)
	if err != nil {
		t.Fatalf("Error assigning requirement: %s", err.Error())
	}
	recurringReq2, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q1.ID, recurringReq2.ID, false)
	if err != nil {
		t.Fatalf("Error assigning requirement: %s", err.Error())
	}
	m1, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification: %s", err.Error())
	}
	// Create and assign another qualification to ensure only the correct member requirements are pulled
	q2, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	ir1, err := qualificationStore.AddRequirement(testutils.RandomInitialRequirement(false, ""))
	if err != nil {
		t.Fatalf("Error adding requirement: %s", err.Error())
	}
	rr1, err := qualificationStore.AddRequirement(testutils.RandomRecurringRequirement())
	if err != nil {
		t.Fatalf("Error adding requirement: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q2.ID, ir1.ID, true)
	if err != nil {
		t.Fatalf("Error assigning requirement: %s", err.Error())
	}
	err = qualificationStore.AssignRequirementToQualification(q2.ID, rr1.ID, false)
	if err != nil {
		t.Fatalf("Error assigning requirement: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q2.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification: %s", err.Error())
	}
	// Add another user that has both quals assigned as well to ensure filtering on user
	m2, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m2.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m2.ID, q2.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification: %s", err.Error())
	}
	initialMemberRequirements, err := mqStore.GetInitialMemberRequirementsForQualification(m1.ID, q1.ID)
	if err != nil {
		t.Fatalf("Error getting initial member requirements: %s", err.Error())
	}

	recurringMemberRequirements, err := mqStore.GetRecurringMemberRequirementsForQualification(m1.ID, q1.ID)
	if err != nil {
		t.Fatalf("Error getting recurring member requirements: %s", err.Error())
	}

	expectedInitialRequirements := []types.InitialMemberRequirement{
		types.InitialMemberRequirement{
			MemberID:      m1.ID,
			RequirementID: initialReq1.ID,
			CompletedDate: types.Never,
			AssignedBy:    supervisor.ID,
		},
		types.InitialMemberRequirement{
			MemberID:      m1.ID,
			RequirementID: initialReq2.ID,
			CompletedDate: types.Never,
			AssignedBy:    supervisor.ID,
		},
	}
	slices.SortFunc(initialMemberRequirements, func(a types.InitialMemberRequirement, b types.InitialMemberRequirement) int {
		if a.RequirementID > b.RequirementID {
			return 1
		}
		return -1
	})
	slices.SortFunc(expectedInitialRequirements, func(a types.InitialMemberRequirement, b types.InitialMemberRequirement) int {
		if a.RequirementID > b.RequirementID {
			return 1
		}
		return -1
	})
	if !slices.Equal(expectedInitialRequirements, initialMemberRequirements) {
		t.Errorf("Expected initial requirements: %+v\nGot: %+v\n", expectedInitialRequirements, initialMemberRequirements)
	}
	expectedRecurringRequirements := []types.RecurringMemberRequirement{
		types.RecurringMemberRequirement{
			ID:                "",
			MemberID:          m1.ID,
			RequirementID:     recurringReq1.ID,
			AssignedBy:        supervisor.ID,
			CompletionHistory: nil,
		},
		types.RecurringMemberRequirement{
			ID:                "",
			MemberID:          m1.ID,
			RequirementID:     recurringReq2.ID,
			AssignedBy:        supervisor.ID,
			CompletionHistory: nil,
		},
	}
	slices.SortFunc(recurringMemberRequirements, func(a types.RecurringMemberRequirement, b types.RecurringMemberRequirement) int {
		if a.RequirementID > b.RequirementID {
			return 1
		}
		return -1
	})
	slices.SortFunc(expectedRecurringRequirements, func(a types.RecurringMemberRequirement, b types.RecurringMemberRequirement) int {
		if a.RequirementID > b.RequirementID {
			return 1
		}
		return -1
	})
	if !slices.EqualFunc(recurringMemberRequirements, expectedRecurringRequirements, func(r1 types.RecurringMemberRequirement, r2 types.RecurringMemberRequirement) bool {
		// We don't care about the id's here, everything else should match
		r1.ID = r2.ID
		return reflect.DeepEqual(r1, r2)
	}) {
		t.Errorf("Expected recurring requirements: %+v\nGot: %+v\n", expectedRecurringRequirements, recurringMemberRequirements)
	}
}

func TestDeleteMemberQualification(t *testing.T) {
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
		t.Fatalf("Error adding member: %s", err.Error())
	}
	m1, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification: %s", err.Error())
	}

	tc := []struct {
		Name            string
		MemberID        string
		QualificationID string
		ExpectedError   error
	}{
		{
			Name:            "Successful delete",
			MemberID:        m1.ID,
			QualificationID: q1.ID,
			ExpectedError:   nil,
		},
		{
			Name:            "Member not found",
			MemberID:        uuid.NewString(),
			QualificationID: q1.ID,
			ExpectedError:   backend.ErrMemberQualificationNotFound,
		},
		{
			Name:            "Qualification Not Found",
			MemberID:        m1.ID,
			QualificationID: uuid.NewString(),
			ExpectedError:   backend.ErrMemberQualificationNotFound,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			err = mqStore.RemoveMemberQualification(tt.MemberID, tt.QualificationID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error, got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError, err)
			}
			if tt.ExpectedError == nil {
				_, err = mqStore.GetMemberQualification(tt.MemberID, tt.QualificationID)
				if !errors.Is(err, backend.ErrMemberQualificationNotFound) {
					t.Errorf("Expected error: %s, got: %s", backend.ErrMemberQualificationNotFound, err)
				}
			}
		})
	}
}

func TestGetMemberQualification(t *testing.T) {
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
	c := constantTimeClock{}
	qualificationStore := qualificationstore.New(qualificationProvider, logger)
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, c)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	m1, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
	if err != nil {
		t.Fatalf("Error assigning member qualification: %s", err.Error())
	}

	tc := []struct {
		Name                        string
		MemberID                    string
		QualificationID             string
		ExpectedMemberQualification types.MemberQualification
		ExpectedError               error
		SetupFunc                   func(t *testing.T)
	}{
		{
			Name:            "Successful get",
			MemberID:        m1.ID,
			QualificationID: q1.ID,
			ExpectedMemberQualification: types.MemberQualification{
				MemberID:        m1.ID,
				QualificationID: q1.ID,
				DateAssigned:    c.Now(),
				AssignedByID:    supervisor.ID,
			},
		},
		{
			Name:            "Assigning Member Disabled",
			MemberID:        m1.ID,
			QualificationID: q1.ID,
			ExpectedMemberQualification: types.MemberQualification{
				MemberID:        m1.ID,
				QualificationID: q1.ID,
				DateAssigned:    c.Now(),
				AssignedByID:    supervisor.ID,
			},
			ExpectedError: nil,
			SetupFunc: func(t *testing.T) {
				err = memberStore.DisableMember(supervisor.ID)
				if err != nil {
					t.Fatalf("Error disabling member: %s", err.Error())
				}
			},
		},
		{
			Name:                        "Member not found",
			MemberID:                    uuid.NewString(),
			QualificationID:             q1.ID,
			ExpectedMemberQualification: types.MemberQualification{},
			ExpectedError:               backend.ErrMemberQualificationNotFound,
			SetupFunc:                   nil,
		},
		{
			Name:                        "Qualification not found",
			MemberID:                    m1.ID,
			QualificationID:             uuid.NewString(),
			ExpectedMemberQualification: types.MemberQualification{},
			ExpectedError:               backend.ErrMemberQualificationNotFound,
			SetupFunc:                   nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SetupFunc != nil {
				tt.SetupFunc(t)
			}
			memberQual, err := mqStore.GetMemberQualification(tt.MemberID, tt.QualificationID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s\n", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				if memberQual != tt.ExpectedMemberQualification {
					t.Errorf("Expected member qualification: %v\nGot: %v", tt.ExpectedMemberQualification, memberQual)
				}
			}
		})
	}
}

func TestGetAllMemberQualifications(t *testing.T) {
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
	c := constantTimeClock{}
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, c)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	m1, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	m2, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	q2, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	var addedMemberquals []types.MemberQualification

	type testCase struct {
		Name                         string
		ExpectedMemberQualifications []types.MemberQualification
		SetupFunc                    func(t *testing.T, tc *testCase)
	}
	tc := []testCase{
		{
			Name:                         "No member qualifications",
			ExpectedMemberQualifications: []types.MemberQualification{},
		},
		{
			Name:                         "One member, one qual",
			ExpectedMemberQualifications: []types.MemberQualification{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
				if err != nil {
					t.Fatalf("Error assigning member qualification: %s", err.Error())
				}
				addedMemberquals = append(addedMemberquals, types.MemberQualification{
					MemberID:        m1.ID,
					QualificationID: q1.ID,
					DateAssigned:    c.Now(),
					AssignedByID:    supervisor.ID,
				})
				tc.ExpectedMemberQualifications = addedMemberquals
			},
		},
		{
			Name:                         "One member, two quals",
			ExpectedMemberQualifications: []types.MemberQualification{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				err = mqStore.AssignMemberQualification(m1.ID, q2.ID, supervisor.ID)
				if err != nil {
					t.Fatalf("Error assigning member qualification: %s", err.Error())
				}
				addedMemberquals = append(addedMemberquals, types.MemberQualification{
					MemberID:        m1.ID,
					QualificationID: q2.ID,
					DateAssigned:    c.Now(),
					AssignedByID:    supervisor.ID,
				})
				tc.ExpectedMemberQualifications = addedMemberquals
			},
		},
		{
			Name:                         "Two members, three quals",
			ExpectedMemberQualifications: []types.MemberQualification{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				err = mqStore.AssignMemberQualification(m2.ID, q1.ID, supervisor.ID)
				if err != nil {
					t.Fatalf("Error assigning member qualification: %s", err.Error())
				}
				addedMemberquals = append(addedMemberquals, types.MemberQualification{
					MemberID:        m2.ID,
					QualificationID: q1.ID,
					DateAssigned:    c.Now(),
					AssignedByID:    supervisor.ID,
				})
				tc.ExpectedMemberQualifications = addedMemberquals
			},
		},
		{
			Name:                         "Two members, four quals",
			ExpectedMemberQualifications: []types.MemberQualification{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				err = mqStore.AssignMemberQualification(m2.ID, q2.ID, supervisor.ID)
				if err != nil {
					t.Fatalf("Error assigning member qualification: %s", err.Error())
				}
				addedMemberquals = append(addedMemberquals, types.MemberQualification{
					MemberID:        m2.ID,
					QualificationID: q2.ID,
					DateAssigned:    c.Now(),
					AssignedByID:    supervisor.ID,
				})
				tc.ExpectedMemberQualifications = addedMemberquals
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SetupFunc != nil {
				tt.SetupFunc(t, &tt)
			}
			memberQuals, err := mqStore.GetAllMemberQualifications()
			if err != nil {
				t.Errorf("Error getting member qualifications: %s", err.Error())
			}
			if !slices.Equal(memberQuals, tt.ExpectedMemberQualifications) {
				t.Errorf("Expected member qualifications: %v\nGot: %v", tt.ExpectedMemberQualifications, memberQuals)
			}
		})
	}
}

func TestGetQualificationsForMember(t *testing.T) {
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
	c := constantTimeClock{}
	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, logger, c)
	supervisor, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	m1, err := memberStore.AddMember(testutils.RandomMember(false))
	if err != nil {
		t.Fatalf("Error adding member: %s", err.Error())
	}
	q1, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	q2, err := qualificationStore.AddQualification(testutils.RandomQualification())
	if err != nil {
		t.Fatalf("Error adding qualification: %s", err.Error())
	}
	var addedMemberQuals []types.MemberQualification

	type testCase struct {
		Name                         string
		MemberID                     string
		ExpectedMemberQualifications []types.MemberQualification
		ExpectedError                error
		SetupFunc                    func(t *testing.T, tc *testCase)
	}

	tc := []testCase{
		{
			Name:                         "No qualifications",
			MemberID:                     m1.ID,
			ExpectedMemberQualifications: []types.MemberQualification{},
		},
		{
			Name:                         "One qualification",
			MemberID:                     m1.ID,
			ExpectedMemberQualifications: []types.MemberQualification{},
			SetupFunc: func(t *testing.T, tc *testCase) {
				err = mqStore.AssignMemberQualification(m1.ID, q1.ID, supervisor.ID)
				if err != nil {
					t.Fatalf("Error assigning member qualification: %s", err.Error())
				}
				addedMemberQuals = append(addedMemberQuals, types.MemberQualification{
					MemberID:        m1.ID,
					QualificationID: q1.ID,
					DateAssigned:    c.Now(),
					AssignedByID:    supervisor.ID,
				})
				tc.ExpectedMemberQualifications = addedMemberQuals
			},
		},
		{
			Name:                         "Two qualifications",
			MemberID:                     m1.ID,
			ExpectedMemberQualifications: []types.MemberQualification{},
			ExpectedError:                nil,
			SetupFunc: func(t *testing.T, tc *testCase) {
				err = mqStore.AssignMemberQualification(m1.ID, q2.ID, supervisor.ID)
				if err != nil {
					t.Fatalf("Error assigning member qualification: %s", err.Error())
				}
				addedMemberQuals = append(addedMemberQuals, types.MemberQualification{
					MemberID:        m1.ID,
					QualificationID: q2.ID,
					DateAssigned:    c.Now(),
					AssignedByID:    supervisor.ID,
				})
				tc.ExpectedMemberQualifications = addedMemberQuals
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.SetupFunc != nil {
				tt.SetupFunc(t, &tt)
			}
			quals, err := mqStore.GetQualificationsForMember(tt.MemberID)
			if tt.ExpectedError == nil && err != nil {
				t.Errorf("Expected no error but got: %s", err.Error())
			}
			if tt.ExpectedError != nil && !errors.Is(err, tt.ExpectedError) {
				t.Errorf("Expected error: %s\nGot: %s", tt.ExpectedError.Error(), err.Error())
			}
			if tt.ExpectedError == nil {
				// Sort to make sure the order matches up
				slices.SortFunc(quals, func(i types.MemberQualification, j types.MemberQualification) int {
					return strings.Compare(i.QualificationID, j.QualificationID)
				})
				slices.SortFunc(tt.ExpectedMemberQualifications, func(i types.MemberQualification, j types.MemberQualification) int {
					return strings.Compare(i.QualificationID, j.QualificationID)
				})
				if !slices.Equal(quals, tt.ExpectedMemberQualifications) {
					t.Errorf("Expected qualifications: %v\nGot: %v", tt.ExpectedMemberQualifications, quals)
				}
			}
		})
	}
}
