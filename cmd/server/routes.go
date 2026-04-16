package main

import (
	"github.com/gabriel-vasile/mimetype"
	"log"
	"net/http"
)

func routeIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	err = Templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println("failure to serve index", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func routeMetrics(w http.ResponseWriter, r *http.Request) {
	var data Metrics
	var err error

	data, err = repoMetrics()
	if err != nil {
		log.Println("failure to retrieve metrics", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = Templates.ExecuteTemplate(w, "metrics.html", data)
	if err != nil {
		log.Println("failure to serve metrics", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func routeAPIUploadMultipart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("non-allowed method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("failure to retrieve file", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size > 1*1e9 {
		log.Println("file too big")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		log.Println("failure to detect mime type", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	allowed := []string{"text/plain"}
	if !mimetype.EqualsAny(mime.String(), allowed...) {
		log.Println("unsupported mime type", mime)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
