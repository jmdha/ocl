package main

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
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

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("routeAPIUpload faield with ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create("tmp/logs" + uuid.New().String())
	if err != nil {
		log.Println("routeAPIUpload faield with ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)
	w.WriteHeader(http.StatusOK)
}
