package sqlite

import (
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func (b *Backend) AddQualification(q types.Qualification) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Inserting qualification into database", slog.Any("qualification", q))
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args")
	if err := q.CheckForMissingArgs(); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Missing required arguments", slog.String("error", err.Error()))
		return err
	}
	_, err := b.Db.Exec(insertQualificationQuery, q.ID, q.Name, q.Notes, q.Expires, q.ExpirationDays)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting qualification into database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (b *Backend) GetQualification(id string) (types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualification from database", slog.String("id", id))
	row := b.Db.QueryRow(getQualificationQuery, id)
	var q types.Qualification
	err := row.Scan(&q.ID, &q.Name, &q.Notes, &q.Expires, &q.ExpirationDays)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Could not find qualification with given id")
		return types.Qualification{}, types.ErrQualificationNotFound
	}
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Error scanning qualification into struct", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	return q, nil
}

func (b *Backend) GetAllQualifications() ([]types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all qualifications from database")
	rows, err := b.Db.Query(getAllQualificationsQuery)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all qualifications from database", slog.String("error", err.Error()))
		return nil, err
	}
	var quals []types.Qualification
	for rows.Next() {
		var qual types.Qualification
		err := rows.Scan(&qual.ID, &qual.Name, &qual.Notes, &qual.Expires, &qual.ExpirationDays)
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning qualification into struct", slog.String("error", err.Error()))
			continue
		}
		quals = append(quals, qual)
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d qualifications in database", len(quals)))
	return quals, nil
}

func (b *Backend) UpdateQualification(q types.Qualification) error {
	if q.Expires && q.ExpirationDays == 0 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Invalid expiration days")
		return fmt.Errorf("%w: invalid expiration days", types.ErrBadUpdate)
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating qualification", slog.Any("qualification", q))
	res, err := b.Db.Exec(updateQualificationQuery, q.Name, q.Notes, q.Expires, q.ExpirationDays, q.ID)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating qualification in database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Expected to update 1 qualification but got 0")
		return types.ErrQualificationNotFound
	}
	return nil
}

func (b *Backend) DeleteQualification(id string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification from database", slog.String("id", id))
	res, err := b.Db.Exec(deleteQualificationQuery, id)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "error deleting qualification from database", slog.String("error", err.Error()))
		return err
	}
	if count, _ := res.RowsAffected(); count != 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "no qualification with that ID exists to be deleted")
		return types.ErrQualificationNotFound
	}
	return nil
}
