package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
)

func routeAPILogs(w http.ResponseWriter, r *http.Request) {
}

func routeAPILogsID(w http.ResponseWriter, r *http.Request) {
}

func routeAPIUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	var reader io.ReadCloser
	var err error

	log.Println("Content-Encoding:", r.Header.Get("Content-Encoding"))

	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "invalid gzip", http.StatusBadRequest)
			return
		}
		defer reader.Close()
	default:
		reader = r.Body
	}
	defer r.Body.Close()

	_, err = io.ReadAll(reader)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
