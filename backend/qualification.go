package backend

import (
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (b Backend) AddQualification(q types.Qualification) (types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID for new qualification")
	q.ID = uuid.NewString()
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args")
	if err := CheckQualificationForMissingArgs(q); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Missing required arguments", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	if q.Expires && q.ExpirationDays < 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided expiration days is invalid", slog.Int("days", q.ExpirationDays))
		return types.Qualification{}, ErrInvalidQualExpiration
	}
	return q, b.qualificationProvider.AddQualification(q)
}

func (b Backend) GetQualification(id string) (types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualification", slog.String("id", id))
	return b.qualificationProvider.GetQualification(id)
}

func (b Backend) GetAllQualifications() ([]types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all qualifications")
	quals, err := b.qualificationProvider.GetAllQualifications()
	if err != nil {
		return []types.Qualification{}, err
	}
	return quals, nil
}

func (b Backend) UpdateQualification(q types.Qualification) (types.Qualification, error) {
	if q.Expires && q.ExpirationDays == 0 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Invalid expiration days")
		return types.Qualification{}, fmt.Errorf("%w: invalid expiration days", ErrBadUpdate)
	}
	err := b.qualificationProvider.UpdateQualification(q)
	if err != nil {
		return types.Qualification{}, err
	}
	newQual, err := b.qualificationProvider.GetQualification(q.ID)
	if err != nil {
		return types.Qualification{}, err
	}
	return newQual, nil
}

func (b Backend) DeleteQualification(id string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification")
	return b.qualificationProvider.DeleteQualification(id)
}
