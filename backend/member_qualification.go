package backend

import (
	"PORTal/types"
	"context"
	"log/slog"
)

func (b Backend) AssignMemberQualification(memberID, qualificationID string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding qualification to member",
		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
	return b.memberProvider.AssignMemberQualification(memberID, qualificationID)
}

func (b Backend) GetMemberQualification(memberID, qualificationID string) (types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualification for member",
		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
	return b.memberProvider.GetMemberQualification(memberID, qualificationID)
}

func (b Backend) GetMemberQualifications(memberID string) ([]types.Qualification, error) {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualifications assigned to member",
		slog.String("member_id", memberID))
	return b.memberProvider.GetMemberQualifications(memberID)
}

func (b Backend) RemoveMemberQualification(memberID, qualificationID string) error {
	b.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member qualification",
		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
	return b.memberProvider.RemoveMemberQualification(memberID, qualificationID)
}
