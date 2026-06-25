package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
		err = db.QueryRow(`
		update import
		set status = 'processing'
		where id = (
			select id
			from import
			where status = 'pending'
			limit 1
		) returning id, path;
		`).Scan(&job, &path)
		if errors.Is(err, sql.ErrNoRows) {
			time.Sleep(time.Second)
			continue
		}

		log.Printf("[worker %d] processing job %d", id, job)
		err = clean(job)
		if err != nil {
			log.Printf("[worker %d] clean %v", id, err)
			db.Exec(`
				update import 
				set status = 'error'
				where id = ?
			`, job)
			continue
		}
		err = process(job, path)
		log.Printf("[worker %d] finish %d %v", id, job, err)
		if err != nil {
			clean(job)
			db.Exec(`
				update import 
				set status = 'failed'
				where id = ?
			`, job)
			continue
		}
		db.Exec(`
			update import 
			set status = 'done'
			where id = ?
		`, job)
	}
}

func clean(job int64) error {
	_, err := db.Exec(`
		delete from encounter where import_id = ?
	`, job)
	if err != nil {
		return fmt.Errorf("%v: remove encounters", err)
	}

	_, err = db.Exec(`
		delete from challenge where import_id = ?
	`, job)
	if err != nil {
		return fmt.Errorf("%v: remove challenges", err)
	}

	_, err = db.Exec(`
		delete from event_damage where import_id = ?
	`, job)
	if err != nil {
		return fmt.Errorf("%v: remove event_damage", err)
	}

	return nil
}

func process(job int64, path string) error {
	in, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%v: in file open fail", err)
	}
	defer in.Close()

	decoder, err := zstd.NewReader(in)
	if err != nil {
		return fmt.Errorf("%v: create decoder failed", err)
	}
	defer decoder.Close()

	var scanner *bufio.Scanner
	var event wcl.Event

	scanner = bufio.NewScanner(decoder)
	for scanner.Scan() {
		err = wcl.Parse(&event, scanner.Text())
		if err != nil {
			return nil
		}

		err = process_event(job, &event)
		if err != nil {
			return nil
		}
	}

	if scanner.Err() != nil {
		return nil
	}

	return nil
}

func process_event(job int64, e *wcl.Event) error {
	var err error
	switch e.Kind {
	case wcl.KindSpellDamage:
		_, err = db.Exec(`
			insert into event_damage (import_id, timestamp, source_id, target_id, spell_id, amount)
			values (?, ?, ?, ?, ?, ?)
		`,
			job,
			e.SpellDamage.Timestamp,
			400,
			200,
			e.SpellDamage.SpellID,
			e.SpellDamage.Damage,
		)
	}
	return err
}
