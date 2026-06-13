package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"ocl/db"
	"ocl/wcl"
	"time"

	"github.com/klauspost/compress/zstd"
)

func worker(id int) {
	var job int64
	var data []byte
	var err error
	for {
		job, data, err = db.ImportClaim()
		if errors.Is(err, sql.ErrNoRows) {
			time.Sleep(time.Second)
			continue
		}

		log.Printf("[worker %d] processing %d", id, job)
		err = process(job, data)
		if err != nil {
			log.Printf("[worker %d] failed to process job %d: %v", id, job, err)
			err = db.ImportSetFailed(job, err.Error())
			if err != nil {
				log.Printf("[worker %d] failed to error %d: %v", id, job, err)
			}
			continue
		}

		err = db.ImportSetDone(job)
		if err != nil {
			log.Printf("[worker %d] failed to mark job done %d: %v", id, job, err)
		}
		log.Printf("[worker %d] finished processing %d", id, job)
	}
}

func process(job int64, data []byte) error {
	byteReader := bytes.NewReader(data)
	zstdReader, err := zstd.NewReader(byteReader)
	if err != nil {
		return err
	}
	defer zstdReader.Close()

	var event wcl.Event
	scanner := bufio.NewScanner(zstdReader)
	for scanner.Scan() {
		err = wcl.Parse(&event, scanner.Text())
		if err != nil {
			return fmt.Errorf("parse error: %v", err)
		}

		err = process_event(job, event)
		if err != nil {
			return fmt.Errorf("event process error: %v", err)
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("scanner error: %v", scanner.Err())
	}

	return nil
}

func process_event(job int64, event wcl.Event) error {
	switch event.Kind {
	case wcl.KindSpellDamage:
		sourceID, err := db.CharacterID(event.SpellDamage.SourceGUID, event.SpellDamage.SourceName)
		if err != nil {
			return err
		}

		targetID, err := db.CharacterID(event.SpellDamage.TargetGUID, event.SpellDamage.TargetName)
		if err != nil {
			return err
		}

		return db.InsertEventDamage(
			job,
			event.SpellDamage.Timestamp,
			sourceID,
			targetID,
			0,
			event.SpellDamage.Damage,
		)
	case wcl.KindEncounterStart:
		return db.EncounterAdd(
			job,
			event.EncounterStart.Timestamp,
			event.EncounterStart.EncounterID,
			event.EncounterStart.EncounterName,
			event.EncounterStart.GroupSize,
			event.EncounterStart.DifficultyID,
		)
	case wcl.KindEncounterEnd:
		return db.EncounterEnd(
			job,
			event.EncounterEnd.Timestamp,
			event.EncounterEnd.EncounterID,
			event.EncounterEnd.Success,
		)
	}
	return nil
}
