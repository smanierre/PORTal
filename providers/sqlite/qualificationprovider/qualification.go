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
	_, err := q.db.Exec(insertQualificationQuery, qual.ID, qual.Name, qual.Notes, qual.Expires, qual.ExpirationInterval)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting qualification into database", slog.String("error", err.Error()))
		return err
	}
	for _, initialRequirement := range qual.InitialRequirements {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding initial requirement to Qualification",
			slog.String("qualification_id", qual.ID), slog.String("requirement_id", initialRequirement.ID))
		_, err = q.db.Exec(insertQualificationInitialRequirementQuery, qual.ID, initialRequirement.ID)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding requirement to qualification", slog.String("error", err.Error()))
			return err
		}
	}
	for _, recurringRequirement := range qual.RecurringRequirements {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding recurring requirement to Qualification",
			slog.String("qualification_id", qual.ID), slog.String("requirement_id", recurringRequirement.ID))
		_, err = q.db.Exec(insertQualificationRecurringRequirementQuery, qual.ID, recurringRequirement.ID)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding requirement to qualification", slog.String("error", err.Error()))
			return err
		}
	}
	return nil
}

func (q QualificationProvider) GetQualification(id string) (types.Qualification, error) {
	row := q.db.QueryRow(getQualificationQuery, id)
	var qual types.Qualification
	err := row.Scan(&qual.ID, &qual.Name, &qual.Notes, &qual.Expires, &qual.ExpirationInterval)
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
	var initialRequirements []types.Requirement
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
	var recurringRequirements []types.Requirement
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
	tx, err := q.db.Begin()
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error creating transaction for UpdateQualification", slog.String("error", err.Error()))
	}
	res, err := tx.Exec(updateQualificationQuery, qual.Name, qual.Notes, qual.Expires, qual.ExpirationInterval, qual.ID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating qualification in database", slog.String("error", err.Error()))
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Expected to update 1 qualification but got 0")
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return backend.ErrQualificationNotFound
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting existing initial requirement IDs")
	var existingInitialIDs []string
	rows, err := q.db.Query(getInitialRequirementIdsQuery, qual.ID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting existing initial requirements for qualification", slog.String("error", err.Error()))
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning ID into string", slog.String("error", err.Error()))
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		existingInitialIDs = append(existingInitialIDs, id)
	}
	rows.Close()
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d existing initial requirements", len(existingInitialIDs)))
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting existing recurring requirement IDs")
	var existingRecurringIds []string
	rows, err = q.db.Query(getRecurringRequirementIdsQuery, qual.ID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting existing recurring requirements for qualification", slog.String("error", err.Error()))
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning ID into string", slog.String("error", err.Error()))
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		existingRecurringIds = append(existingRecurringIds, id)
	}
	rows.Close()
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d existing recurring requirements", len(existingRecurringIds)))
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Determining initial qualifications to be removed")
	var initialIdsToBeRemoved []string
	for _, id := range existingInitialIDs {
		found := false
		for _, newReq := range qual.InitialRequirements {
			if id == newReq.ID {
				found = true
			}
		}
		if !found {
			initialIdsToBeRemoved = append(initialIdsToBeRemoved, id)
		}
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Initial requirements set to be removed: %s", initialIdsToBeRemoved))
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Determining recurring qualifications to be removed")
	var recurringIdsToBeRemoved []string
	for _, id := range existingRecurringIds {
		found := false
		for _, newReq := range qual.RecurringRequirements {
			if id == newReq.ID {
				found = true
			}
		}
		if !found {
			recurringIdsToBeRemoved = append(recurringIdsToBeRemoved, id)
		}
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Recurring requirements set to be removed: %s", recurringIdsToBeRemoved))
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Removing initial requirements")
	for _, id := range initialIdsToBeRemoved {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Removing initial requirement: %s", id))
		_, err = tx.Exec(deleteQualificationInitialRequirementQuery, id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting initial requirement from qualification", slog.String("error", err.Error()))
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully removed initial requirement")
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Removing recurring requirements")
	for _, id := range recurringIdsToBeRemoved {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Removing recurring requirement: %s", id))
		_, err = tx.Exec(deleteQualificationRecurringRequirementQuery, id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting recurring requirement from qualification", slog.String("error", err.Error()))
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
			if err = tx.Rollback(); err != nil {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
				return err
			}
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully removed recurring requirement")
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding new initial requirements")
	for _, newReq := range qual.InitialRequirements {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Adding initial requirement: %s", newReq.ID))
		_, err := tx.Exec(insertQualificationInitialRequirementQuery, qual.ID, newReq.ID)
		if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Requirement already assigned to qualification, skipping")
			continue
		}
		if err != nil {
			var errToReturn error
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Initial requirement to be added not found")
				errToReturn = fmt.Errorf("%w: %s", backend.ErrRequirementNotFound, newReq.ID)
			} else {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding new initial requirement to query", slog.String("error", err.Error()))
				q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
				errToReturn = err
			}
			if err := tx.Rollback(); err != nil {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction",
					slog.String("error", err.Error()), slog.String("original_error", err.Error()))
				return err
			}
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return errToReturn
		}
	}
	for _, newReq := range qual.RecurringRequirements {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Adding recurring requirement: %s", newReq.ID))
		_, err := tx.Exec(insertQualificationRecurringRequirementQuery, qual.ID, newReq.ID)
		if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Requirement already assigned to qualification, skipping")
			continue
		}
		if err != nil {
			var errToReturn error
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Recurring requirement to be added not found")
				errToReturn = fmt.Errorf("%w: %s", backend.ErrRequirementNotFound, newReq.ID)
			} else {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error adding new recurring requirement to query", slog.String("error", err.Error()))
				q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
				errToReturn = err
			}
			if err := tx.Rollback(); err != nil {
				q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction",
					slog.String("error", err.Error()), slog.String("original_error", err.Error()))
				return err
			}
			q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
			return errToReturn
		}
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Committing transaction")
	if err = tx.Commit(); err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error committing transaction", slog.String("error", err.Error()))
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Rolling back transaction")
		if err = tx.Rollback(); err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error rolling back transaction", slog.String("error", err.Error()))
			return err
		}
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully rolled back transaction")
		return err
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully committed transaction")
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

func (q QualificationProvider) GetMemberQualification(memberID, qualificationID string) (types.Qualification, error) {
	// Verify that member has qualification assigned
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking to see if member has qualification assigned")
	row := q.db.QueryRow(checkMemberQualificationQuery, memberID, qualificationID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error checking if member has qualification assigned", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	if count == 0 {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "Member with given qualification not found")
		return types.Qualification{}, backend.ErrMemberQualificationNotFound
	}
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "User has qualification assigned")
	return q.GetQualification(qualificationID)
}

func (q QualificationProvider) GetMemberQualifications(memberID string) ([]types.Qualification, error) {
	rows, err := q.db.Query(getMemberQualificationIDsQuery, memberID)
	if err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting qualification IDs for member", slog.String("error", err.Error()))
		return nil, err
	}
	var ids []string
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			q.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning id into string", slog.String("error", err.Error()))
			return nil, err
		}
		ids = append(ids, id)
	}
	quals := make([]types.Qualification, 0, len(ids))
	var qual types.Qualification
	for _, id := range ids {
		qual, err = q.GetQualification(id)
		if err != nil {
			return nil, err
		}
		quals = append(quals, qual)
	}
	return quals, nil
}
