package qualificationstore

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (q QualificationStore) AddQualification(qual types.Qualification) (types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID for new qualification")
	qual.ID = uuid.NewString()
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing args")
	if err := checkQualificationForMissingArgs(qual); err != nil {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Missing required arguments", slog.String("error", err.Error()))
		return types.Qualification{}, err
	}
	if qual.Expires && qual.ExpirationInterval < 1 {
		q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided expiration days is invalid", slog.Duration("days", qual.ExpirationInterval))
		return types.Qualification{}, backend.ErrInvalidQualExpiration
	}
	if qual.Expires {
		qual.ExpirationInterval = -1
	}
	// Check all the requirements for any new ones. If there are new ones, create them first before updating the qualification
	for i, req := range qual.InitialRequirements {
		if req.ID == "" {
			addedRequirement, err := q.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new initial requirement for qualification: %w", err)
			}
			qual.InitialRequirements[i] = addedRequirement
		}
	}

	for i, req := range qual.RecurringRequirements {
		if req.ID == "" {
			addedRequirement, err := q.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new recurring requirement for qualification: %w", err)
			}
			qual.RecurringRequirements[i] = addedRequirement
		}
	}

	return qual, q.provider.AddQualification(qual)
}

func (q QualificationStore) GetQualification(id string) (types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualification", slog.String("id", id))
	return q.provider.GetQualification(id)
}

func (q QualificationStore) GetAllQualifications() ([]types.Qualification, error) {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all qualifications")
	quals, err := q.provider.GetAllQualifications()
	if err != nil {
		return []types.Qualification{}, err
	}
	return quals, nil
}

func (q QualificationStore) UpdateQualification(qual types.Qualification) (types.Qualification, error) {
	// Expiration interval should never equal 0 or less if the qualification expires.
	if qual.Expires && qual.ExpirationInterval < 1 {
		q.logger.LogAttrs(context.Background(), slog.LevelWarn, "Invalid expiration days")
		return types.Qualification{}, fmt.Errorf("%w: invalid expiration days", backend.ErrBadUpdate)
	}
	// If the qualification doesn't expire, explicitly set the expiration to -1
	if !qual.Expires {
		qual.ExpirationInterval = -1
	}

	// Check all the requirements for any new ones. If there are new ones, create them first before updating the qualification
	for i, req := range qual.InitialRequirements {
		if req.ID == "" {
			addedRequirement, err := q.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new initial requirement for qualification: %w", err)
			}
			qual.InitialRequirements[i] = addedRequirement
		}
	}

	for i, req := range qual.RecurringRequirements {
		if req.ID == "" {
			addedRequirement, err := q.AddRequirement(req)
			if err != nil {
				return types.Qualification{}, fmt.Errorf("error adding new recurring requirement for qualification: %w", err)
			}
			qual.RecurringRequirements[i] = addedRequirement
		}
	}
	err := q.provider.UpdateQualification(qual)
	if err != nil {
		return types.Qualification{}, err
	}
	newQual, err := q.provider.GetQualification(qual.ID)
	if err != nil {
		return types.Qualification{}, err
	}
	return newQual, nil
}

func (q QualificationStore) DeleteQualification(id string) error {
	q.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting qualification")
	return q.provider.DeleteQualification(id)
}
