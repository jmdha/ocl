package main

import (
	"log"
	"net/http"
)

func routeIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	err = Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println("routeIndex faield with ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func routeMetrics(w http.ResponseWriter, r *http.Request) {
	var data Metrics
	var err error

	data, err = repoMetrics()
	if err != nil {
		log.Println("routeMetrics faield with ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = Templates.ExecuteTemplate(w, "metrics.html", data)
	if err != nil {
		log.Println("routeMetrics faield with ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func routeAPIUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("routeAPIUpload faield with ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size > 1*1e9 {
		log.Println("routeAPIUpload faield with max file size exceeded")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
