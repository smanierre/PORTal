package sqlite

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

func (p Provider) AddMember(m types.Member) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Inserting member into database", slog.Any("member", m))
	var err error
	if m.SupervisorID == "" {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "No supervisor provided, setting to null in database")
		_, err = p.Db.Exec(insertMemberQuery, m.ID, m.FirstName, m.LastName, m.Rank, m.Username, nil, m.Admin, m.Hash)
	} else {
		_, err = p.Db.Exec(insertMemberQuery, m.ID, m.FirstName, m.LastName, m.Rank, m.Username, m.SupervisorID, m.Admin, m.Hash)
	}
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Provided supervisor id doesn't exist", slog.String("supervisor_id", m.ID))
		return fmt.Errorf("%w: %s", backend.ErrSupervisorNotFound, m.SupervisorID)
	} else if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed: member.user_name") {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Provided username is taken")
		return backend.ErrDuplicateUsername
	} else if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting member into database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (p Provider) GetMember(identifier string, method backend.ProviderMethod) (types.Member, error) {
	var row *sql.Row
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from database")
	switch method {
	case backend.ById:
		row = p.Db.QueryRow(getMemberQuery, identifier)
		break
	case backend.ByUsername:
		row = p.Db.QueryRow(getMemberByUsernameQuery, identifier)
	default:
		return types.Member{}, errors.New(fmt.Sprintf("unexpected retrieval method: %d", method))
	}
	var m types.Member
	supervisorId := sql.NullString{}
	err := row.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Rank, &m.Username, &supervisorId, &m.Admin, &m.Hash)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "No user found with given identifier")
		return types.Member{}, backend.ErrMemberNotFound
	} else if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning member into struct", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	if supervisorId.Valid {
		m.SupervisorID = supervisorId.String
	}
	return m, nil
}

func (p Provider) GetAllMembers() ([]types.Member, error) {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all members from database")
	rows, err := p.Db.Query(getAllMembersQuery)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all members from database", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var members []types.Member
	for rows.Next() {
		m := types.Member{}
		supervisorId := sql.NullString{}
		err = rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Rank, &m.Username, &supervisorId, &m.Admin, &m.Hash)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning member into struct", slog.String("error", err.Error()))
			continue
		}
		if supervisorId.Valid {
			m.SupervisorID = supervisorId.String
		}
		members = append(members, m)
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Retreived %d members from database", len(members)))
	return members, nil
}

func (p Provider) GetSubordinates(memberID string) ([]types.Member, error) {
	rows, err := p.Db.Query(getSubordinatesQuery, memberID)
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting subordinates for member", slog.String("error", err.Error()))
		return nil, err
	}
	var subordinates []types.Member
	var subordinate types.Member
	for rows.Next() {
		err = rows.Scan(&subordinate.ID, &subordinate.FirstName, &subordinate.LastName, &subordinate.Rank, &subordinate.Username, &subordinate.SupervisorID, &subordinate.Admin, &subordinate.Hash)
		if err != nil {
			p.logger.LogAttrs(context.Background(), slog.LevelError, "Error when scanning subordinate into struct", slog.String("error", err.Error()))
			return nil, err
		}
		subordinates = append(subordinates, subordinate)
	}
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Found %d subordinates for member", len(subordinates)))
	return subordinates, nil
}

func (p Provider) UpdateMember(m types.Member) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating member", slog.Any("member", m))
	var res sql.Result
	var err error
	if m.SupervisorID == "" {
		p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Supervisory ID is empty, inserting as null in database")
		res, err = p.Db.Exec(updateMemberQuery, m.FirstName, m.LastName, m.Rank, nil, m.Admin, m.Hash, m.ID)
	} else {
		res, err = p.Db.Exec(updateMemberQuery, m.FirstName, m.LastName, m.Rank, m.SupervisorID, m.Admin, m.Hash, m.ID)
	}
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Attempting to update member with non-existent supervisor")
		return backend.ErrSupervisorNotFound
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating member", slog.String("error", err.Error()))
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated for member, got 0")
		return backend.ErrMemberNotFound
	}
	return nil
}

func (p Provider) DeleteMember(identifier string, method backend.ProviderMethod) error {
	p.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member", slog.String("identifier", identifier))
	var res sql.Result
	var err error
	switch method {
	case backend.ById:
		res, err = p.Db.Exec(deleteMemberQuery, identifier)
	case backend.ByUsername:
		res, err = p.Db.Exec(deleteMemberByUsernameQuery, identifier)
	}
	if err != nil {
		p.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting member", slog.String("error", err.Error()))
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		p.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated for member, got 0")
		return backend.ErrMemberNotFound
	}
	return nil
}
