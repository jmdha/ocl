package main

import (
	"net/http"
)

func routeIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	err = Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func routeMetrics(w http.ResponseWriter, r *http.Request) {
	var data Metrics
	var err error

	data, err = repoMetrics()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = Templates.ExecuteTemplate(w, "metrics.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
