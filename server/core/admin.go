package core

import (
	"PORTal/backend"
	"PORTal/types"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AdminMembersData struct {
	Members              []types.Member `json:"members,omitempty"`
	SelectedMember       types.Member   `json:"selected_member,omitempty"`
	Ranks                types.RankMap  `json:"ranks,omitempty"`
	PotentialSupervisors []types.Member `json:"potential_supervisors,omitempty"`
	DisabledMembers      bool           `json:"disabled_members,omitempty"`
	MemberQuery          string         `json:"member_query,omitempty"`
	NewMember            bool           `json:"new_member,omitempty"`
}

func (c Core) AdminMemberPage(r *http.Request) (AdminMembersData, error) {
	var members []types.Member
	var err error

	disabledMembers := r.URL.Query().Get("disabled") == "true"
	memberQuery := r.URL.Query().Get("member_query")
	selectedMemberID := r.URL.Query().Get("selected")

	if disabledMembers {
		members, err = c.memberStore.GetDisabledMembers()
		if err != nil {
			return AdminMembersData{}, err
		}
	} else {
		members, err = c.memberStore.GetAllMembers()
		if err != nil {
			return AdminMembersData{}, nil
		}
	}
	if memberQuery != "" {
		members = types.FilterMembers(members, func(m types.Member) bool {
			return strings.Contains(
				strings.ToLower(fmt.Sprintf("%s %s", m.FirstName, m.LastName)),
				strings.ToLower(memberQuery),
			)
		})
	}
	var selectedMember types.Member
	var potentialSupervisors []types.Member
	if selectedMemberID != "" {
		if disabledMembers {
			selectedMember, err = c.memberStore.GetDisabledMember(selectedMemberID)
		} else {
			selectedMember, err = c.memberStore.GetMember(selectedMemberID)
		}
		if err != nil && !errors.Is(err, backend.ErrMemberNotFound) {
			return AdminMembersData{}, err
		}
		// Disabled members don't have supervisors
		if !disabledMembers {
			potentialSupervisors, err = c.memberStore.GetPotentialSupervisors(selectedMember.ID, selectedMember.Grade)
			if err != nil {
				return AdminMembersData{}, err
			}
		}
	}

	return AdminMembersData{
		Members:              members,
		SelectedMember:       selectedMember,
		Ranks:                c.Ranks,
		PotentialSupervisors: potentialSupervisors,
		DisabledMembers:      disabledMembers,
		MemberQuery:          memberQuery,
	}, nil
}

type AdminQualificationData struct {
	Qualifications        []types.Qualification
	SelectedQualification types.Qualification
	NewQualification      bool
	QualificationQuery    string
}

func (c Core) AdminQualificationPage(r *http.Request) (AdminQualificationData, error) {
	var err error

	qualificationQuery := r.URL.Query().Get("qualification_query")
	selectedQualificationID := r.URL.Query().Get("selected")

	qualifications, err := c.qualificationStore.GetAllQualifications()
	if err != nil {
		return AdminQualificationData{}, err
	}
	if qualificationQuery != "" {
		qualifications = types.FilterQualifications(qualifications, func(q types.Qualification) bool {
			return strings.Contains(strings.ToLower(q.Name), strings.ToLower(qualificationQuery))
		})
	}
	var selectedQualification types.Qualification
	if selectedQualificationID != "" {
		selectedQualification, err = c.qualificationStore.GetQualification(selectedQualificationID)
		if err != nil {
			return AdminQualificationData{}, err
		}
	}
	return AdminQualificationData{
		Qualifications:        qualifications,
		SelectedQualification: selectedQualification,
		QualificationQuery:    qualificationQuery,
	}, nil
}
