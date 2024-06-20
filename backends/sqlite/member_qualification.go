package sqlite

import (
	"PORTal/types"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func (b *Backend) AddMemberQualification(qualID, memberID string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding qualification to member", slog.String("qualification_id", qualID), slog.String("member_id", memberID))

	// Add qualification to member with qual inactive and active time set to never.
	_, err := b.Db.Exec(addMemberQualificationQuery, memberID, qualID, false, types.Never)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: member_qualification.member_id, member_qualification.qualification_id") {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Member already assigned qualification")
		return fmt.Errorf("%w: member_id=%s qualification_id=%s", types.ErrQualificationAlreadyAssigned, memberID, qualID)
	} else if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		if _, err = b.GetMember(memberID); err != nil {
			return err
		}
		if _, err = b.GetQualification(qualID); err != nil {
			return err
		}
	} else if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding qualification to member", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (b *Backend) GetMemberQualifications(memberID string) ([]types.MemberQualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualifications for member", slog.String("member_id", memberID))
	if _, err := b.GetMember(memberID); errors.Is(err, types.ErrMemberNotFound) {
		return nil, err
	}
	rows, err := b.Db.Query(getMemberQualificationsQuery, memberID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting qualifications for member", slog.String("error", err.Error()))
		return nil, err
	}
	var quals []types.MemberQualification
	var qualID string
	for rows.Next() {
		mq := types.MemberQualification{}
		err := rows.Scan(&mq.MemberID, &qualID, &mq.Active, &mq.ActiveDate)
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning MemberQualification into struct", slog.String("error", err.Error()))
			continue
		}
		mq.Qualification, err = b.GetQualification(qualID)
		if err != nil {
			continue
		}
		mq.ActiveDate = mq.ActiveDate.In(time.Local)
		quals = append(quals, mq)
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d qualifications for member", len(quals)))
	return quals, nil
}

func (b *Backend) UpdateMemberQualification(mq types.MemberQualification) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating qualification for member", slog.Any("member_qualification", mq))
	res, err := b.Db.Exec(updateMemberQualificationQuery, mq.Active, mq.ActiveDate, mq.MemberID, mq.Qualification.ID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating member qualification", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "0 rows updated, determining if it was member or qualification related")
		if _, err := b.GetMember(mq.MemberID); errors.Is(err, types.ErrMemberNotFound) {
			return err
		}
		if _, err := b.GetQualification(mq.Qualification.ID); errors.Is(err, types.ErrQualificationNotFound) {
			return err
		}
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Unable to determine why no rows were updated")
		return err
	}
	return nil
}

func (b *Backend) DeleteMemberQualification(qualID, memberID string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Removing qualification from member", slog.String("member_id", memberID), slog.String("qualification_id", qualID))
	res, err := b.Db.Exec(deleteMemberQualificationQuery, memberID, qualID)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Did not find any qualifications that matched member and qualification combo")
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking to see if member exists")
		if _, err := b.GetMember(memberID); errors.Is(err, types.ErrMemberNotFound) {
			return err
		}
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking to see if qualification exists")
		if _, err := b.GetQualification(qualID); errors.Is(err, types.ErrQualificationNotFound) {
			return err
		}
		return types.ErrMemberQualificationNotFound
	}
	return nil
}
