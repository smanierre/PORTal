package memberprovider

import (
	"PORTal/backend"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func (m MemberProvider) AssignMemberQualification(memberID, qualificationID string) error {
	_, err := m.db.Exec(addMemberQualificationQuery, memberID, qualificationID)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: member_qualification.member_id, member_qualification.qualification_id") {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Member already assigned qualification")
		return fmt.Errorf("%w: member_id=%s qualification_id=%s", backend.ErrQualificationAlreadyAssigned, memberID, qualificationID)
	} else if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		if _, err = m.GetMember(memberID, backend.ById); err != nil {
			return err
		}
		return backend.ErrQualificationNotFound
	} else if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding qualification to member", slog.String("error", err.Error()))
		return err
	}
	return nil
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
