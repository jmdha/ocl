package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"ocl/db"
	"ocl/web"
	"text/template"
	"time"

	"github.com/klauspost/compress/zstd"
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
		ReadTimeout:    0 * time.Second,
		WriteTimeout:   0 * time.Second,
		IdleTimeout:    0 * time.Second,
		MaxHeaderBytes: 64 << 10, // 64 kb
	}

	// Start server
	log.Fatal(srv.ListenAndServe())
}

func GetRouteIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	type Data struct {
		QueueActive   int
		QueueTotal    int
		StorageUsed   float64
		StorageTotal  float64
		Requests24H   int
		RequestsTotal int
		Visitors24H   int
		VisitorsTotal int
	}

	StorageUsed, err := db.Size()
	QueueActive, err := db.QueueActive()
	QueueTotal, err := db.QueueTotal()
	Requests24H, err := db.Requests24H()
	RequestsTotal, err := db.RequestsTotal()
	Visitors24H, err := db.Visitors24H()
	VisitorsTotal, err := db.VisitorsTotal()

	data := Data{
		QueueActive,
		QueueTotal,
		float64(StorageUsed) / 1024 / 1024 / 1024,
		4,
		Requests24H,
		RequestsTotal,
		Visitors24H,
		VisitorsTotal,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func PostUpload(w http.ResponseWriter, r *http.Request) {
	const MaxSize = 2 << 30

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

	var buf bytes.Buffer
	encoder, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		db.ImportMarkFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(encoder, r.Body)
	if err != nil {
		db.ImportMarkFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = encoder.Close()
	if err != nil {
		db.ImportMarkFailed(importID, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.ImportMarkPending(importID, buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
