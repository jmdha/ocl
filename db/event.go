package db

import "time"

func InsertEventEncounterStart(
	import_id int64,
	line int64,
	timestamp time.Time,
	encounter_id int64,
	difficulty_id int64,
	group_size int64,
) error {
	var err error

	_, err = db.Exec(`
		insert into event_encounter_start (
			import_id,
			line,
			timestamp,
			encounter_id,
			difficulty_id,
			group_size
		) values (?, ?, ?, ?, ?, ?)`,
		import_id,
		line,
		timestamp.UnixMilli(),
		encounter_id,
		difficulty_id,
		group_size,
	)

	return err
}

func InsertEventEncounterEnd(
	import_id int64,
	line int64,
	timestamp time.Time,
	encounter_id int64,
	difficulty_id int64,
	group_size int64,
	success bool,
	duration time.Duration,
) error {
	var err error

	_, err = db.Exec(`
		insert into event_encounter_end (
			import_id,
			line,
			timestamp,
			encounter_id,
			difficulty_id,
			group_size,
			success,
			duration
		) values (?, ?, ?, ?, ?, ?, ?, ?)`,
		import_id,
		line,
		timestamp.UnixMilli(),
		encounter_id,
		difficulty_id,
		group_size,
		success,
		duration.Milliseconds(),
	)

	return err
}

func InsertEventDamage(
	import_id int64,
	line int64,
	timestamp time.Time,
	spell_id int64,
	source_id int64,
	target_id int64,
	amount float64,
) error {
	var err error

	_, err = db.Exec(`
		insert into event_damage (
			import_id,
			line,
			timestamp,
			spell_id,
			source_id,
			target_id,
			amount
		) values (?, ?, ?, ?, ?, ?, ?)`,
		import_id,
		line,
		timestamp.UnixMilli(),
		spell_id,
		source_id,
		target_id,
		amount,
	)

	return err
}
