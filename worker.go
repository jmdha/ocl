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
)

func worker(id int) {
	var job int64
	var err error
	for {
		job, err = db.ImportClaim()
		if errors.Is(err, sql.ErrNoRows) {
			time.Sleep(time.Second)
			continue
		}

		err = process(job, fmt.Sprintf("uploads/%d", job))
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

	line = 0
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		err = wcl.Parse(&event, scanner.Text())
		if err != nil {
			return fmt.Errorf("failed to parse file %v line %d", err, line)
		}

		switch event.Kind {
		case wcl.KindVersion:
			err = db.InsertWOWEventVersion(
				job,
				line,
				event.Version.Timestamp,
				event.Version.Log,
				event.Version.Advanced,
				event.Version.Major,
				event.Version.Minor,
				event.Version.Patch,
				event.Version.Project,
			)
		case wcl.KindUnknown:
			err = fmt.Errorf("line %d unknown event %s", line, scanner.Text())
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
