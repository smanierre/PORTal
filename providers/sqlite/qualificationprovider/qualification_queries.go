package qualificationprovider

const (
	insertQualificationQuery                     = "INSERT INTO qualification(id, name, notes) VALUES(?, ?, ?);"
	getQualificationQuery                        = "SELECT * FROM qualification WHERE id=?;"
	getAllQualificationIDsQuery                  = "SELECT id FROM qualification;"
	updateQualificationQuery                     = "UPDATE qualification SET name=?, notes=? WHERE ID=?;"
	deleteQualificationQuery                     = "DELETE FROM qualification WHERE id=?;"
	insertQualificationInitialRequirementQuery   = "INSERT INTO qualification_initial_requirement(qualification_id, requirement_id) VALUES(?, ?);"
	insertQualificationRecurringRequirementQuery = "INSERT INTO qualification_recurring_requirement(qualification_id, requirement_id) VALUES(?, ?);"
	getInitialRequirementIdsQuery                = "SELECT requirement_id FROM qualification_initial_requirement WHERE qualification_id=?;"
	getRecurringRequirementIdsQuery              = "SELECT requirement_id FROM qualification_recurring_requirement WHERE qualification_id=?;"
	deleteQualificationRecurringRequirementQuery = "DELETE FROM qualification_recurring_requirement WHERE requirement_id=?;"
	deleteQualificationInitialRequirementQuery   = "DELETE FROM qualification_initial_requirement WHERE requirement_id=?;"
	checkMemberQualificationQuery                = "SELECT COUNT(*) FROM member_qualification WHERE member_id=? AND qualification_id=?;"
	getMemberQualificationIDsQuery               = "SELECT qualification_id FROM member_qualification WHERE member_id=?;"
)
