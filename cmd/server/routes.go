package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
)

func routeAPILogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "%d", 27)
}

func routeAPIRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "%d", 27)
}

func routeAPIPlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "%d", 27)
}

func routeAPILogsID(w http.ResponseWriter, r *http.Request) {
}

func routeAPIUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Encoding") != "gzip" {
		http.Error(w, "invalid encoding", http.StatusBadRequest)
		return
	}

	_, err := gzip.NewReader(r.Body)
	if err != nil {
		http.Error(w, "invalid gzip", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func routeAPIQueue(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", 0)
	w.WriteHeader(http.StatusOK)
}
