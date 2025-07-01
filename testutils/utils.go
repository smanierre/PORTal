package testutils

import (
	"PORTal/types"
	"fmt"
	"math/rand/v2"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func RandomString() string {
	chars := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	s := strings.Builder{}
	for range (rand.Int() % 30) + 10 {
		s.WriteString(chars[rand.Int()%(len(chars)-1)])
	}
	return s.String()
}

func RandomMember(admin bool) types.Member {
	return types.Member{
		ID:           "",
		FirstName:    RandomString(),
		LastName:     RandomString(),
		Username:     RandomString(),
		Grade:        types.E4,
		SupervisorID: "",
		Admin:        admin,
		Password:     RandomString(),
		Hash:         "",
	}
}

func RandomQualification() types.Qualification {
	return types.Qualification{
		ID:                    "",
		Name:                  RandomString(),
		InitialRequirements:   []types.Requirement{},
		RecurringRequirements: []types.Requirement{},
		Notes:                 RandomString(),
	}
}

func RandomInitialRequirement(canBeQualType bool, id string) types.Requirement {
	var requirementTypes []types.RequirementType
	if canBeQualType && id != "" {
		requirementTypes = []types.RequirementType{types.GradeType, types.WbtType, types.QualificationType}
	} else {
		requirementTypes = []types.RequirementType{types.GradeType, types.WbtType}
	}

	var grades = []types.Grade{types.E1, types.E2, types.E3, types.E4, types.E5, types.E6, types.E7, types.E8, types.E9}
	n := rand.IntN(len(requirementTypes))
	g := rand.IntN(len(grades))
	switch requirementTypes[n] {
	case types.GradeType:
		return types.Requirement{
			Name:            RandomString(),
			Reference:       RandomString(),
			Initial:         true,
			QualificationID: "",
			Grade:           grades[g],
			Notes:           "",
			DaysValidFor:    0,
			Type:            types.GradeType,
		}
	case types.QualificationType:
		return types.Requirement{
			Name:            RandomString(),
			QualificationID: id,
			Initial:         true,
			Reference:       RandomString(),
			Type:            types.QualificationType,
		}
	case types.WbtType:
		return types.Requirement{
			Name:      RandomString(),
			Initial:   true,
			Reference: RandomString(),
			Notes:     RandomString(),
			Type:      types.WbtType,
		}
	}
	return types.Requirement{}
}

func RandomRecurringRequirement() types.Requirement {
	requirementTypes := []types.RequirementType{types.WbtType, types.ProficiencyType}
	n := rand.IntN(len(requirementTypes))
	switch requirementTypes[n] {
	case types.WbtType:
		return types.Requirement{
			Name:         RandomString(),
			Reference:    RandomString(),
			Notes:        RandomString(),
			DaysValidFor: rand.IntN(1000) + 1,
			Type:         types.WbtType,
		}
	case types.ProficiencyType:
		return types.Requirement{
			Name:         RandomString(),
			Reference:    RandomString(),
			Notes:        RandomString(),
			DaysValidFor: rand.IntN(1000) + 1,
			Type:         types.ProficiencyType,
		}
	default:
		return types.Requirement{}
	}
}

func VerifyUpdatedUserNoHash(updates, returned types.Member, t *testing.T) {
	updates.Hash = ""
	returned.Hash = ""
	VerifyUpdatedUser(updates, returned, t)
}

func VerifyUpdatedUser(updates, returned types.Member, t *testing.T) {
	// New password was provided, so hashes will be different. Blank out both
	if updates.Password != "" {
		updates.Password = ""
		returned.Hash = ""
	}
	if updates != returned {
		t.Errorf("Expected Member: %+v\nGot: %+v\n", updates, returned)
	}
}

func CompareQuals(got, wanted types.Qualification) bool {
	// Pull out slices to sort later and set passed values to nil so reflect.DeepEqual and determine the rest of the fields.
	gotInitialReqs := got.InitialRequirements
	wantInitialReqs := wanted.InitialRequirements
	gotRecurringReqs := got.RecurringRequirements
	wantRecurringReqs := wanted.RecurringRequirements
	got.InitialRequirements = nil
	wanted.InitialRequirements = nil
	got.RecurringRequirements = nil
	wanted.RecurringRequirements = nil
	sort.Slice(gotInitialReqs, func(i, j int) bool {
		return gotInitialReqs[i].ID < gotInitialReqs[j].ID
	})
	sort.Slice(wantInitialReqs, func(i, j int) bool {
		return wantInitialReqs[i].ID < wantInitialReqs[j].ID
	})
	sort.Slice(gotRecurringReqs, func(i, j int) bool {
		return gotRecurringReqs[i].ID < gotRecurringReqs[j].ID
	})
	sort.Slice(wantRecurringReqs, func(i, j int) bool {
		return wantRecurringReqs[i].ID < wantRecurringReqs[j].ID
	})

	return reflect.DeepEqual(got, wanted) && reflect.DeepEqual(gotInitialReqs, wantInitialReqs) && reflect.DeepEqual(gotRecurringReqs, wantRecurringReqs)
}

func CompareRequirements(r1, r2 types.Requirement) bool {
	if r1.ID != r2.ID {
		return false
	}
	if r1.Name != r2.Name {
		return false
	}
	if r1.Notes != r2.Notes {
		return false
	}
	if r1.DaysValidFor != r2.DaysValidFor {
		return false
	}
	if r1.Reference != r2.Reference {
		return false
	}
	return true
}

func GetDbString() string {
	return fmt.Sprintf("%s.db", uuid.NewString())
}
