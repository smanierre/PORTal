package qualificationprovider

const (
	addRequirementQuery                  = "INSERT INTO requirement(id, name, initial, reference, notes, type, qualification_id, grade, days_valid_for) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);"
	getRequirementQuery                  = "SELECT id, name, initial, reference, notes, type, qualification_id, grade, days_valid_for FROM requirement WHERE id=?;"
	getAllRequirementsQuery              = "SELECT id, name, initial, reference, notes, type, qualification_id, grade, days_valid_for FROM requirement WHERE id IS NOT NULL;"
	getQualificationsForRequirementQuery = "SELECT qualification_id FROM qualification_initial_requirement WHERE requirement_id=? UNION SELECT qualification_id FROM qualification_recurring_requirement WHERE requirement_id=?;"
	updateRequirementQuery               = "UPDATE requirement SET name=?, initial=?, reference=?, notes=?, type=?, qualification_id=?, grade=?, days_valid_for=? WHERE id=?;"
	deleteRequirementQuery               = "DELETE FROM requirement WHERE id=?;"
)
