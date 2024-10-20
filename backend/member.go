package backend

import (
	"PORTal/types"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (b Backend) AddMember(m types.Member) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding member", slog.Any("member", m))
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID")
	m.ID = uuid.NewString()

	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing fields")
	if err := CheckMemberForMissingArgs(m); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Required arguments missing for user creation", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	if len(m.Password) < MinimumPwLength {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Password length %d does not meet minimum length of %d", len(m.Password), MinimumPwLength))
		return types.Member{}, ErrWeakPassword
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Hashing password")
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Password), b.config.BcryptCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided password is too long", slog.Int("length", len(m.Password)))
		return types.Member{}, ErrPasswordTooLong
	}
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error hashing password for new user", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	m.Hash = string(hash)
	m.Password = ""
	err = b.memberProvider.AddMember(m)
	if err != nil {
		return types.Member{}, err
	}
	return m, nil
}

func (b Backend) GetMember(identifier string) (types.Member, error) {
	l := b.logger.With(slog.String("identifier", identifier))
	l.LogAttrs(context.Background(), slog.LevelInfo, "Determining method to get member with")
	var m ProviderMethod
	if _, err := uuid.Parse(identifier); err != nil {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Using method ByUsername")
		m = ByUsername
	} else {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Using method ById")
		m = ById
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from database")
	return b.memberProvider.GetMember(identifier, m)
}

func (b Backend) GetAllMembers() ([]types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all members")
	return b.memberProvider.GetAllMembers()
}

func (b Backend) GetSubordinates(memberID string) ([]types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting subordinates for member", slog.String("member_id", memberID))
	return b.memberProvider.GetSubordinates(memberID)
}

func (b Backend) UpdateMember(m types.Member) (types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating member")
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting previous member to determine updates")
	previousMember, err := b.memberProvider.GetMember(m.ID, ById)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Unable to get previous member to compare updates")
		return types.Member{}, err
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Merging members to determine updates")
	updateMember := previousMember.MergeIn(m)
	if updateMember.Password != "" {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "New password provided, verifying it meets requirements")
		if len(m.Password) < MinimumPwLength {
			b.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Password length %d does not meet minimum length of %d", len(m.Password), MinimumPwLength))
			return types.Member{}, ErrWeakPassword
		}
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Hashing new password")
		hash, err := bcrypt.GenerateFromPassword([]byte(updateMember.Password), b.config.BcryptCost)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided password is too long", slog.Int("length", len(m.Password)))
				return types.Member{}, ErrPasswordTooLong
			}
			if err != nil {
				b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error hashing password for new user", slog.String("error", err.Error()))
				return types.Member{}, err
			}
		}
		updateMember.Hash = string(hash)
		updateMember.Password = ""
	}
	err = b.memberProvider.UpdateMember(updateMember)
	if err != nil {
		return types.Member{}, err
	}
	return updateMember, nil
}

func (b Backend) DeleteMember(identifier string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member", slog.String("identifier", identifier))
	var m ProviderMethod
	if _, err := uuid.Parse(identifier); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Setting delete method to ByUsername")
		m = ByUsername
	} else {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Setting delete method to ById")
		m = ById
	}
	return b.memberProvider.DeleteMember(identifier, m)
}
