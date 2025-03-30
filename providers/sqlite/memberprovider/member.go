package memberprovider

import (
	"PORTal/backend"
	"PORTal/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func (m MemberProvider) AddMember(mem types.Member) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Inserting member into database", slog.Any("member", m))
	var err error
	if mem.SupervisorID == "" {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "No supervisor provided, setting to null in database")
		_, err = m.db.Exec(insertMemberQuery, mem.ID, mem.FirstName, mem.LastName, mem.Grade, mem.Username, nil, mem.Admin, mem.Hash)
	} else {
		_, err = m.db.Exec(insertMemberQuery, mem.ID, mem.FirstName, mem.LastName, mem.Grade, mem.Username, mem.SupervisorID, mem.Admin, mem.Hash)
	}
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Provided supervisor id doesn't exist", slog.String("supervisor_id", mem.ID))
		return fmt.Errorf("%w: %s", backend.ErrSupervisorNotFound, mem.SupervisorID)
	} else if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: member.user_name") {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided username is taken")
		return backend.ErrDuplicateUsername
	} else if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting member into database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (m MemberProvider) GetMember(identifier string, method backend.ProviderMethod) (types.Member, error) {
	var row *sql.Row
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from database")
	switch method {
	case backend.ById:
		row = m.db.QueryRow(getMemberQuery, identifier)
		break
	case backend.ByUsername:
		row = m.db.QueryRow(getMemberByUsernameQuery, identifier)
	default:
		return types.Member{}, errors.New(fmt.Sprintf("unexpected retrieval method: %d", method))
	}
	var mem types.Member
	supervisorId := sql.NullString{}
	err := row.Scan(&mem.ID, &mem.FirstName, &mem.LastName, &mem.Grade, &mem.Username, &supervisorId, &mem.Admin, &mem.Hash, &mem.Disabled)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "No user found with given identifier")
		return types.Member{}, backend.ErrMemberNotFound
	} else if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning member into struct", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	if supervisorId.Valid {
		mem.SupervisorID = supervisorId.String
	}
	return mem, nil
}

func (m MemberProvider) GetMemberFromSession(sessionID string) (types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from session", slog.String("session_id", sessionID))
	row := m.db.QueryRow(getMemberFromSessionQuery, sessionID)
	var memberID string
	err := row.Scan(&memberID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting member id from session", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	return m.GetMember(memberID, backend.ById)
}

func (m MemberProvider) GetAllMembers() ([]types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all members from database")
	rows, err := m.db.Query(getAllMembersQuery)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all members from database", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var members []types.Member
	for rows.Next() {
		mem := types.Member{}
		supervisorId := sql.NullString{}
		err = rows.Scan(&mem.ID, &mem.FirstName, &mem.LastName, &mem.Grade, &mem.Username, &supervisorId, &mem.Admin, &mem.Hash)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning member into struct", slog.String("error", err.Error()))
			continue
		}
		if supervisorId.Valid {
			mem.SupervisorID = supervisorId.String
		}
		members = append(members, mem)
	}
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Retreived %d members from database", len(members)))
	return members, nil
}

func (m MemberProvider) GetDisabledMembers() ([]types.Member, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all disabled members from database")
	rows, err := m.db.Query(getDisabledMembersQuery)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error when getting disabled members from database", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var members []types.Member
	var mem types.Member
	var supervisorId sql.NullString
	for rows.Next() {
		err := rows.Scan(&mem.ID, &mem.FirstName, &mem.LastName, &mem.Grade, &mem.Username, &supervisorId, &mem.Admin, &mem.Hash)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error when scanning disabled member into struct", slog.String("error", err.Error()))
			return nil, err
		}
		if supervisorId.Valid {
			mem.SupervisorID = supervisorId.String
		}
		members = append(members, mem)
	}
	return members, nil
}

func (m MemberProvider) GetSubordinates(memberID string) ([]types.Member, error) {
	rows, err := m.db.Query(getSubordinatesQuery, memberID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting subordinates for member", slog.String("error", err.Error()))
		return nil, err
	}
	var subordinates []types.Member
	var subordinate types.Member
	for rows.Next() {
		err = rows.Scan(&subordinate.ID, &subordinate.FirstName, &subordinate.LastName, &subordinate.Grade, &subordinate.Username, &subordinate.SupervisorID, &subordinate.Admin, &subordinate.Hash)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error when scanning subordinate into struct", slog.String("error", err.Error()))
			return nil, err
		}
		subordinates = append(subordinates, subordinate)
	}
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d subordinates for member", len(subordinates)))
	return subordinates, nil
}

// RemoveSubordinates takes a member's ID and blanks out the SupervisorID of any members that have it as their supervisor
func (m MemberProvider) RemoveSubordinates(id string) error {
	_, err := m.db.Exec(removeSubordinatesQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func (m MemberProvider) UpdateMember(mem types.Member) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating member", slog.Any("member", m))
	var res sql.Result
	var err error
	if mem.SupervisorID == "" {
		m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Supervisor ID is empty, inserting as null in database")
		res, err = m.db.Exec(updateMemberQuery, mem.FirstName, mem.LastName, mem.Grade, nil, mem.Admin, mem.Hash, mem.ID)
	} else {
		res, err = m.db.Exec(updateMemberQuery, mem.FirstName, mem.LastName, mem.Grade, mem.SupervisorID, mem.Admin, mem.Hash, mem.ID)
	}
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Attempting to update member with non-existent supervisor")
		return backend.ErrSupervisorNotFound
	}
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating member", slog.String("error", err.Error()))
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated for member, got 0")
		return backend.ErrMemberNotFound
	}
	return nil
}

func (m MemberProvider) DeleteMember(identifier string, method backend.ProviderMethod) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member", slog.String("identifier", identifier))
	var res sql.Result
	var err error
	switch method {
	case backend.ById:
		res, err = m.db.Exec(deleteMemberQuery, identifier)
	case backend.ByUsername:
		res, err = m.db.Exec(deleteMemberByUsernameQuery, identifier)
	}
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting member", slog.String("error", err.Error()))
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		m.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated for member, got 0")
		return backend.ErrMemberNotFound
	}
	return nil
}

func (m MemberProvider) DisableMember(id string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Disabling member", slog.String("identifier", id))
	_, err := m.db.Exec(disableMemberQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func (m MemberProvider) EnableMember(id string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Enabling member", slog.String("identifier", id))
	res, err := m.db.Exec(enableMemberQuery, id)
	if err != nil {
		return err
	}
	if updated, _ := res.RowsAffected(); updated == 0 {
		return backend.ErrMemberNotFound
	}
	return nil
}
