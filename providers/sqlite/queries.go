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
    FOREIGN KEY (supervisor_id) REFERENCES member(id) ON DELETE SET NULL
);

CREATE TABLE qualification(
    id string PRIMARY KEY,
    name string UNIQUE,
    notes string,
    expires integer,
    expiration_days integer
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
    name string UNIQUE,
    description string,
    notes string,
    days_valid_for integer,
    reference_id string,
    FOREIGN KEY (reference_id) REFERENCES reference(id) ON DELETE SET NULL 
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

CREATE TABLE session(
    id string PRIMARY KEY,
    expiration datetime,
    user_agent string
);

CREATE TABLE member_session(
    member_id string,
    session_id string,
    FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
    PRIMARY KEY (member_id, session_id)
);

CREATE TABLE reference(
    id string PRIMARY KEY,
    name string UNIQUE,
    volume int,
    paragraph string
);

INSERT INTO versions VALUES(1);`

	insertMemberQuery           = "INSERT INTO member(id, first_name, last_name, rank, user_name, supervisor_id, admin, hash) VALUES($1, $2, $3, $4, $5, $6, $7, $8);"
	getMemberQuery              = "SELECT * FROM member WHERE id=$1;"
	getMemberByUsernameQuery    = "SELECT * FROM member where user_name=$1;"
	getAllMembersQuery          = "SELECT * FROM member;"
	getSubordinatesQuery        = "SELECT * FROM member WHERE supervisor_id=$1;"
	updateMemberQuery           = "UPDATE member SET first_name=$1, last_name=$2, rank=$3, supervisor_id=$4, admin=$5, hash=$6 WHERE ID=$7;"
	deleteMemberQuery           = "DELETE FROM member WHERE id=$1;"
	deleteMemberByUsernameQuery = "DELETE FROM member WHERE user_name=$1;"

	insertQualificationQuery                     = "INSERT INTO qualification(id, name, notes, expires, expiration_days) VALUES($1, $2, $3, $4, $5);"
	getQualificationQuery                        = "SELECT * FROM qualification WHERE id=$1;"
	getAllQualificationIDsQuery                  = "SELECT id FROM qualification;"
	updateQualificationQuery                     = "UPDATE qualification SET name=$1, notes=$2, expires=$3, expiration_days=$4 WHERE ID=$6;"
	deleteQualificationQuery                     = "DELETE FROM qualification WHERE id=$1;"
	insertQualificationInitialRequirementQuery   = "INSERT INTO qualification_initial_requirement(qualification_id, requirement_id) VALUES($1, $2);"
	insertQualificationRecurringRequirementQuery = "INSERT INTO qualification_recurring_requirement(qualification_id, requirement_id) VALUES($1, $2);"
	getInitialRequirementIdsQuery                = "SELECT requirement_id FROM qualification_initial_requirement WHERE qualification_id=$1;"
	getRecurringRequirementIdsQuery              = "SELECT requirement_id FROM qualification_recurring_requirement WHERE qualification_id=$1;"
	deleteQualificationRecurringRequirementQuery = "DELETE FROM qualification_recurring_requirement WHERE requirement_id=$1;"
	deleteQualificationInitialRequirementQuery   = "DELETE FROM qualification_initial_requirement WHERE requirement_id=$1;"

	addMemberQualificationQuery    = "INSERT INTO member_qualification(member_id, qualification_id) VALUES($1, $2);"
	checkMemberQualificationQuery  = "SELECT COUNT(*) FROM member_qualification WHERE member_id=$1 AND qualification_id=$2;"
	getMemberQualificationIDsQuery = "SELECT qualification_id FROM member_qualification WHERE member_id=$1;"
	removeMemberQualificationQuery = "DELETE FROM member_qualification WHERE member_id=$1 AND qualification_ID=$2;"

	addRequirementQuery                  = "INSERT INTO requirement(id, name, description, notes, days_valid_for, reference_id) VALUES($1, $2, $3, $4, $5, $6);"
	getRequirementQuery                  = "SELECT * FROM requirement r FULL JOIN reference re ON r.reference_id = re.id WHERE r.id = $1;"
	getAllRequirementsQuery              = "SELECT * FROM requirement r FULL JOIN reference re ON r.reference_id = re.id;"
	getQualificationsForRequirementQuery = "SELECT qualification_id FROM qualification_initial_requirement  WHERE requirement_id=$1 UNION SELECT qualification_id FROM qualification_recurring_requirement WHERE requirement_id=$1;"
	updateRequirementQuery               = "UPDATE requirement SET name=$1, description=$2, notes=$3, days_valid_for=$4, reference_id=$5 WHERE id=$6;"
	deleteRequirementQuery               = "DELETE FROM requirement WHERE id=$1;"

	addReferenceQuery    = "INSERT INTO reference(id, name, volume, paragraph) VALUES($1, $2, $3, $4);"
	getReferenceQuery    = "SELECT * FROM reference WHERE id=$1;"
	getReferencesQuery   = "SELECT * FROM reference;"
	updateReferenceQuery = "UPDATE reference SET name=$1, volume=$2, paragraph=$3 WHERE id=$4;"
	deleteReferenceQuery = "DELETE FROM reference WHERE id=$1;"

	insertSessionQuery       = "INSERT INTO session(id, expiration, user_agent) VALUES($1, $2, $3);"
	insertMemberSessionQuery = "INSERT INTO member_session(member_id, session_id) VALUES($1, $2);"
	getSessionQuery          = "SELECT * FROM session WHERE id=$1;"
	deleteSessionQuery       = "DELETE FROM session WHERE id=$1;"
	getMemberSessionQuery    = "SELECT * FROM member_session WHERE member_id=$1 AND session_id=$2;"
)
