package memberstore

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (m MemberStore) Login(username, password string) (types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Attempting to login member", slog.String("username", username))
	member, err := m.GetMember(username)
	if err != nil || member.Disabled {
		return types.Member{}, backend.ErrAuthenticationFailed
	}
	if err = bcrypt.CompareHashAndPassword([]byte(member.Hash), []byte(password)); err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Password validation failed", slog.String("error", err.Error()))
		return types.Member{}, backend.ErrAuthenticationFailed
	}
	return member, nil
}

func (m MemberStore) AddMember(mem types.Member) (types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding member", slog.Any("member", mem))
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Generating ID")
	mem.ID = uuid.NewString()

	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing fields")
	if err := checkMemberForMissingArgs(mem); err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Required arguments missing for user creation", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	if len(mem.Password) < backend.MinimumPwLength {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Password length %d does not meet minimum length of %d", len(mem.Password), backend.MinimumPwLength))
		return types.Member{}, backend.ErrWeakPassword
	}
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Hashing password")
	hash, err := bcrypt.GenerateFromPassword([]byte(mem.Password), m.hashCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided password is too long", slog.Int("length", len(mem.Password)))
		return types.Member{}, backend.ErrPasswordTooLong
	}
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error hashing password for new user", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	mem.Hash = string(hash)
	mem.Password = ""
	err = m.provider.AddMember(mem)
	if err != nil {
		return types.Member{}, err
	}
	return mem, nil
}

func (m MemberStore) GetMember(identifier string) (types.Member, error) {
	l := m.logger.With(slog.String("identifier", identifier))
	l.LogAttrs(context.Background(), slog.LevelInfo, "Determining method to get member with")
	var method backend.ProviderMethod
	if _, err := uuid.Parse(identifier); err != nil {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Using method ByUsername")
		method = backend.ByUsername
	} else {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Using method ById")
		method = backend.ById
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from database")
	member, err := m.provider.GetMember(identifier, method)
	if err != nil {
		return member, err
	}
	if member.Disabled {
		return types.Member{}, backend.ErrMemberDisabled
	}
	return member, nil
}

func (m MemberStore) GetPotentialSupervisors(member types.Member, grade types.Grade) ([]types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting potential supervisors")
	allMembers, err := m.GetAllMembers()
	if err != nil {
		return nil, err
	}
	var potentialSupervisors []types.Member
	for _, mem := range allMembers {
		if member.ID == mem.ID {
			continue
		}
		if mem.Grade.CanSupervise(grade) {
			potentialSupervisors = append(potentialSupervisors, mem)
		}
	}
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d potential supervisors", len(potentialSupervisors)))
	return potentialSupervisors, nil
}

func (m MemberStore) GetMemberFromSession(sessionID string) (types.Member, error) {
	return m.provider.GetMemberFromSession(sessionID)
}

func (m MemberStore) GetDisabledMember(identifier string) (types.Member, error) {
	l := m.logger.With(slog.String("identifier", identifier))
	l.LogAttrs(context.Background(), slog.LevelInfo, "Getting disabled member from database")
	var method backend.ProviderMethod
	if _, err := uuid.Parse(identifier); err != nil {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Using method ByUsername")
		method = backend.ByUsername
	} else {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Using method ById")
		method = backend.ById
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from database")
	member, err := m.provider.GetMember(identifier, method)
	if err != nil {
		return member, err
	}
	if !member.Disabled {
		return types.Member{}, backend.ErrMemberNotFound
	}
	return member, nil
}

func (m MemberStore) GetAllMembers() ([]types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all members")
	return m.provider.GetAllMembers()
}

func (m MemberStore) GetDisabledMembers() ([]types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all disabled members")
	members, err := m.provider.GetDisabledMembers()
	if err != nil {
		return nil, err
	}
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d disabled members", len(members)))
	return members, nil
}

func (m MemberStore) GetSubordinates(memberID string) []types.Member {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting subordinates for member", slog.String("member_id", memberID))
	return m.provider.GetSubordinates(memberID)
}

func (m MemberStore) UpdateMember(mem types.Member) (types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating member")
	if mem.Password != "" {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "New password provided, verifying it meets requirements")
		if len(mem.Password) < backend.MinimumPwLength {
			m.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Password length %d does not meet minimum length of %d", len(mem.Password), backend.MinimumPwLength))
			return types.Member{}, backend.ErrWeakPassword
		}
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Hashing new password")
		hash, err := bcrypt.GenerateFromPassword([]byte(mem.Password), m.hashCost)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided password is too long", slog.Int("length", len(mem.Password)))
				return types.Member{}, backend.ErrPasswordTooLong
			}
			if err != nil {
				m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Error hashing password for new user", slog.String("error", err.Error()))
				return types.Member{}, err
			}
		}
		mem.Hash = string(hash)
		mem.Password = ""
	} else {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting previous hash so it doesn't get blanked out.")
		currentMember, err := m.provider.GetMember(mem.ID, backend.ById)
		if err != nil {
			return types.Member{}, fmt.Errorf("error getting previous hash for member: %w", err)
		}
		mem.Hash = currentMember.Hash
	}
	err := m.provider.UpdateMember(mem)
	if err != nil {
		return types.Member{}, err
	}
	updatedMember, err := m.provider.GetMember(mem.ID, backend.ById)
	if err != nil {
		return types.Member{}, err
	}
	return updatedMember, nil
}

func (m MemberStore) DeleteMember(identifier string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member", slog.String("identifier", identifier))
	var method backend.ProviderMethod
	if _, err := uuid.Parse(identifier); err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Setting delete method to ByUsername")
		method = backend.ByUsername
	} else {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Setting delete method to ById")
		method = backend.ById
	}
	return m.provider.DeleteMember(identifier, method)
}

func (m MemberStore) DisableMember(id string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Disabling member", slog.String("id", id))
	err := m.provider.DisableMember(id)
	if err != nil {
		return err
	}
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Removing member as supervisor for any other members")
	return m.provider.RemoveSubordinates(id)
}

func (m MemberStore) EnableMember(id string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Enabling member", slog.String("id", id))
	err := m.provider.EnableMember(id)
	if err != nil {
		return err
	}
	return nil
}
