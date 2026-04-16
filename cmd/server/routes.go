package main

import (
	"bytes"
	"compress/gzip"
	"io"
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

	if handler.Size > 1e9 {
		log.Println("file too big")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tableSize int64
	err = DB.QueryRow(`
		select sum(length(data)) from logs;
	`).Scan(&tableSize)

	if err != nil {
		log.Println("failed to retrieve logs size")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if handler.Size+tableSize > 1e9 {
		log.Println("max log size exceeded")
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}

	var buf bytes.Buffer
	bufWriter := gzip.NewWriter(&buf)
	defer bufWriter.Close()

	_, err = io.Copy(bufWriter, file)
	if err != nil {
		log.Println("failed to compress file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec(`
		insert into logs (ip, size, compress, data)
		values (?, ?, ?, ?)`,
		getIP(r),
		handler.Size,
		"gzip",
		buf.Bytes(),
	)
	if err != nil {
		log.Println("failed to save file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
