package memberprovider

const (
	addMemberQualificationQuery     = "INSERT INTO member_qualification(member_id, qualification_id, date_assigned, assigned_by_id) VALUES(?, ?, ?, ?);"
	removeMemberQualificationQuery  = "DELETE FROM member_qualification WHERE member_id=? AND qualification_ID=?;"
	getMemberQualificationQuery     = "SELECT member_id, qualification_id, date_assigned, assigned_by_id FROM member_qualification WHERE member_id=? AND qualification_id=?;"
	getMemberQualificationsQuery    = "SELECT member_id, qualification_id, date_assigned, assigned_by_id FROM member_qualification;"
	getQualificationsForMemberQuery = "SELECT member_id, qualification_id, date_assigned, assigned_by_id FROM member_qualification WHERE member_id=?;"
)
