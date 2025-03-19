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
	if q.Expires && q.ExpirationInterval < 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided expiration days is invalid", slog.Duration("days", q.ExpirationInterval))
		return types.Qualification{}, ErrInvalidQualExpiration
	}
	if q.Expires {
		q.ExpirationInterval = -1
	}
	// Check all the requirements for any new ones. If there are new ones, create them first before updating the qualification
	for i, req := range q.InitialRequirements {
		if req.ID == "" {
			addedRequirement, err := b.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new initial requirement for qualification: %w", err)
			}
			q.InitialRequirements[i] = addedRequirement
		}
	}

	for i, req := range q.RecurringRequirements {
		if req.ID == "" {
			addedRequirement, err := b.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new recurring requirement for qualification: %w", err)
			}
			q.RecurringRequirements[i] = addedRequirement
		}
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
	// Expiration interval should never equal 0 or less if the qualification expires.
	if q.Expires && q.ExpirationInterval < 1 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Invalid expiration days")
		return types.Qualification{}, fmt.Errorf("%w: invalid expiration days", ErrBadUpdate)
	}
	// If the qualification doesn't expire, explicitly set the expiration to -1
	if !q.Expires {
		q.ExpirationInterval = -1
	}

	// Check all the requirements for any new ones. If there are new ones, create them first before updating the qualification
	for i, req := range q.InitialRequirements {
		if req.ID == "" {
			addedRequirement, err := b.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new initial requirement for qualification: %w", err)
			}
			q.InitialRequirements[i] = addedRequirement
		}
	}

	for i, req := range q.RecurringRequirements {
		if req.ID == "" {
			addedRequirement, err := b.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new recurring requirement for qualification: %w", err)
			}
			q.RecurringRequirements[i] = addedRequirement
		}
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
