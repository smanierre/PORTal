package memberprovider

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func (m MemberProvider) AssignMemberQualification(memberID, qualificationID string, dateAssigned time.Time, assignedBy string) error {
	_, err := m.db.Exec(addMemberQualificationQuery, memberID, qualificationID, dateAssigned, assignedBy)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: member_qualification.member_id, member_qualification.qualification_id") {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Member already assigned qualification")
		return fmt.Errorf("%w: member_id=%s qualification_id=%s", backend.ErrQualificationAlreadyAssigned, memberID, qualificationID)
	} else if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		if _, err = m.GetMember(memberID, backend.ById); err != nil {
			return fmt.Errorf("%w: member_id=%s", backend.ErrMemberNotFound, memberID)
		}
		if _, err = m.GetMember(assignedBy, backend.ById); err != nil {
			return fmt.Errorf("%w: assigned_by=%s", backend.ErrMemberNotFound, assignedBy)
		}
		return fmt.Errorf("%w: qualification_id=%s", backend.ErrQualificationNotFound, qualificationID)
	} else if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding qualification to member", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (m MemberProvider) GetMemberQualification(memberID, qualificationID string) (types.MemberQualification, error) {
	var mq types.MemberQualification
	row := m.db.QueryRow(getMemberQualificationQuery, memberID, qualificationID)
	err := row.Scan(&mq.MemberID, &mq.QualificationID, &mq.DateAssigned, &mq.AssignedByID)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return types.MemberQualification{}, backend.ErrMemberQualificationNotFound
	}
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error getting member qualification from database", slog.String("error", err.Error()))
		return types.MemberQualification{}, err
	}
	return mq, nil
}

func (m MemberProvider) GetQualificationsForMember(memberID string) ([]types.MemberQualification, error) {
	var memberQuals []types.MemberQualification
	rows, err := m.db.Query(getQualificationsForMemberQuery, memberID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting qualification IDs for member", slog.String("error", err.Error()))
		return nil, err
	}
	for rows.Next() {
		var mq types.MemberQualification
		err = rows.Scan(&mq.MemberID, &mq.QualificationID, &mq.DateAssigned, &mq.AssignedByID)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning qualification into struct", slog.String("error", err.Error()))
			return nil, err
		}
		memberQuals = append(memberQuals, mq)
	}
	return memberQuals, nil
}

func (m MemberProvider) GetAllMemberQualifications() ([]types.MemberQualification, error) {
	rows, err := m.db.Query(getMemberQualificationsQuery)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting qualification IDs for member", slog.String("error", err.Error()))
		return nil, err
	}
	var memberQuals []types.MemberQualification
	for rows.Next() {
		var mq types.MemberQualification
		err = rows.Scan(&mq.MemberID, &mq.QualificationID, &mq.DateAssigned, &mq.AssignedByID)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning qualification into struct", slog.String("error", err.Error()))
			continue
		}
		memberQuals = append(memberQuals, mq)
	}
	return memberQuals, nil
}

func (m MemberProvider) RemoveMemberQualification(memberId, qualificationId string) error {
	res, err := m.db.Exec(removeMemberQualificationQuery, memberId, qualificationId)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error removing qualification from member", slog.String("error", err.Error()))
		return err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Could not find member qualification to remove")
		return fmt.Errorf("%w: member_id: %s, qualification_id: %s", backend.ErrMemberQualificationNotFound, memberId, qualificationId)
	}
	return nil
}
