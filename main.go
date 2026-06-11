package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"ocl/db"
	"ocl/web"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var templates *template.Template

func main() {
	var addr string
	var port int
	var conn string
	var workers int
	var err error

	// Parse program args
	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port, "p", 8080, "port to operate on")
	flag.StringVar(&conn, "c", "ocl.sqlite", "path to db")
	flag.IntVar(&workers, "w", 1, "number of processing workers")
	flag.Parse()

	// Init db
	err = db.Init(conn)
	if err != nil {
		log.Fatalf("db init failed with %v", err)
	}

	// Init templates
	templates = template.Must(template.ParseFS(web.Templates, "templates/*.html"))

	// Start processing worker(s)
	for i := 0; i < workers; i++ {
		go worker(i)
	}

	// Register routes
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(web.Static)))

	mux.HandleFunc("GET  /", GetRouteIndex)
	mux.HandleFunc("POST /upload", PostUpload)

	// Create server
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", addr, port),
		Handler:        MiddlewareLogging(mux),
		ReadTimeout:    8 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	log.Fatal(srv.ListenAndServe())
}

func GetRouteIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	err = templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func PostUpload(w http.ResponseWriter, r *http.Request) {
	const MaxSize = 4 << 30 // 4 GB

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpPath := filepath.Join("uploads", fmt.Sprintf("%d.tmp", importID))
	outPath := filepath.Join("uploads", fmt.Sprintf("%d", importID))

	f, err := os.Create(tmpPath)
	if err != nil {
		db.ImportMarkFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, copyErr := io.Copy(f, r.Body)
	closeErr := f.Close()

	if copyErr != nil {
		os.Remove(tmpPath)
		db.ImportMarkFailed(importID, copyErr.Error())
		http.Error(w, copyErr.Error(), http.StatusBadRequest)
		return
	}

	if closeErr != nil {
		os.Remove(tmpPath)
		db.ImportMarkFailed(importID, closeErr.Error())
		http.Error(w, closeErr.Error(), http.StatusInternalServerError)
		return
	}

	err = os.Rename(tmpPath, outPath)
	if err != nil {
		db.ImportMarkFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.ImportMarkPending(importID, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"import_id": importID,
		"status":    "pending",
		"size":      n,
	})
}
