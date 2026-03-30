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

var DB *sql.DB

func main() {
	var addr string
	var port int
	var conn string

	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.StringVar(&conn, "c", "file:ocl.sqlite", "path to db")
	flag.Parse()

	db_init(conn)

	sub, _ := fs.Sub(web.Static, "static")

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.FS(sub)))
	mux.HandleFunc("GET  /api/logs", routeAPILogs)
	mux.HandleFunc("GET  /api/runs", routeAPIRuns)
	mux.HandleFunc("GET  /api/players", routeAPIPlayers)
	mux.HandleFunc("GET  /api/queue", routeAPIQueue)
	mux.HandleFunc("POST /api/upload", routeAPIUpload)
	mux.HandleFunc("GET /api/metrics/routes", routeAPIMetricsRoutes)

	l, err := logger.NewLogger(DB)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	handler := l.Middleware(mux)

	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), handler)
}

func db_init(conn string) {
	var err error

	DB, err = sql.Open("sqlite3", fmt.Sprintf("%s", conn))
	if err != nil {
		log.Fatalf("db creation failed with conn %s error %v", conn, err)
	}

	_, err = DB.Exec(`PRAGMA journal_mode=WAL`)
	if err != nil {
		log.Fatalf("setting WAL failed with error %v", err)
	}

	DB.SetMaxOpenConns(1)
}
