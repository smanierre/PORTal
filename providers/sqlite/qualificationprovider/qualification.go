package qualificationprovider

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func (q QualificationProvider) AddQualification(qual types.Qualification) error {
	_, err := q.db.Exec(insertQualificationQuery, qual.ID, qual.Name, qual.Notes)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting qualification into database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (q QualificationProvider) AssignRequirementToQualification(qualificationID string, requirementID string, initial bool) error {
	var err error
	if initial {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding initial requirement to Qualification")
		_, err = q.db.Exec(insertQualificationInitialRequirementQuery, qualificationID, requirementID)
	} else {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding recurring requirement to Qualification")
		_, err = q.db.Exec(insertQualificationRecurringRequirementQuery, qualificationID, requirementID)
	}
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding requirement to qualification", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (q QualificationProvider) GetQualification(id string) (types.Qualification, error) {
	row := q.db.QueryRow(getQualificationQuery, id)
	var qual types.Qualification
	err := row.Scan(&qual.ID, &qual.Name, &qual.Notes)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "Could not find qualification with given id")
		return types.Qualification{}, backend.ErrQualificationNotFound
	}
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning qualification into struct", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Retrieving initial requirements")
	rows, err := q.db.Query(getInitialRequirementIdsQuery, qual.ID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting initial requirement IDs for qualification", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	initialRequirements := []types.Requirement{}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning initial requirement ID string", slog.String("error", err.Error()))
			return types.Qualification{}, err
		}
		req, err := q.GetRequirement(id)
		if err != nil {
			return types.Qualification{}, err
		}
		initialRequirements = append(initialRequirements, req)
	}
	rows.Close()
	qual.InitialRequirements = initialRequirements
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Retrieving recurring requirements")
	rows, err = q.db.Query(getRecurringRequirementIdsQuery, qual.ID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting recurring requirement IDs for qualification", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	recurringRequirements := []types.Requirement{}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning recurring requirement ID string", slog.String("error", err.Error()))
			return types.Qualification{}, err
		}
		req, err := q.GetRequirement(id)
		if err != nil {
			return types.Qualification{}, err
		}
		recurringRequirements = append(recurringRequirements, req)
	}
	rows.Close()
	qual.RecurringRequirements = recurringRequirements
	return qual, nil
}

func (q QualificationProvider) GetAllQualifications() ([]types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all qualification IDs from database")
	rows, err := q.db.Query(getAllQualificationIDsQuery)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all qualifications from database", slog.String("error", err.Error()))
		return nil, err
	}
	var quals []types.Qualification
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning id into string", slog.String("error", err.Error()))
			continue
		}
		qual, err := q.GetQualification(id)
		if err != nil {
			return nil, err
		}
		quals = append(quals, qual)
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d qualifications in database", len(quals)))
	return quals, nil
}

func (q QualificationProvider) UpdateQualification(qual types.Qualification) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating qualification", slog.Any("qualification", qual))
	res, err := q.db.Exec(updateQualificationQuery, qual.Name, qual.Notes, qual.ID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating qualification in database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Expected to update 1 qualification but got 0")
		return backend.ErrQualificationNotFound
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully updated requirement")
	return nil
}

func (q QualificationProvider) DeleteQualification(id string) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification from database", slog.String("id", id))
	res, err := q.db.Exec(deleteQualificationQuery, id)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "error deleting qualification from database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "no qualification with that ID exists to be deleted")
		return backend.ErrQualificationNotFound
	}
	return nil
}
