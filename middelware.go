package main

import (
	"log"
	"net"
	"net/http"
	"ocl/db"
	"time"
)

func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()

		next.ServeHTTP(w, r)

		ip := r.Header.Get("X-Forwarded-For")

		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}

		if ip == "" {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				ip = host
			} else {
				ip = r.RemoteAddr
			}
		}

		err := db.RequestsAdd(
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			ip,
			r.UserAgent(),
			int64(time.Since(t0).Nanoseconds()),
		)
		if err != nil {
			log.Printf("request logging failed: %v", err)
		}
	})
}
