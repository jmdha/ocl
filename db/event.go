package db

import "time"

func InsertEventDamage(
	import_id int64,
	timestamp time.Time,
	source_id int64,
	target_id int64,
	spell_id int64,
	amount int64,
) error {
	var err error

	_, err = db.Exec(`
		insert into event_damage (
			import_id,
			timestamp,
			source_id,
			target_id,
			spell_id,
			amount
		) values (?, ?, ?, ?, ?, ?)`,
		import_id,
		timestamp.UnixNano(),
		source_id,
		target_id,
		spell_id,
		amount,
	)

	return err
}
