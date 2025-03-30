package qualificationprovider

const (
	addRequirementQuery                  = "INSERT INTO requirement(id, name, notes, grade, days_valid_for, reference, type) VALUES(?, ?, ?, ?, ?, ?, ?);"
	getRequirementQuery                  = "SELECT * FROM requirement WHERE id=?;"
	getAllRequirementsQuery              = "SELECT * FROM requirement WHERE id IS NOT NULL;"
	getQualificationsForRequirementQuery = "SELECT qualification_id FROM qualification_initial_requirement WHERE requirement_id=? UNION SELECT qualification_id FROM qualification_recurring_requirement WHERE requirement_id=?;"
	updateRequirementQuery               = "UPDATE requirement SET name=?, notes=?, grade=?, days_valid_for=?, reference=?, type=? WHERE id=?;"
	deleteRequirementQuery               = "DELETE FROM requirement WHERE id=?;"
)
