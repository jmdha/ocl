package db

import (
	"database/sql"
	"embed"
	"fmt"

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

func Size() (int64, error) {
	var pageCount int64
	var pageSize int64
	var err error

	err = db.QueryRow(`PRAGMA page_count`).Scan(&pageCount)
	if err != nil {
		return 0, err
	}

	err = db.QueryRow(`PRAGMA page_size`).Scan(&pageSize)
	if err != nil {
		return 0, err
	}

	return pageCount * pageSize, nil
}
