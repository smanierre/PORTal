package sqlite

import (
	"PORTal/types"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func (b *Backend) AddRequirement(r types.Requirement) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args...")
	if err := r.CheckForMissingArgs(); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Required arguments missing", slog.String("error", err.Error()))
		return err
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding requirement to database", slog.Any("requirement", r))
	_, err := b.Db.Exec(addRequirementQuery, r.ID, r.Name, r.Description, r.Notes, r.DaysValidFor)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting requirement into database: %s", slog.String("error", err.Error()))
	}
	return nil
}

func (b *Backend) GetRequirement(requirementID string) (types.Requirement, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting requirement from database", slog.String("requirement_id", requirementID))
	row := b.Db.QueryRow(getRequirementQuery, requirementID)
	r := types.Requirement{}
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.Notes, &r.DaysValidFor)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "No results found for requirement with given id")
		return types.Requirement{}, fmt.Errorf("%w: requirement_id=%s", types.ErrRequirementNotFound, requirementID)
	}
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning requirement into struct", slog.String("error", err.Error()))
		return types.Requirement{}, err
	}
	return r, nil
}

func (b *Backend) GetAllRequirements() ([]types.Requirement, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all requirements from database")
	rows, err := b.Db.Query(getAllRequirementsQuery)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all requirements from databse: %s", slog.String("error", err.Error()))
		return nil, err
	}
	var reqs []types.Requirement
	var r types.Requirement
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.Notes, &r.DaysValidFor)
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning requirement into struct", slog.String("error", err.Error()))
			continue
		}
		reqs = append(reqs, r)
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d requirements", len(reqs)))
	return reqs, nil
}

func (b *Backend) UpdateRequirement(r types.Requirement) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating requirement", slog.Any("new_requirement", r))
	res, err := b.Db.Exec(updateRequirementQuery, r.Name, r.Description, r.Notes, r.DaysValidFor, r.ID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating requirement in database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated but didn't get that")
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking to see if requirement exists")
		if _, err := b.GetRequirement(r.ID); errors.Is(err, types.ErrRequirementNotFound) {
			return err
		}
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Can't determine why row wasn't updated")
		return errors.New("couldn't update requirement")
	}
	return nil
}

func (b *Backend) DeleteRequirement(requirementID string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting requirement from database", slog.String("requirement_id", requirementID))
	res, err := b.Db.Exec(deleteRequirementQuery, requirementID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting requirement from database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, fmt.Sprintf("Expected 1 row to be updated, but got %d", count))
		if count == 0 {
			return types.ErrRequirementNotFound
		}
	}
	return nil
}
