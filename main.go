package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"ocl/db"
	"ocl/web"
	"text/template"
	"time"
)

var templates *template.Template

func main() {
	var addr string
	var port int
	var conn string
	var workers int
	var err error

	// Parse program args
	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.StringVar(&conn, "c", "ocl.sqlite", "path to db")
	flag.IntVar(&workers, "w", 1, "number of processing workers")
	flag.Parse()

	// Init db
	err = db.Init(conn)
	if err != nil {
		log.Fatalf("db init failed with %v", err)
	}

	// Init templates
	templates = template.Must(template.ParseFS(web.Templates, "templates/*.html"))

	// Start processing worker(s)
	for i := 0; i < workers; i++ {
		go worker(i)
	}

	// Register routes
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(web.Static)))

	mux.HandleFunc("/", RouteIndex)
	mux.HandleFunc("/encounters", RouteEncounters)
	mux.HandleFunc("/encounters/{id}", RouteEncountersID)
	mux.HandleFunc("/logs", RouteLogs)
	mux.HandleFunc("/logs/{id}", RouteLogsID)
	mux.HandleFunc("/logs/{id}/{entry}", RouteLogsIDEntry)
	mux.HandleFunc("/upload", RouteUpload)

	// Create server
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", addr, port),
		Handler:        MiddlewareLogging(mux),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 64 << 10, // 64 kb
	}

	// Start server
	log.Fatal(srv.ListenAndServe())
}
