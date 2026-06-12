package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*
var migrations embed.FS

var db *sql.DB

func Init(conn string) error {
	var d source.Driver
	var m *migrate.Migrate
	var err error

	d, err = iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err = migrate.NewWithSourceInstance(
		"iofs",
		d,
		fmt.Sprintf("sqlite://%s", conn),
	)
	if err != nil {
		return err
	}
	err = m.Up()

	db, err = sql.Open("sqlite", conn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(1)

	_, err = db.Exec(`pragma busy_timeout = 10000;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`PRAGMA journal_mode = WAL;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`PRAGMA synchronous = normal;`)
	if err != nil {
		return err
	}

	return nil
}

func RequestsAdd(
	method string,
	path string,
	query string,
	ip string,
	agent string,
	duration time.Duration,
) error {
	var err error

	_, err = db.Exec(`
		insert into requests (timestamp, method, path, query, ip, agent, duration)
		values (?, ?, ?, ?, ?, ?, ?)`,
		time.Now().UTC().UnixMilli(),
		method,
		path,
		query,
		ip,
		agent,
		duration.Nanoseconds(),
	)

	return err
}
