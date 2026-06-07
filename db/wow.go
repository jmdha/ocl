package db

import "time"

func InsertWOWEventVersion(
	import_id int64,
	line int64,
	timestamp time.Time,
	log int64,
	advanced int64,
	major int64,
	minor int64,
	patch int64,
	project int64,
) error {
	_, err := db.Exec(`
		insert into wow_event_version (
			import_id,
			line,
			timestamp,
			log,
			advanced,
			major,
			minor,
			patch,
			project
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		import_id,
		line,
		timestamp,
		log,
		advanced,
		major,
		minor,
		patch,
		project,
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertWOWEncounterStart(
	import_id int64,
	line int64,
	timestamp time.Time,
	encounterID int64,
	encounterName string,
	difficultyID int64,
	groupSize int64,
	instanceID int64,
) error {
	_, err := db.Exec(`
		insert into wow_event_encounter_start (
			import_id        ,
			line             ,
			timestamp        ,
			encounter_id     ,
			encounter_name   ,
			difficulty_id    ,
			group_size       ,
			instance_id      
		)
		values (?, ?, ?, ?, ?, ?, ?, ?)`,
		import_id,
		line,
		timestamp,
		encounterID,
		encounterName,
		difficultyID,
		groupSize,
		instanceID,
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertWOWEncounterEnd(
	import_id int64,
	line int64,
	timestamp time.Time,
	encounterID int64,
	encounterName string,
	difficultyID int64,
	groupSize int64,
	success bool,
	duration int64,
) error {
	_, err := db.Exec(`
		insert into wow_event_encounter_end (
			import_id        ,
			line             ,
			timestamp        ,
			encounter_id     ,
			encounter_name   ,
			difficulty_id    ,
			group_size       ,
			success          ,
			duration         
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		import_id,
		line,
		timestamp,
		encounterID,
		encounterName,
		difficultyID,
		groupSize,
		success,
		duration,
	)
	if err != nil {
		return err
	}

	return nil
}
