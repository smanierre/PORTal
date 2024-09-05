package sqlite

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func (p Provider) AddQualification(q types.Qualification) error {
	_, err := p.Db.Exec(insertQualificationQuery, q.ID, q.Name, q.Notes, q.Expires, q.ExpirationDays)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting qualification into database", slog.String("error", err.Error()))
		return err
	}
	for _, initialRequirement := range q.InitialRequirements {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding initial requirement to Qualification",
			slog.String("qualification_id", q.ID), slog.String("requirement_id", initialRequirement.ID))
		_, err = p.Db.Exec(insertQualificationInitialRequirementQuery, q.ID, initialRequirement.ID)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding requirement to qualification", slog.String("error", err.Error()))
			return err
		}
	}
	for _, recurringRequirement := range q.RecurringRequirements {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding recurring requirement to Qualification",
			slog.String("qualification_id", q.ID), slog.String("requirement_id", recurringRequirement.ID))
		_, err = p.Db.Exec(insertQualificationRecurringRequirementQuery, q.ID, recurringRequirement.ID)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding requirement to qualification", slog.String("error", err.Error()))
			return err
		}
	}
	return nil
}

func (p Provider) GetQualification(id string) (types.Qualification, error) {
	row := p.Db.QueryRow(getQualificationQuery, id)
	var q types.Qualification
	err := row.Scan(&q.ID, &q.Name, &q.Notes, &q.Expires, &q.ExpirationDays)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Could not find qualification with given id")
		return types.Qualification{}, backend.ErrQualificationNotFound
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Error scanning qualification into struct", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Retrieving initial requirements")
	rows, err := p.Db.Query(getInitialRequirementIdsQuery, q.ID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting initial requirement IDs for qualification", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	var initialRequirements []types.Requirement
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning initial requirement ID string", slog.String("error", err.Error()))
			return types.Qualification{}, err
		}
		req, err := p.GetRequirement(id)
		if err != nil {
			return types.Qualification{}, err
		}
		initialRequirements = append(initialRequirements, req)
	}
	rows.Close()
	q.InitialRequirements = initialRequirements
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Retreiving recurring requirements")
	rows, err = p.Db.Query(getRecurringRequirementIdsQuery, q.ID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting recurring requirement IDs for qualification", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	var recurringRequirements []types.Requirement
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning recurring requirement ID string", slog.String("error", err.Error()))
			return types.Qualification{}, err
		}
		req, err := p.GetRequirement(id)
		if err != nil {
			return types.Qualification{}, err
		}
		recurringRequirements = append(recurringRequirements, req)
	}
	rows.Close()
	q.RecurringRequirements = recurringRequirements
	return q, nil
}

func (p Provider) GetAllQualifications() ([]types.Qualification, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all qualification IDs from database")
	rows, err := p.Db.Query(getAllQualificationIDsQuery)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all qualifications from database", slog.String("error", err.Error()))
		return nil, err
	}
	var quals []types.Qualification
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning id into string", slog.String("error", err.Error()))
			continue
		}
		qual, err := p.GetQualification(id)
		if err != nil {
			return nil, err
		}
		quals = append(quals, qual)
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d qualifications in database", len(quals)))
	return quals, nil
}

func (p Provider) UpdateQualification(q types.Qualification) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating qualification", slog.Any("qualification", q))
	tx, err := p.Db.Begin()
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error creating transaction for UpdateQualification", slog.String("error", err.Error()))
	}
	res, err := tx.Exec(updateQualificationQuery, q.Name, q.Notes, q.Expires, q.ExpirationDays, q.ID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating qualification in database", slog.String("error", err.Error()))
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Expected to update 1 qualification but got 0")
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return backend.ErrQualificationNotFound
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting existing initial requirement IDs")
	var existingInitialIDs []string
	rows, err := p.Db.Query(getInitialRequirementIdsQuery, q.ID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting existing initial requirements for qualification", slog.String("error", err.Error()))
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning ID into string", slog.String("error", err.Error()))
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		existingInitialIDs = append(existingInitialIDs, id)
	}
	rows.Close()
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d existing initial requirements", len(existingInitialIDs)))
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting existing recurring requirement IDs")
	var existingRecurringIds []string
	rows, err = p.Db.Query(getRecurringRequirementIdsQuery, q.ID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting existing recurring requirements for qualification", slog.String("error", err.Error()))
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning ID into string", slog.String("error", err.Error()))
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		existingRecurringIds = append(existingRecurringIds, id)
	}
	rows.Close()
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d existing recurring requirements", len(existingRecurringIds)))
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Determining initial qualifications to be removed")
	var initialIdsToBeRemoved []string
	for _, id := range existingInitialIDs {
		found := false
		for _, newReq := range q.InitialRequirements {
			if id == newReq.ID {
				found = true
			}
		}
		if !found {
			initialIdsToBeRemoved = append(initialIdsToBeRemoved, id)
		}
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Initial requirements set to be removed: %s", initialIdsToBeRemoved))
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Determining recurring qualifications to be removed")
	var recurringIdsToBeRemoved []string
	for _, id := range existingRecurringIds {
		found := false
		for _, newReq := range q.RecurringRequirements {
			if id == newReq.ID {
				found = true
			}
		}
		if !found {
			recurringIdsToBeRemoved = append(recurringIdsToBeRemoved, id)
		}
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Recurring requirements set to be removed: %s", recurringIdsToBeRemoved))
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Removing initial requirements")
	for _, id := range initialIdsToBeRemoved {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Removing initial requirement: %s", id))
		_, err = tx.Exec(deleteQualificationInitialRequirementQuery, id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting initial requirement from qualification", slog.String("error", err.Error()))
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully removed initial requirement")
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Removing recurring requirements")
	for _, id := range recurringIdsToBeRemoved {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Removing recurring requirement: %s", id))
		_, err = tx.Exec(deleteQualificationRecurringRequirementQuery, id)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting recurring requirement from qualification", slog.String("error", err.Error()))
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully removed recurring requirement")
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding new initial requirements")
	for _, newReq := range q.InitialRequirements {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Adding initial requirement: %s", newReq.ID))
		_, err := tx.Exec(insertQualificationInitialRequirementQuery, q.ID, newReq.ID)
		if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Requirement already assigned to qualification, skipping")
			continue
		}
		if err != nil {
			var errToReturn error
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Initial requirement to be added not found")
				errToReturn = fmt.Errorf("%w: %s", backend.ErrRequirementNotFound, newReq.ID)
			} else {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding new initial requirement to query", slog.String("error", err.Error()))
				p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
				errToReturn = err
			}
			if err := tx.Rollback(); err != nil {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction",
					slog.String("error", err.Error()), slog.String("original_error", err.Error()))
				return err
			}
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return errToReturn
		}
	}
	for _, newReq := range q.RecurringRequirements {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Adding recurring requirement: %s", newReq.ID))
		_, err := tx.Exec(insertQualificationRecurringRequirementQuery, q.ID, newReq.ID)
		if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Requirement already assigned to qualification, skipping")
			continue
		}
		if err != nil {
			var errToReturn error
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Recurring requirement to be added not found")
				errToReturn = fmt.Errorf("%w: %s", backend.ErrRequirementNotFound, newReq.ID)
			} else {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding new recurring requirement to query", slog.String("error", err.Error()))
				p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
				errToReturn = err
			}
			if err := tx.Rollback(); err != nil {
				p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction",
					slog.String("error", err.Error()), slog.String("original_error", err.Error()))
				return err
			}
			p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return errToReturn
		}
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Committing transaction")
	if err = tx.Commit(); err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error committing transaction", slog.String("error", err.Error()))
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully committed transaction")
	return nil
}

func (p Provider) DeleteQualification(id string) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification from database", slog.String("id", id))
	res, err := p.Db.Exec(deleteQualificationQuery, id)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "error deleting qualification from database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "no qualification with that ID exists to be deleted")
		return backend.ErrQualificationNotFound
	}
	return nil
}
