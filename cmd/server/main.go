package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"ocl/web"
)

func main() {
	var addr string
	var port int

	flag.StringVar(&addr, "a", "localhost", "address to operate on")
	flag.IntVar(&port,    "p", 8080,        "port to operate on")
	flag.Parse()

	sub, _ := fs.Sub(web.Static, "static")
	http.Handle("GET /", http.FileServer(http.FS(sub)))

	http.HandleFunc("GET  /api/logs",    routeAPILogs)
	http.HandleFunc("GET  /api/runs",    routeAPIRuns)
	http.HandleFunc("GET  /api/players", routeAPIPlayers)
	http.HandleFunc("POST /api/upload",  routeAPIUpload)

	
	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
}
