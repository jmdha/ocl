package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"ocl/web"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var Templates *template.Template

func main() {
	var addr string
	var port int
	var conn string

	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.StringVar(&conn, "c", "file:ocl.sqlite", "path to db")
	flag.Parse()

	db_init(conn)
	Templates = template.Must(template.ParseFS(web.Templates, "templates/*.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET  /", routeIndex)
	mux.HandleFunc("GET  /metrics", routeMetrics)
	mux.HandleFunc("POST /api/upload/multipart", routeAPIUploadMultipart)

	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), req_log(mux))
}

func db_init(conn string) {
	var err error

	DB, err = sql.Open("sqlite3", fmt.Sprintf("%s", conn))
	if err != nil {
		log.Fatalf("db creation failed with conn %s error %v", conn, err)
	}

	_, err = DB.Exec(`pragma journal_mode=WAL`)
	if err != nil {
		log.Fatalf("setting WAL failed with error %v", err)
	}

	_, err = DB.Exec(`pragma busy_timeout=60000`)
	if err != nil {
		log.Fatalf("setting WAL failed with error %v", err)
	}

	_, err = DB.Exec(`
		create table if not exists requests (
			id        integer primary key autoincrement,
			timestamp datetime not null default current_timestamp,
			ip        text not null,
			method    text not null,
            		path      text not null,
			query     text not null,
			agent     text not null,
			duration  integer not null
		);
	`)
	if err != nil {
		log.Fatalf("create requests failed with error %v", err)
	}

	_, err = DB.Exec(`
		create index if not exists requests_index on requests (
			method,
			path
		);
	`)
	if err != nil {
		log.Fatalf("failed to create index %v\n", err)
	}

	DB.SetMaxOpenConns(1)
}

func req_log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(t0).Nanoseconds()
		_, err := DB.Exec(`
			insert into requests (method, path, query, ip, agent, duration)
			values (?, ?, ?, ?, ?, ?)`,
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			getIP(r),
			r.UserAgent(),
			duration,
		)
		if err == nil {
			log.Printf("%s %s", r.Method, r.URL.Path)
		} else {
			log.Printf("Failed to log request: %v", err)
		}
	})
}

// cannot just use request.RemoteAddr as the server may be behind a proxy
func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
