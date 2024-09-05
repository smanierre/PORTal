package sqlite

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func (p Provider) AssignMemberQualification(memberID, qualificationID string) error {
	_, err := p.Db.Exec(addMemberQualificationQuery, memberID, qualificationID)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: member_qualification.member_id, member_qualification.qualification_id") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Member already assigned qualification")
		return fmt.Errorf("%w: member_id=%s qualification_id=%s", backend.ErrQualificationAlreadyAssigned, memberID, qualificationID)
	} else if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		if _, err = p.GetMember(memberID, backend.ById); err != nil {
			return err
		}
		if _, err = p.GetQualification(qualificationID); err != nil {
			return err
		}
	} else if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding qualification to member", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p Provider) GetMemberQualification(memberID, qualificationID string) (types.Qualification, error) {
	// Verify that member has qualification assigned
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking to see if member has qualification assigned")
	row := p.Db.QueryRow(checkMemberQualificationQuery, memberID, qualificationID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error checking if member has qualification assigned", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	if count == 0 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Member with given qualification not found")
		return types.Qualification{}, backend.ErrMemberQualificationNotFound
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "User has qualification assigned")
	return p.GetQualification(qualificationID)
}

func (p Provider) GetMemberQualifications(memberID string) ([]types.Qualification, error) {
	rows, err := p.Db.Query(getMemberQualificationIDsQuery, memberID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting qualification IDs for member", slog.String("error", err.Error()))
		return nil, err
	}
	var ids []string
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning id into string", slog.String("error", err.Error()))
			return nil, err
		}
		ids = append(ids, id)
	}
	quals := make([]types.Qualification, 0, len(ids))
	var qual types.Qualification
	for _, id := range ids {
		qual, err = p.GetQualification(id)
		if err != nil {
			return nil, err
		}
		quals = append(quals, qual)
	}
	return quals, nil
}

func (p Provider) RemoveMemberQualification(memberId, qualificationId string) error {
	res, err := p.Db.Exec(removeMemberQualificationQuery, memberId, qualificationId)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error removing qualification from member", slog.String("error", err.Error()))
		return err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Could not find member qualification to remove")
		return fmt.Errorf("%w: member_id: %s, qualification_id: %s", backend.ErrMemberQualificationNotFound, memberId, qualificationId)
	}
	return nil
}
