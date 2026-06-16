package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"ocl/db"
	"ocl/web"
	"strconv"

	"github.com/klauspost/compress/zstd"
)

func RouteIndex(w http.ResponseWriter, r *http.Request) {
	var data web.DataIndex
	var err error

	data, err = db.WebDataIndex()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteCharacters(w http.ResponseWriter, r *http.Request) {
	var data web.DataCharacters
	var err error

	data, err = db.WebDataCharacters()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "characters.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteCharactersID(w http.ResponseWriter, r *http.Request) {
	var data web.DataCharactersID
	var err error

	data, err = db.WebDataCharactersID()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "characters_id.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteEncounters(w http.ResponseWriter, r *http.Request) {
	var data web.DataEncounters
	var err error

	data, err = db.WebDataEncounters()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "encounters.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteEncountersID(w http.ResponseWriter, r *http.Request) {
	var data web.DataEncountersID
	var err error

	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil || id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err = db.WebDataEncountersID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "encounters_id.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteLogs(w http.ResponseWriter, r *http.Request) {
	var data web.DataLogs
	var err error

	data, err = db.WebDataLogs(500, 0)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "logs.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteLogsID(w http.ResponseWriter, r *http.Request) {
	var data web.DataLogsID
	var err error

	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil || id < 0 {
		log.Printf("invalid id %s: %v", sid, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err = db.WebDataLogsID(id)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "logs_id.html", data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func RouteLogsIDEntry(w http.ResponseWriter, r *http.Request) {
	var data web.DataLogsIDEntry
	var err error

	data, err = db.WebDataLogsIDEntry()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	importID, err := db.ImportAdd()
	if err != nil {
		log.Printf("upload failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	encoder, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedDefault))
	if err != nil {
		log.Printf("upload failed: %v", err)
		db.ImportSetFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer encoder.Close()

	_, err = io.Copy(encoder, r.Body)
	if err != nil {
		log.Printf("upload failed: %v", err)
		db.ImportSetFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = encoder.Close()
	if err != nil {
		log.Printf("upload failed: %v", err)
		db.ImportSetFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.ImportSetPending(importID, buf.Bytes())
	if err != nil {
		log.Printf("upload failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
