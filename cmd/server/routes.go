package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
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

func routeAPIMetricsRoutes(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("select method, path, count(*), avg(duration) from requests group by method, path order by count(*) desc;")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type rowdata struct {
		Method string  `json:"method"`
		Path   string  `json:"path"`
		Calls  uint    `json:"calls"`
		Avg    float64 `json:"avg"`
	}
	var data []rowdata
	for rows.Next() {
		var rd rowdata
		err = rows.Scan(&rd.Method, &rd.Path, &rd.Calls, &rd.Avg)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		rd.Avg = rd.Avg / 1000
		data = append(data, rd)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
