package sqlite

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"log/slog"
	"strings"
)

func (p Provider) AddReference(r types.Reference) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding reference to database", slog.Any("reference", r))
	_, err := p.Db.Exec(addReferenceQuery, r.ID, r.Name, r.Volume, r.Paragraph)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: reference.name") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Reference with that name already exists")
		return backend.ErrDuplicateReference
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting reference into database: %s", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p Provider) GetReference(id string) (types.Reference, error) {
	row := p.Db.QueryRow(getReferenceQuery, id)
	var ref types.Reference
	err := row.Scan(&ref.ID, &ref.Name, &ref.Volume, &ref.Paragraph)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Unable to find reference with given ID")
		return types.Reference{}, backend.ErrReferenceNotFound
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting reference from database", slog.String("error", err.Error()))
		return types.Reference{}, err
	}
	return ref, nil
}

func (p Provider) GetReferences() ([]types.Reference, error) {
	rows, err := p.Db.Query(getReferencesQuery)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting references from database", slog.String("error", err.Error()))
		return nil, err
	}
	var refs []types.Reference
	var ref types.Reference
	for rows.Next() {
		err = rows.Scan(&ref.ID, &ref.Name, &ref.Volume, &ref.Paragraph)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning reference into struct", slog.String("error", err.Error()))
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (p Provider) UpdateReference(r types.Reference) error {
	_, err := p.Db.Exec(updateReferenceQuery, r.Name, r.Volume, r.Paragraph, r.ID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating reference", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p Provider) DeleteReference(id string) error {
	res, err := p.Db.Exec(deleteReferenceQuery, id)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting reference from database", slog.String("error", err.Error()))
		return err
	}
	if updated, _ := res.RowsAffected(); updated != 1 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Didn't get expected 1 row updated, qualification mostly not found")
		return backend.ErrReferenceNotFound
	}
	return nil
}
