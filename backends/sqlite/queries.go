package sqlite

const (
	createStructureQuery = `CREATE TABLE versions(version float PRIMARY KEY);
CREATE TABLE member(
    id string PRIMARY KEY,
     first_name string,
     last_name string,
     rank string,
     user_name string UNIQUE,
     supervisor_id string,
     hash string,
     FOREIGN KEY (supervisor_id) REFERENCES member(id) ON DELETE SET NULL
);
CREATE TABLE qualification(
    id string PRIMARY KEY,
    name string,
    notes string,
    expires integer,
    expiration_days integer
);

CREATE TABLE member_qualification(
    member_id string,
    qualification_id string,
    active integer,
    active_time datetime,
	PRIMARY KEY (member_id, qualification_id),
	FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
	FOREIGN KEY (qualification_id) REFERENCES qualification(id) ON DELETE CASCADE
);

CREATE TABLE requirement(
    id string PRIMARY KEY,
    name string,
    description string,
    notes string,
    days_valid_for integer
);

CREATE TABLE session(
    id string PRIMARY KEY,
    expiration datetime,
    ip_address string
);

CREATE TABLE member_session(
    member_id string,
    session_id string,
    FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
    PRIMARY KEY (member_id, session_id)
);

INSERT INTO versions VALUES(1);`

	insertMemberQuery        = "INSERT INTO member(id, first_name, last_name, rank, user_name, supervisor_id, hash) VALUES($1, $2, $3, $4, $5, $6, $7);"
	getMemberQuery           = "SELECT * FROM member WHERE id=$1;"
	getMemberByUsernameQuery = "SELECT * FROM member where user_name=$1;"
	getAllMembersQuery       = "SELECT * FROM member;"
	updateMemberQuery        = "UPDATE member SET first_name=$1, last_name=$2, rank=$3, supervisor_id=$4, hash=$5 WHERE ID=$6;"
	deleteMemberQuery        = "DELETE FROM member WHERE id=$1;"

	insertQualificationQuery  = "INSERT INTO qualification(id, name, notes, expires, expiration_days) VALUES($1, $2, $3, $4, $5);"
	getQualificationQuery     = "SELECT * FROM qualification WHERE id=$1;"
	getAllQualificationsQuery = "SELECT * FROM qualification;"
	updateQualificationQuery  = "UPDATE qualification SET name=$1, notes=$2, expires=$3, expiration_days=$4 WHERE ID=$6;"
	deleteQualificationQuery  = "DELETE FROM qualification WHERE id=$1;"

	addMemberQualificationQuery    = "INSERT INTO member_qualification(member_id, qualification_id, active, active_time) VALUES($1, $2, $3, $4);"
	getMemberQualificationsQuery   = "SELECT * FROM member_qualification WHERE member_id=$1;"
	updateMemberQualificationQuery = "UPDATE member_qualification SET active=$1, active_time=$2 WHERE member_id=$3 AND qualification_id=$4;"
	deleteMemberQualificationQuery = "DELETE FROM member_qualification WHERE member_id=$1 AND qualification_ID=$2;"

	addRequirementQuery     = "INSERT INTO requirement(id, name, description, notes, days_valid_for) VALUES($1, $2, $3, $4, $5);"
	getRequirementQuery     = "SELECT * FROM requirement WHERE id=$1;"
	getAllRequirementsQuery = "SELECT * FROM requirement;"
	updateRequirementQuery  = "UPDATE requirement SET name=$1, description=$2, notes=$3, days_valid_for=$4 WHERE id=$5;"
	deleteRequirementQuery  = "DELETE FROM requirement WHERE id=$1;"

	insertSessionQuery       = "INSERT INTO session(id, expiration, ip_address) VALUES($1, $2, $3);"
	insertMemberSessionQuery = "INSERT INTO member_session(member_id, session_id) VALUES($1, $2);"
	getSessionForMemberQuery = "SELECT * FROM session WHERE id=(SELECT session_id FROM member_session WHERE session_id=$1 AND member_id=$2);"
)
