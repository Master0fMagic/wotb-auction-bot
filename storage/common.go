package storage

const (
	createTableQuery = `
		CREATE TABLE IF NOT EXISTS monitoring_data (
			vehicle_name TEXT,
			chat_id INTEGER,
			minimal_count INTEGER,
			PRIMARY KEY (vehicle_name, chat_id)
		);
	`

	insertDataQuery = `
		INSERT OR REPLACE INTO monitoring_data (vehicle_name, chat_id, minimal_count) 
		VALUES (?, ?, ?);
`

	removeDataQuery = `
		DELETE FROM monitoring_data 
		WHERE vehicle_name = ? AND chat_id = ?;
	`

	getAllQuery = `
		SELECT vehicle_name, chat_id, minimal_count 
		FROM monitoring_data;
	`

	getAllByVehicleAndCountGteQuery = `
		SELECT vehicle_name, chat_id, minimal_count 
		FROM monitoring_data 
		WHERE vehicle_name = ? AND minimal_count >= ?;
	`
)
