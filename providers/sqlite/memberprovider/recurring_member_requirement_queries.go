package memberprovider

const (
	addRecurringMemberRequirementsQuery = `INSERT INTO recurring_member_requirement(id, member_id, requirement_id, assigned_by) VALUES(?, ?, ?, ?);`

	getRecurringMemberRequirementQuery = `
		SELECT id, member_id, requirement_id, assigned_by, completed_by, completed_date
		FROM recurring_member_requirement
		WHERE id=?;`

	getRecurringMemberRequirementsQuery = `
		SELECT id, member_id, requirement_id, assigned_by, completed_by, completed_date
		FROM recurring_member_requirement
		WHERE member_id=?;
	`
	getRecurringMemberRequirementsForQualificationQuery = `
		SELECT r.id, r.member_id, r.requirement_id, r.assigned_by, r.completed_by, r.completed_date
		FROM recurring_member_requirement r 
		INNER JOIN qualification_recurring_requirement q 
		ON r.requirement_id=q.requirement_id 
		AND r.member_id=? AND q.qualification_id=?`

	getRecurringMemberRequirementCompletionQuery = `
		SELECT id, completed_date, completed_by, completed_by, completed_date
		FROM recurring_member_requirement_completion 
		WHERE recurring_member_requirement_id=?;`

	completeRecurringMemberRequirementQuery = `
		INSERT INTO recurring_member_requirement_completion(id, recurring_member_requirement_id, completed_date, completed_by, completed_date)
		VALUES(?, ?, ?, ?, ?);
	`
)
