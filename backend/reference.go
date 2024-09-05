package backend

import (
	"PORTal/types"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

func (b Backend) AddReference(r types.Reference) (types.Reference, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID for new reference")
	r.ID = uuid.NewString()
	if err := CheckReferenceForMissingArgs(r); err != nil {
		return types.Reference{}, err
	}
	return r, b.requirementProvider.AddReference(r)
}

func (b Backend) GetReference(id string) (types.Reference, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting reference", slog.String("reference_id", id))
	return b.requirementProvider.GetReference(id)
}

func (b Backend) GetReferences() ([]types.Reference, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all references")
	return b.requirementProvider.GetReferences()
}

func (b Backend) UpdateReference(r types.Reference, overrideNoVolume bool) (types.Reference, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting reference to determine updates")
	ref, err := b.GetReference(r.ID)
	if err != nil {
		return types.Reference{}, err
	}
	ref = ref.MergeIn(r, overrideNoVolume)
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating reference", slog.Any("new_reference", ref))
	if err := b.requirementProvider.UpdateReference(ref); err != nil {
		return types.Reference{}, err
	}
	return ref, nil
}

func (b Backend) DeleteReference(id string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting reference", slog.String("reference_id", id))
	return b.requirementProvider.DeleteReference(id)
}
