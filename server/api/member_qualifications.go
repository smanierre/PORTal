package api

import (
	"PORTal/backend"
	"PORTal/types"
	"encoding/json"
	"errors"
	"net/http"
)

func memberQualificationToApiMemberQualification(mq types.MemberQualification) MemberQualification {
	return MemberQualification{
		AssignedBy:      mq.AssignedByID,
		DateAssigned:    mq.DateAssigned,
		MemberId:        mq.MemberID,
		QualificationId: mq.QualificationID,
	}
}

func memberQualificationsToApiMemberQualifications(mqs []types.MemberQualification) []MemberQualification {
	m := []MemberQualification{}
	for _, mq := range mqs {
		m = append(m, memberQualificationToApiMemberQualification(mq))
	}
	return m
}

func (a api) GetApiMemberQualifications(w http.ResponseWriter, r *http.Request, p GetApiMemberQualificationsParams) {
	memberRequirements, err := a.memberQualificationStore.GetQualificationsForMember(p.MemberId)
	if errors.Is(err, backend.ErrMemberNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(memberQualificationsToApiMemberQualifications(memberRequirements))
}

func (a api) PostApiMemberQualification(w http.ResponseWriter, r *http.Request, p PostApiMemberQualificationParams) {
	err := a.memberQualificationStore.AssignMemberQualification(p.MemberId, p.QualificationId, p.AssignedBy)
	if errors.Is(err, backend.ErrMemberNotFound) || errors.Is(err, backend.ErrQualificationNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mq, err := a.memberQualificationStore.GetMemberQualification(p.MemberId, p.QualificationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(memberQualificationToApiMemberQualification(mq))
}
