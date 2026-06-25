package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"ocl/web"
	"text/template"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*
var migrations embed.FS
var db *sql.DB
var templates *template.Template

func main() {
	var addr string
	var port int
	var conn string
	var workers int
	var storage string

	// Init logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse program args
	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.StringVar(&conn, "c", "ocl.sqlite", "path to db")
	flag.IntVar(&workers, "w", 1, "number of processing workers")
	flag.StringVar(&storage, "s", "", "scp storage path")
	flag.Parse()

	// Init db
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		d,
		fmt.Sprintf("sqlite://%s", conn),
	)
	if err != nil {
		log.Fatal(err)
	}

	m.Up()

	db, err = sql.Open("sqlite", fmt.Sprintf("%s?_pragma=busy_timeout(10000)&_pragma=journal_mode(wal)&_pragma=synchronous(normal)&_time_integer_format=unix_milli", conn))
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	// Init templates
	templates = template.Must(template.ParseFS(web.Public, "public/*"))

	_, err = db.Exec(`update import set status = 'pending'`)
	if err != nil {
		log.Fatal(err)
	}

	// Start processing worker(s)
	for i := range workers {
		go worker(i)
	}

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET  /", RouteIndex)
	mux.HandleFunc("GET  /characters", RouteCharacters)
	mux.HandleFunc("GET  /characters/{id}", RouteCharactersID)
	mux.HandleFunc("GET  /encounters", RouteEncounters)
	mux.HandleFunc("GET  /encounters/{id}", RouteEncountersID)
	mux.HandleFunc("GET  /logs", RouteLogs)
	mux.HandleFunc("GET  /logs/{id}", RouteLogsID)
	mux.HandleFunc("GET  /logs/{id}/{entry}", RouteLogsIDEntry)
	mux.HandleFunc("POST /upload", RouteUpload)

	// Create server
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", addr, port),
		Handler:        mux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 64 << 10, // 64 kb
	}

	// Start server
	log.Fatal(srv.ListenAndServe())
}
