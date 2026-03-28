package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"net/http"
	"ocl/pkg/logger"
	"ocl/web"
)

func main() {
	var addr string
	var port int
	var conn string

	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.StringVar(&conn, "c", "file:ocl.sqlite", "path to db")
	flag.Parse()

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s", conn))
	if err != nil {
		log.Fatalf("failed to open db conn: %v %v", conn, err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		log.Fatalf("failed to set db WAL: %v", err)
	}

	db.SetMaxOpenConns(1)

	sub, _ := fs.Sub(web.Static, "static")

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.FS(sub)))
	mux.HandleFunc("GET  /api/logs", routeAPILogs)
	mux.HandleFunc("GET  /api/runs", routeAPIRuns)
	mux.HandleFunc("GET  /api/players", routeAPIPlayers)
	mux.HandleFunc("GET  /api/queue", routeAPIQueue)
	mux.HandleFunc("POST /api/upload", routeAPIUpload)

	l, err := logger.NewLogger(db)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	handler := l.Middleware(mux)

	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), handler)
}
