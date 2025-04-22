package memberprovider

const (
	insertMemberQuery           = "INSERT INTO member(id, first_name, last_name, rank, user_name, supervisor_id, admin, hash, disabled) VALUES(?, ?, ?, ?, ?, ?, ?, ?, 00);"
	getMemberQuery              = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash, disabled FROM member WHERE id=?;"
	getMemberByUsernameQuery    = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash, disabled FROM member WHERE user_name=?;"
	getAllMembersQuery          = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash FROM member WHERE disabled != 1;"
	getDisabledMembersQuery     = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash FROM member WHERE disabled=1;"
	getSubordinatesQuery        = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash FROM member WHERE supervisor_id=? AND disabled <> 1;"
	removeSubordinatesQuery     = "UPDATE member SET supervisor_id=null WHERE supervisor_id=?;"
	updateMemberQuery           = "UPDATE member SET user_name=?, first_name=?, last_name=?, rank=?, supervisor_id=?, admin=?, hash=? WHERE ID=?;"
	deleteMemberQuery           = "DELETE FROM member WHERE id=?;"
	deleteMemberByUsernameQuery = "DELETE FROM member WHERE user_name=?;"
	disableMemberQuery          = "UPDATE member SET disabled=1 WHERE id=?;"
	enableMemberQuery           = "UPDATE member set disabled=0 WHERE id=?;"
	getMemberFromSessionQuery   = "SELECT member_id FROM member_session WHERE session_id=?;"
)
