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
    admin integer,
    hash string,
    disabled integer,
    FOREIGN KEY (supervisor_id) REFERENCES member(id) ON DELETE SET NULL
);

CREATE TABLE qualification(
    id string PRIMARY KEY,
    name string UNIQUE,
    notes string,
    expires integer,
    expiration_interval integer
);

CREATE TABLE member_qualification(
    member_id string,
    qualification_id string,
	PRIMARY KEY (member_id, qualification_id),
	FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
	FOREIGN KEY (qualification_id) REFERENCES qualification(id) ON DELETE CASCADE
);

CREATE TABLE requirement(
    id string PRIMARY KEY,
    name string,
    notes string,
    grade string,
    days_valid_for integer,
    reference string,
    type string NOT NULL
);

CREATE TABLE member_requirement(
    member_id string,
    requirement_id string,
    initial_completion datetime,
    most_recent_completion datetime,
    PRIMARY KEY (member_id, requirement_id),
    FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
    FOREIGN KEY (requirement_id) REFERENCES requirement(id) ON DELETE CASCADE
);

CREATE TABLE qualification_initial_requirement(
    qualification_id string,
    requirement_id string,
    PRIMARY KEY (qualification_id, requirement_id),
    FOREIGN KEY (qualification_id) REFERENCES qualification(id) ON DELETE CASCADE,
    FOREIGN KEY (requirement_id) REFERENCES requirement(id) ON DELETE CASCADE
);

CREATE TABLE qualification_recurring_requirement(
    qualification_id string,
    requirement_id string,
    PRIMARY KEY (qualification_id, requirement_id),
    FOREIGN KEY (qualification_id) REFERENCES qualification(id) ON DELETE CASCADE,
    FOREIGN KEY (requirement_id) REFERENCES requirement(id) ON DELETE CASCADE
);

CREATE TABLE qualification_requirement(
    requirement_id string,
    qualification_id string,
    PRIMARY KEY (requirement_id, qualification_id),
    FOREIGN KEY (requirement_id) REFERENCES requirement(id) ON DELETE CASCADE,
    FOREIGN KEY (qualification_id) REFERENCES qualification(id) ON DELETE CASCADE
);

CREATE TABLE session(
    id string PRIMARY KEY,
    user_agent string,
    ip_address string,
    expiration datetime
);

CREATE TABLE member_session(
    member_id string,
    session_id string,
    FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
    PRIMARY KEY (member_id, session_id)
);

INSERT INTO versions VALUES(1);`

	insertMemberQuery           = "INSERT INTO member(id, first_name, last_name, rank, user_name, supervisor_id, admin, hash, disabled) VALUES(?, ?, ?, ?, ?, ?, ?, ?, 00);"
	getMemberQuery              = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash, disabled FROM member WHERE id=?;"
	getMemberByUsernameQuery    = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash, disabled FROM member WHERE user_name=?;"
	getAllMembersQuery          = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash FROM member WHERE disabled != 1;"
	getDisabledMembersQuery     = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash FROM member WHERE disabled=1;"
	getSubordinatesQuery        = "SELECT id, first_name, last_name, rank, user_name, supervisor_id, admin, hash FROM member WHERE supervisor_id=? AND disabled != 1;"
	removeSubordinatesQuery     = "UPDATE member SET supervisor_id=null WHERE supervisor_id=?;"
	updateMemberQuery           = "UPDATE member SET first_name=?, last_name=?, rank=?, supervisor_id=?, admin=?, hash=? WHERE ID=?;"
	deleteMemberQuery           = "DELETE FROM member WHERE id=?;"
	deleteMemberByUsernameQuery = "DELETE FROM member WHERE user_name=?;"
	disableMemberQuery          = "UPDATE member SET disabled=1 WHERE id=?;"
	enableMemberQuery           = "UPDATE member set disabled=0 WHERE id=?;"

	insertQualificationQuery                     = "INSERT INTO qualification(id, name, notes, expires, expiration_interval) VALUES(?, ?, ?, ?, ?);"
	getQualificationQuery                        = "SELECT * FROM qualification WHERE id=?;"
	getAllQualificationIDsQuery                  = "SELECT id FROM qualification;"
	updateQualificationQuery                     = "UPDATE qualification SET name=?, notes=?, expires=?, expiration_interval=? WHERE ID=?;"
	deleteQualificationQuery                     = "DELETE FROM qualification WHERE id=?;"
	insertQualificationInitialRequirementQuery   = "INSERT INTO qualification_initial_requirement(qualification_id, requirement_id) VALUES(?, ?);"
	insertQualificationRecurringRequirementQuery = "INSERT INTO qualification_recurring_requirement(qualification_id, requirement_id) VALUES(?, ?);"
	getInitialRequirementIdsQuery                = "SELECT requirement_id FROM qualification_initial_requirement WHERE qualification_id=?;"
	getRecurringRequirementIdsQuery              = "SELECT requirement_id FROM qualification_recurring_requirement WHERE qualification_id=?;"
	deleteQualificationRecurringRequirementQuery = "DELETE FROM qualification_recurring_requirement WHERE requirement_id=?;"
	deleteQualificationInitialRequirementQuery   = "DELETE FROM qualification_initial_requirement WHERE requirement_id=?;"

	addMemberQualificationQuery    = "INSERT INTO member_qualification(member_id, qualification_id) VALUES(?, ?);"
	checkMemberQualificationQuery  = "SELECT COUNT(*) FROM member_qualification WHERE member_id=? AND qualification_id=?;"
	getMemberQualificationIDsQuery = "SELECT qualification_id FROM member_qualification WHERE member_id=?;"
	removeMemberQualificationQuery = "DELETE FROM member_qualification WHERE member_id=? AND qualification_ID=?;"

	addRequirementQuery                  = "INSERT INTO requirement(id, name, notes, grade, days_valid_for, reference, type) VALUES(?, ?, ?, ?, ?, ?, ?);"
	getRequirementQuery                  = "SELECT * FROM requirement WHERE id=?;"
	getAllRequirementsQuery              = "SELECT * FROM requirement WHERE id IS NOT NULL;"
	getQualificationsForRequirementQuery = "SELECT qualification_id FROM qualification_initial_requirement WHERE requirement_id=? UNION SELECT qualification_id FROM qualification_recurring_requirement WHERE requirement_id=?;"
	updateRequirementQuery               = "UPDATE requirement SET name=?, notes=?, grade=?, days_valid_for=?, reference=?, type=? WHERE id=?;"
	deleteRequirementQuery               = "DELETE FROM requirement WHERE id=?;"

	insertSessionQuery        = "INSERT INTO session(id, user_agent, ip_address, expiration) VALUES(?, ?, ?, ?);"
	insertMemberSessionQuery  = "INSERT INTO member_session(member_id, session_id) VALUES(?, ?);"
	getSessionQuery           = "SELECT id, expiration, user_agent, ip_address FROM session WHERE id=?;"
	deleteSessionQuery        = "DELETE FROM session WHERE id=?;"
	getMemberFromSessionQuery = "SELECT member_id FROM member_session WHERE session_id=?;"
)
