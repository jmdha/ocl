package logger

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"
)

type Logger struct {
	db *sql.DB
}

func NewLogger(db *sql.DB) (Logger, error) {
	_, err := db.Exec(`
		create table if not exists requests (
			id        integer primary key autoincrement,
			method    text not null,
            		path      text not null,
			query     text not null,
			ip        text not null,
			user      text not null,
			duration  integer not null,
			status    integer not null,
			timestamp datetime default current_timestamp
		)
	`)
	if err != nil {
		return Logger{}, err
	}
	return Logger{
		db: db,
	}, nil
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		wrapper := responseWriter{ResponseWriter: w}
		next.ServeHTTP(&wrapper, r)
		dur := time.Since(t0).Nanoseconds()
		go func() {
			if strings.Contains(r.URL.Path, "/api/") ||
				strings.Contains(r.URL.Path, ".css") {
				return
			}
			_, err := l.db.Exec(`
				insert into requests (method, path, query, ip, user, duration, status)
				values (?, ?, ?, ?, ?, ?, ?)`,
				r.Method,
				r.URL.Path,
				r.URL.RawQuery,
				r.RemoteAddr,
				r.UserAgent(),
				dur,
				wrapper.status,
			)
			if err != nil {
				log.Printf("Failed to log request: %v", err)
			}
		}()
	})
}
