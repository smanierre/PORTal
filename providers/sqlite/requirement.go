package sqlite

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func (p Provider) AddRequirement(r types.Requirement) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding requirement to database", slog.Any("requirement", r))
	_, err := p.Db.Exec(addRequirementQuery, r.ID, r.Name, r.Description, r.Notes, r.DaysValidFor, r.Reference.ID)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: requirement.name") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Requirement with given name already exists")
		return backend.ErrDuplicateRequirement
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting requirement into database", slog.String("error", err.Error()))
		return err
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Mapping reference to requirement")
	return nil
}

func (p Provider) GetRequirement(id string) (types.Requirement, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting requirement from database", slog.String("requirement_id", id))
	row := p.Db.QueryRow(getRequirementQuery, id)
	r := types.Requirement{}
	unUsedRefId := ""
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.Notes, &r.DaysValidFor, &unUsedRefId, &r.Reference.ID, &r.Reference.Name, &r.Reference.Volume, &r.Reference.Paragraph)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "No results found for requirement with given id")
		return types.Requirement{}, fmt.Errorf("%w: requirement_id=%s", backend.ErrRequirementNotFound, id)
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning requirement into struct", slog.String("error", err.Error()))
		return types.Requirement{}, err
	}
	return r, nil
}

func (p Provider) GetAllRequirements() ([]types.Requirement, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all requirements from database")
	rows, err := p.Db.Query(getAllRequirementsQuery)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all requirements from database: %s", slog.String("error", err.Error()))
		return nil, err
	}
	var reqs []types.Requirement
	var r types.Requirement
	var unUsedRefId string
	for rows.Next() {
		err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.Notes, &r.DaysValidFor, &unUsedRefId, &r.Reference.ID, &r.Reference.Name, &r.Reference.Volume, &r.Reference.Paragraph)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning requirement into struct", slog.String("error", err.Error()))
			continue
		}
		reqs = append(reqs, r)
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d requirements", len(reqs)))
	return reqs, nil
}

func (p Provider) GetQualificationIDsForRequirement(requirementID string) ([]string, error) {
	rows, err := p.Db.Query(getQualificationsForRequirementQuery, requirementID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting qualifications for requirement", slog.String("error", err.Error()))
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
	return ids, nil
}

func (p Provider) UpdateRequirement(r types.Requirement) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating requirement", slog.Any("new_requirement", r))
	res, err := p.Db.Exec(updateRequirementQuery, r.Name, r.Description, r.Notes, r.DaysValidFor, r.Reference.ID, r.ID)
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Provided reference doesn't exist")
		return backend.ErrReferenceNotFound
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating requirement in database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated but didn't get that")
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking to see if requirement exists")
		if _, err := p.GetRequirement(r.ID); errors.Is(err, backend.ErrRequirementNotFound) {
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Can't determine why row wasn't updated")
		return errors.New("couldn't update requirement")
	}
	return nil
}

func (p Provider) DeleteRequirement(id string) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting requirement from database", slog.String("requirement_id", id))
	res, err := p.Db.Exec(deleteRequirementQuery, id)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting requirement from database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, fmt.Sprintf("Expected 1 row to be updated, but got %d", count))
		if count == 0 {
			return backend.ErrRequirementNotFound
		}
	}
	return nil
}
