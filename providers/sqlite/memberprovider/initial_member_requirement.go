package memberprovider

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

func (m MemberProvider) AssignInitialMemberRequirement(memberID, requirementID, assignedBy string) error {
	_, err := m.db.Exec(addInitialMemberRequirementQuery, memberID, requirementID, nil, assignedBy, nil)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error assigning initial requirement to member", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (m MemberProvider) GetInitialMemberRequirement(memberID, requirementID string) (types.InitialMemberRequirement, error) {
	row := m.db.QueryRow(getInitialMemberRequirementQuery, memberID, requirementID)
	var i types.InitialMemberRequirement
	var verifiedByStr sql.NullString
	var completedDate sql.NullTime
	err := row.Scan(&i.MemberID, &i.RequirementID, &completedDate, &i.AssignedBy, &verifiedByStr)
	if errors.Is(err, sql.ErrNoRows) {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Couldn't find initial member requirement", slog.String("error", err.Error()))
		return types.InitialMemberRequirement{}, fmt.Errorf("%w: %s", backend.ErrInitialMemberRequirementNotFound, err.Error())
	}
	if err != nil {
		return types.InitialMemberRequirement{}, err
	}
	if verifiedByStr.Valid {
		i.CompletedBy = verifiedByStr.String
	}
	if completedDate.Valid {
		i.CompletedDate = completedDate.Time
	}
	return i, nil
}

func (m MemberProvider) GetInitialMemberRequirements(memberID string) ([]types.InitialMemberRequirement, error) {
	rows, err := m.db.Query(getInitialMemberRequirementsForQualificationQuery, memberID)
	if err != nil {
		return nil, err
	}
	var requirements []types.InitialMemberRequirement
	for rows.Next() {
		var i types.InitialMemberRequirement
		var verifiedByStr sql.NullString
		var completedDate sql.NullTime
		err := rows.Scan(&i.MemberID, &i.RequirementID, &completedDate, &i.AssignedBy, &verifiedByStr)
		if err != nil {
			return nil, err
		}
		if verifiedByStr.Valid {
			i.CompletedBy = verifiedByStr.String
		}
		if completedDate.Valid {
			i.CompletedDate = completedDate.Time
		}
		requirements = append(requirements, i)
	}
	return requirements, nil
}

func (m MemberProvider) GetInitialMemberRequirementsForQualification(memberID, qualificationID string) ([]types.InitialMemberRequirement, error) {
	rows, err := m.db.Query(getInitialMemberRequirementsForQualificationQuery, memberID, qualificationID)
	if err != nil {
		return nil, err
	}
	var requirements []types.InitialMemberRequirement
	for rows.Next() {
		var i types.InitialMemberRequirement
		var verifiedByStr sql.NullString
		var completedDate sql.NullTime
		err := rows.Scan(&i.MemberID, &i.RequirementID, &completedDate, &i.AssignedBy, &verifiedByStr)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning initial member requirement into struct", slog.String("error", err.Error()))
			return nil, err
		}
		if verifiedByStr.Valid {
			i.CompletedBy = verifiedByStr.String
		}
		if completedDate.Valid {
			i.CompletedDate = completedDate.Time
		}
		requirements = append(requirements, i)
	}
	return requirements, nil
}

func (m MemberProvider) CompleteInitialMemberRequirement(completedDate time.Time, completedBy, memberID, requirementID string) error {
	res, err := m.db.Exec(completeInitialMemberRequirementQuery, completedDate, completedBy, memberID, requirementID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error completing initial member requirement", slog.String("error", err.Error()))
		return err
	}
	if updated, _ := res.RowsAffected(); updated != 1 {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Unable to find row to update", slog.String("MemberID", memberID), slog.String("RequirementID", requirementID))
		return backend.ErrInitialMemberRequirementNotFound
	}
	return nil
}
