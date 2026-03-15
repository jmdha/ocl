package main

import (
	"net/http"
)

func main() {
	http.Handle("GET /", http.FileServer(http.Dir("./web")))

	http.HandleFunc("GET  /api/logs",      routeAPILogs)
	http.HandleFunc("GET  /api/logs/{ID}", routeAPILogsID)
	http.HandleFunc("POST /api/upload",    routeAPIUpload)
	
	http.ListenAndServe(":8080", nil)
}
