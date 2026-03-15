package main

import (
	"io/fs"
	"net/http"
	"ocl/web"
)

func main() {
	sub, _ := fs.Sub(web.Static, "static")
	http.Handle("GET /", http.FileServer(http.FS(sub)))

	http.HandleFunc("GET  /api/logs",      routeAPILogs)
	http.HandleFunc("GET  /api/logs/{ID}", routeAPILogsID)
	http.HandleFunc("POST /api/upload",    routeAPIUpload)
	
	http.ListenAndServe(":8080", nil)
}
