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
    notes string
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
    name string NOT NULL,
    initial integer NOT NULL,
    reference string NOT NULL,
    notes string,
    type string NOT NULL,
    qualification_id string,
    grade string,
    days_valid_for integer
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
)
