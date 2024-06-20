package sqlite

import (
	"PORTal/types"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

func (b *Backend) AddMember(m types.Member) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Checking for missing fields")
	if err := m.CheckForMissingArgs(); err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Required arguments missing for user creation", slog.String("error", err.Error()))
		return err
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Inserting member into database", slog.Any("member", m))
	var err error
	if m.SupervisorID == "" {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "No supervisor provided, setting to null in database")
		_, err = b.Db.Exec(insertMemberQuery, m.ID, m.FirstName, m.LastName, m.Rank, nil, m.Hash)
	} else {
		_, err = b.Db.Exec(insertMemberQuery, m.ID, m.FirstName, m.LastName, m.Rank, m.SupervisorID, m.Hash)
	}
	if err != nil && strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Provided supervisor id doesn't exist", slog.String("supervisor_id", m.ID))
		return types.ErrSupervisorNotFound
	} else if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error inserting member into database", slog.String("error", err.Error()))
	}
	return err
}

func (b *Backend) GetMember(id string) (types.Member, error) {
	l := b.logger.With(slog.String("id", id))
	l.LogAttrs(context.Background(), slog.LevelInfo, "Getting member from database")
	row := b.Db.QueryRow(getMemberQuery, id)
	var m types.Member
	supervisorId := sql.NullString{}
	err := row.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Rank, &supervisorId, &m.Hash)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		l.LogAttrs(context.Background(), slog.LevelWarn, "No user found with given id")
		return types.Member{}, types.ErrMemberNotFound
	} else if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error scanning member into struct", slog.String("error", err.Error()))
		return types.Member{}, err
	}
	if supervisorId.Valid {
		m.SupervisorID = supervisorId.String
	}
	//TODO: implement getting qualifications for member
	return m, nil
}

func (b *Backend) GetAllMembers() ([]types.Member, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all members from database")
	rows, err := b.Db.Query(getAllMembersQuery)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error getting all members from database", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()
	var members []types.Member
	for rows.Next() {
		m := types.Member{}
		supervisorId := sql.NullString{}
		err = rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Rank, &supervisorId, &m.Hash)
		if err != nil {
			b.logger.LogAttrs(context.Background(), slog.LevelError, "Error scanning member into struct", slog.String("error", err.Error()))
			continue
		}
		if supervisorId.Valid {
			m.SupervisorID = supervisorId.String
		}
		members = append(members, m)
	}
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Retreived %d members from database", len(members)))
	return members, nil
}

func (b *Backend) UpdateMember(m types.Member) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Updating member", slog.Any("member", m))
	var res sql.Result
	var err error
	if m.SupervisorID == "" {
		b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Supervisory ID is empty, inserting as null in database")
		res, err = b.Db.Exec(updateMemberQuery, m.FirstName, m.LastName, m.Rank, nil, m.Hash, m.ID)
	} else {
		res, err = b.Db.Exec(updateMemberQuery, m.FirstName, m.LastName, m.Rank, m.SupervisorID, m.Hash, m.ID)
	}
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error updating member", slog.String("error", err.Error()))
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated for member, got 0")
		return types.ErrMemberNotFound
	}
	return nil
}

func (b *Backend) DeleteMember(id string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member", slog.String("id", id))
	res, err := b.Db.Exec(deleteMemberQuery, id)
	if err != nil {
		b.logger.LogAttrs(context.Background(), slog.LevelError, "Error deleting member", slog.String("error", err.Error()))
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		b.logger.LogAttrs(context.Background(), slog.LevelWarn, "Expected 1 row to be updated for member, got 0")
		return types.ErrMemberNotFound
	}
	return nil
}
