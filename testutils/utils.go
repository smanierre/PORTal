package testutils

import (
	"PORTal/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand/v2"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func RandomString() string {
	chars := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	s := strings.Builder{}
	for range (rand.Int() % 30) + 10 {
		s.WriteString(chars[rand.Int()%(len(chars)-1)])
	}
	return s.String()
}

func RandomMember() types.Member {
	return types.Member{
		ApiMember: types.ApiMember{
			ID:           "",
			FirstName:    RandomString(),
			LastName:     RandomString(),
			Username:     RandomString(),
			Rank:         types.E4,
			SupervisorID: "",
		},
		Password: RandomString(),
		Hash:     "",
	}
}

func RandomQualification() types.Qualification {
	expires := rand.IntN(100) > 50
	days := 0
	if expires {
		days = rand.IntN(730)
	}
	return types.Qualification{
		ID:                    "",
		Name:                  RandomString(),
		InitialRequirements:   nil,
		RecurringRequirements: nil,
		Notes:                 RandomString(),
		Expires:               expires,
		ExpirationDays:        days,
	}
}

func RandomRequirement(r types.Reference) types.Requirement {
	return types.Requirement{
		ID:           uuid.NewString(),
		Name:         RandomString(),
		Description:  RandomString(),
		Notes:        RandomString(),
		DaysValidFor: rand.IntN(1000) + 1,
		Reference:    r,
	}
}

func RandomReference() types.Reference {
	return types.Reference{
		Name:      RandomString(),
		Volume:    rand.IntN(10),
		Paragraph: RandomString(),
	}
}

func VerifyUpdatedUser(original, updates, returned types.Member, t *testing.T) {
	if updates.FirstName != "" {
		original.FirstName = updates.FirstName
	}
	if updates.LastName != "" {
		original.LastName = updates.LastName
	}
	if updates.Rank != "" {
		original.Rank = updates.Rank
	}
	if updates.Username != "" {
		original.Username = updates.Username
	}
	if updates.SupervisorID != "" {
		original.SupervisorID = updates.SupervisorID
	}
	if updates.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(returned.Hash), []byte(updates.Password)); err != nil {
			t.Errorf("Updated password does not match returned hash")
		}
	}
	if !reflect.DeepEqual(original.ApiMember, returned.ApiMember) {
		t.Errorf("Expected updated member: %+v\nGot: %+v", original, returned)
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
	if r1.Description != r2.Description {
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
