package memberqualificationstore

import (
	"PORTal/types"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

func (m MemberQualificationStore) AssignMemberQualification(memberID, qualificationID, assignedByID string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding qualification to member",
		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
	err := m.memberProvider.AssignMemberQualification(memberID, qualificationID, m.clock.Now(), assignedByID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error assigning qualification to member", slog.String("error", err.Error()))
		return err
	}
	qual, err := m.qualificationProvider.GetQualification(qualificationID)
	if err != nil {
		m.logger.LogAttrs(context.Background(), slog.LevelError, "Error retrieving qualification from DB", slog.String("error", err.Error()))
		return err
	}
	for _, req := range qual.InitialRequirements {
		err = m.memberRequirementProvider.AssignInitialMemberRequirement(memberID, req.ID, assignedByID)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error assigning initial member requirement", slog.String("error", err.Error()))
			return err
		}
	}
	for _, req := range qual.RecurringRequirements {
		err = m.memberRequirementProvider.AssignRecurringMemberRequirement(uuid.NewString(), memberID, req.ID, assignedByID)
		if err != nil {
			m.logger.LogAttrs(context.Background(), slog.LevelError, "Error assigning recurring member requirement", slog.String("error", err.Error()))
			return err
		}
	}
	return nil
}

func (m MemberQualificationStore) RemoveMemberQualification(memberID, qualificationID string) error {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member qualification",
		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
	return m.memberProvider.RemoveMemberQualification(memberID, qualificationID)
}

func (m MemberQualificationStore) GetMemberQualification(memberID, qualificationID string) (types.MemberQualification, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting member qualification",
		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
	return m.memberProvider.GetMemberQualification(memberID, qualificationID)
}

func (m MemberQualificationStore) GetAllMemberQualifications() ([]types.MemberQualification, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting all member qualifications")
	return m.memberProvider.GetAllMemberQualifications()
}

func (m MemberQualificationStore) GetQualificationsForMember(memberID string) ([]types.MemberQualification, error) {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualifications assigned to member", slog.String("member_id", memberID))
	return m.memberProvider.GetQualificationsForMember(memberID)
}
