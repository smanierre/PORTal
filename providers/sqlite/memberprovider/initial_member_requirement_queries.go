package memberprovider

const (
	addInitialMemberRequirementQuery  = `INSERT INTO initial_member_requirement(member_id, requirement_id, completed_date, assigned_by, completed_by) VALUES(?, ?, ?, ?, ?);`
	getInitialMemberRequirementQuery  = `SELECT member_id, requirement_id, completed_date, assigned_by, completed_by FROM initial_member_requirement WHERE member_id=? AND requirement_id=?;`
	getInitialMemberRequirementsQuery = `
		SELECT i.member_id, i.requirement_id, i.completed_date, i.assigned_by, i.completed_by 
		FROM initial_member_requirement i INNER JOIN qualification_initial_requirement q 
    	ON i.requirement_id = q.requirement_id 
        AND member_id=? 
        AND qualification_id=?;`
	completeInitialMemberRequirementQuery = `UPDATE initial_member_requirement SET completed_date=?, completed_by=? WHERE member_id=? AND requirement_id=?;`
)
