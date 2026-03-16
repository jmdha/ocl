package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"sync"
)

var queue_mu  sync.Mutex
var queue_itr uint64

func routeAPILogs(w http.ResponseWriter, r *http.Request) {
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

	reader, err := gzip.NewReader(r.Body)
	if err != nil {
		http.Error(w, "invalid gzip", http.StatusBadRequest)
		return
	}

	queue_mu.Lock()
	idx := queue_itr + 1
	queue_itr = queue_itr + 1
	queue_mu.Unlock()

	err = CompressToFile(reader, fmt.Sprintf("queue/%d", idx))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save file %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
