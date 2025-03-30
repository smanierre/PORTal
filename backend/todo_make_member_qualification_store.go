package backend

//func (s Backend) AssignMemberQualification(memberID, qualificationID string) error {
//	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Adding qualification to member",
//		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
//	return s.memberProvider.AssignMemberQualification(memberID, qualificationID)
//}
//
//func (s Backend) GetMemberQualification(memberID, qualificationID string) (types.Qualification, error) {
//	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualification for member",
//		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
//	return s.memberProvider.GetMemberQualification(memberID, qualificationID)
//}
//
//func (s Backend) GetMemberQualifications(memberID string) ([]types.Qualification, error) {
//	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Getting qualifications assigned to member",
//		slog.String("member_id", memberID))
//	return s.memberProvider.GetMemberQualifications(memberID)
//}
//
//func (s Backend) RemoveMemberQualification(memberID, qualificationID string) error {
//	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Deleting member qualification",
//		slog.String("member_id", memberID), slog.String("qualification_id", qualificationID))
//	return s.memberProvider.RemoveMemberQualification(memberID, qualificationID)
//}
