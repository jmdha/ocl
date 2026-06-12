package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"ocl/db"
	"ocl/wcl"
	"os"
	"time"

	"github.com/klauspost/compress/zstd"
)

func worker(id int) {
	var job int64
	var path string
	var err error
	for {
		job, path, err = db.ImportClaim()
		if errors.Is(err, sql.ErrNoRows) {
			time.Sleep(time.Second)
			continue
		}

		log.Printf("[worker %d] processing %d", id, job)
		err = process(job, path)
		if err != nil {
			log.Printf("[worker %d] failed to process job %d: %v", id, job, err)
			err = db.ImportMarkFailed(job, err.Error())
			if err != nil {
				log.Printf("[worker %d] failed to error %d: %v", id, job, err)
			}
			continue
		}

		err = db.ImportMarkDone(job)
		if err != nil {
			log.Printf("[worker %d] failed to mark job done %d: %v", id, job, err)
		}
		log.Printf("[worker %d] finished processing %d", id, job)
	}
}

func process(job int64, path string) error {
	var file *os.File
	var scanner *bufio.Scanner
	var line int64
	var event wcl.Event
	var err error
	file, err = os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer file.Close()

	d, err := zstd.NewReader(file)
	if err != nil {
		return err
	}
	defer d.Close()

	line = 0
	scanner = bufio.NewScanner(d)
	for scanner.Scan() {
		err = wcl.Parse(&event, scanner.Text())
		if err != nil {
			return fmt.Errorf("failed to parse file %v line %d", err, line)
		}

		switch event.Kind {
		case wcl.KindEncounterStart:
			err = db.InsertEventEncounterStart(
				job,
				line,
				event.EncounterStart.Timestamp,
				event.EncounterStart.EncounterID,
				event.EncounterStart.DifficultyID,
				event.EncounterStart.GroupSize,
			)
		case wcl.KindEncounterEnd:
			err = db.InsertEventEncounterEnd(
				job,
				line,
				event.EncounterEnd.Timestamp,
				event.EncounterEnd.EncounterID,
				event.EncounterEnd.DifficultyID,
				event.EncounterEnd.GroupSize,
				event.EncounterEnd.Success,
				event.EncounterEnd.Duration,
			)
		case wcl.KindSpellDamage:
			err = db.InsertEventDamage(
				job,
				line,
				event.SpellDamage.Timestamp,
				0,
				0,
				0,
				event.SpellDamage.Damage,
			)
		case wcl.KindUnknown:
			err = nil
		}
		if err != nil {
			return err
		}
		line++
	}

	if scanner.Err() != nil {
		return fmt.Errorf("process scanner error: %v", scanner.Err())
	}

	return nil
}
