package sessionprovider

const (
	insertSessionQuery       = "INSERT INTO session(id, user_agent, ip_address, expiration) VALUES(?, ?, ?, ?);"
	insertMemberSessionQuery = "INSERT INTO member_session(member_id, session_id) VALUES(?, ?);"
	getSessionQuery          = "SELECT id, expiration, user_agent, ip_address FROM session WHERE id=?;"
	deleteSessionQuery       = "DELETE FROM session WHERE id=?;"
)
