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
		log.Printf("%s %s", r.Method, r.URL.Path)
		t0 := time.Now()
		next.ServeHTTP(w, r)
		t1 := time.Now()

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

		err := db.RequestAdd(
			ip,
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			r.UserAgent(),
			t1.Sub(t0),
		)
		if err != nil {
			log.Printf("request logging failed: %v", err)
		}
	})
}
