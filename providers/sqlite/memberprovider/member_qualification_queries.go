package memberprovider

const (
	addMemberQualificationQuery    = "INSERT INTO member_qualification(member_id, qualification_id) VALUES(?, ?);"
	removeMemberQualificationQuery = "DELETE FROM member_qualification WHERE member_id=? AND qualification_ID=?;"
)
