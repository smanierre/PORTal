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
    date_assigned datetime NOT NULL,
    assigned_by_id string NOT NULL,
	PRIMARY KEY (member_id, qualification_id),
	FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
	FOREIGN KEY (qualification_id) REFERENCES qualification(id) ON DELETE CASCADE,
	FOREIGN KEY (assigned_by_id) REFERENCES member(id) ON DELETE NO ACTION 
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

CREATE TABLE initial_member_requirement(
    member_id string,
    requirement_id string,
    completed_date datetime,
    assigned_by string not null,
    completed_by string,
    PRIMARY KEY (member_id, requirement_id),
    FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
    FOREIGN KEY (requirement_id) REFERENCES requirement(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES member(id) ON DELETE NO ACTION,
    FOREIGN KEY (completed_by) REFERENCES member(id) ON DELETE NO ACTION
);

CREATE TABLE recurring_member_requirement(
    id string PRIMARY KEY,
    member_id string,
    requirement_id string,
    assigned_by string,
    completed_by string not null,
    completed_date datetime,
    FOREIGN KEY (member_id) REFERENCES member(id) ON DELETE CASCADE,
    FOREIGN KEY (requirement_id) REFERENCES requirement(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES member(id) ON DELETE NO ACTION
);

CREATE TABLE recurring_member_requirement_completion(
    id string PRIMARY KEY,
    recurring_member_requirement_id string,
    completed_date datetime not null,
    completed_by string not null,
    FOREIGN KEY (recurring_member_requirement_id) REFERENCES recurring_member_requirement(id) ON DELETE CASCADE,
    FOREIGN KEY (completed_by) REFERENCES member(id) ON DELETE NO ACTION
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
