package db

import (
	"fmt"
	"time"
)

func EncounterAdd(
	importID int64,
	timestampStart time.Time,
	encounterID int64,
	encounterName string,
	groupSize int64,
	difficulty int64,
) error {
	var err error

	_, err = db.Exec(`
		insert into encounter (import_id, timestamp_start, encounter_id, encounter_name, group_size, difficulty)
		values (?, ?, ?, ?, ?, ?)`,
		importID,
		timestampStart.UnixNano(),
		encounterID,
		encounterName,
		groupSize,
		difficulty,
	)

	return err
}

func EncounterEnd(
	importID int64,
	timestampEnd time.Time,
	encounterID int64,
	success bool,
) error {
	var rowID int64
	var timeStart int64
	var oldEncounterID int64
	var err error

	err = db.QueryRow(`
		select id, timestamp_start, encounter_id
		from encounter
		where import_id = ?
		order by timestamp_start desc
		limit 1`,
		importID,
	).Scan(&rowID, &timeStart, &oldEncounterID)
	if err != nil {
		return err
	}

	if oldEncounterID != encounterID {
		return fmt.Errorf("encounter id mismatch: old %d new %d", oldEncounterID, encounterID)
	}

	if timeStart > timestampEnd.UnixNano() {
		return fmt.Errorf("end before begin: old %d new %d", timeStart, timestampEnd.UnixNano())
	}

	_, err = db.Exec(`
		update encounter set timestamp_end = ?, success = ?
		where id = ?`,
		timestampEnd.UnixNano(),
		success,
		rowID,
	)

	return err
}
