package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/klauspost/compress/zstd"
)

func RouteIndex(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteCharacters(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	err = templates.ExecuteTemplate(w, "characters.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteCharactersID(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	err = templates.ExecuteTemplate(w, "characters_id.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteEncounters(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	err = templates.ExecuteTemplate(w, "encounters.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteEncountersID(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil || id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = templates.ExecuteTemplate(w, "encounters_id.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteLogs(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	err = templates.ExecuteTemplate(w, "logs.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteLogsID(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil || id < 0 {
		log.Printf("invalid id %s: %v", sid, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = templates.ExecuteTemplate(w, "logs_id.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteLogsIDEntry(w http.ResponseWriter, r *http.Request) {
	type Data struct {
	}

	var data Data
	var err error

	err = templates.ExecuteTemplate(w, "logs_id_entry.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteUpload(w http.ResponseWriter, r *http.Request) {
	const MaxSize = 2 * 1024 * 1024 * 1024 // 2 GB

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.ContentLength > MaxSize {
		http.Error(w, "too large", http.StatusRequestEntityTooLarge)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxSize)
	defer r.Body.Close()

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	file, err := os.CreateTemp("logs", "")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder, err := zstd.NewWriter(file)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(encoder, r.Body)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = encoder.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		`insert into import(timestamp, path, status) values (?, ?, ?)`,
		time.Now(),
		file.Name(),
		"pending",
	)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
