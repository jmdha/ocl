package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

	if handler.Size > 1e9 {
		log.Println("file too big")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = os.MkdirAll("uploads", 0755)
	if err != nil {
		log.Println("failed to create dir", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()
	outPath := filepath.Join("uploads", id+".gz")
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Println("failed to create file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	gz := gzip.NewWriter(outFile)
	defer gz.Close()

	_, err = io.Copy(gz, file)
	if err != nil {
		log.Println("failed to save file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
