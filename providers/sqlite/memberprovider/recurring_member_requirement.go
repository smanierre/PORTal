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

func (m MemberProvider) AssignRecurringMemberRequirement(id, memberID string, requirementID string, assignedBy string) error {
	_, err := m.db.Exec(addRecurringMemberRequirementsQuery, id, memberID, requirementID, assignedBy)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error assigning recurring requirement to member in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (m MemberProvider) GetRecurringMemberRequirement(id string) (types.RecurringMemberRequirement, error) {
	row := m.db.QueryRow(getRecurringMemberRequirementQuery, id)
	var req types.RecurringMemberRequirement
	var completedBy sql.NullString
	var completedAt sql.NullTime
	err := row.Scan(&req.ID, &req.MemberID, &req.RequirementID, &req.AssignedBy, &completedBy, &completedAt)
	if errors.Is(err, sql.ErrNoRows) {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Unable to find recurring member requirement", slog.String("error", err.Error()))
		return types.RecurringMemberRequirement{}, fmt.Errorf("%w, %w", backend.ErrRecurringMemberRequirementNotFound, err)
	}
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting recurring member requirement from database", slog.String("error", err.Error()))
		return types.RecurringMemberRequirement{}, err
	}
	completionHistory, err := m.getRecurringMemberRequirementCompletions(req.ID)
	if err != nil {
		return types.RecurringMemberRequirement{}, err
	}
	req.CompletionHistory = completionHistory
	if completedBy.Valid {
		req.CompletedBy = completedBy.String
	}
	if completedAt.Valid {
		req.CompletedDate = completedAt.Time
	}
	return req, nil
}

func (m MemberProvider) getRecurringMemberRequirementCompletions(memberRequirementID string) ([]types.RecurringMemberRequirementCompletion, error) {
	historyRows, err := m.db.Query(getRecurringMemberRequirementCompletionQuery, memberRequirementID)
	defer historyRows.Close()
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting requirement completion history from database", slog.String("error", err.Error()))
		return nil, err
	}
	var completionHistory []types.RecurringMemberRequirementCompletion
	for historyRows.Next() {
		var history types.RecurringMemberRequirementCompletion
		err = historyRows.Scan(&history.ID, &history.CompletionDate, &history.CompletedBy)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning RecurringMemberRequirementCompletion into struct", slog.String("error", err.Error()))
			return nil, err
		}
		completionHistory = append(completionHistory, history)
	}
	return completionHistory, nil
}

func (m MemberProvider) GetRecurringMemberRequirements(memberID string) ([]types.RecurringMemberRequirement, error) {
	rows, err := m.db.Query(getRecurringMemberRequirementsQuery, memberID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting recurring member requirements for qualification", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var memberRequirements []types.RecurringMemberRequirement
	for rows.Next() {
		var memberRequirement types.RecurringMemberRequirement
		var completedBy sql.NullString
		var completedAt sql.NullTime
		err = rows.Scan(&memberRequirement.ID, &memberRequirement.MemberID, &memberRequirement.RequirementID, &memberRequirement.AssignedBy, &completedBy, &completedAt)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning recurring member requirement to struct", slog.String("error", err.Error()))
			return nil, err
		}
		completionHistory, err := m.getRecurringMemberRequirementCompletions(memberID, memberRequirement.ID)
		if err != nil {
			return nil, err
		}
		memberRequirement.CompletionHistory = completionHistory
		if completedBy.Valid {
			memberRequirement.CompletedBy = completedBy.String
		}
		if completedAt.Valid {
			memberRequirement.CompletedDate = completedAt.Time
		}
		memberRequirements = append(memberRequirements, memberRequirement)
	}
	return memberRequirements, nil
}

func (m MemberProvider) GetRecurringMemberRequirementsForQualification(memberID, qualificationID string) ([]types.RecurringMemberRequirement, error) {
	rows, err := m.db.Query(getRecurringMemberRequirementsForQualificationQuery, memberID, qualificationID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting recurring member requirements for qualification", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var memberRequirements []types.RecurringMemberRequirement
	for rows.Next() {
		var memberRequirement types.RecurringMemberRequirement
		var completedBy sql.NullString
		var completedAt sql.NullTime
		err = rows.Scan(&memberRequirement.ID, &memberRequirement.MemberID, &memberRequirement.RequirementID, &memberRequirement.AssignedBy, &completedBy, &completedAt)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning recurring member requirement to struct", slog.String("error", err.Error()))
			return nil, err
		}
		completionHistory, err := m.getRecurringMemberRequirementCompletions(memberRequirement.ID)
		if err != nil {
			return nil, err
		}
		if completedBy.Valid {
			memberRequirement.CompletedBy = completedBy.String
		}
		if completedAt.Valid {
			memberRequirement.CompletedDate = completedAt.Time
		}
		memberRequirement.CompletionHistory = completionHistory
		memberRequirements = append(memberRequirements, memberRequirement)
	}
	return memberRequirements, nil
}

func (m MemberProvider) CompleteRecurringMemberRequirement(id, recurringRequirementID, completedBy string, completedDate time.Time) error {
	//update this to add a completion history before updating the current one
	r, err := m.GetRecurringMemberRequirement(id)
	if err != nil {
		return err
	}
	
	_, err := m.db.Exec(completeRecurringMemberRequirementQuery, id, recurringRequirementID, completedDate, completedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return backend.ErrRecurringMemberRequirementNotFound
	} else if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error completing recurring member requirement", slog.String("error", err.Error()))
		return err
	}
	return nil
}
