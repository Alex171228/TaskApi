package httpserver

import (
	"TaskAPI2/internal/httpx"
	"TaskAPI2/internal/taskhttp"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, request *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func New(addr string, handler *taskhttp.Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.TasksHandler)
	mux.HandleFunc("/tasks/{id}", handler.TaskbyIdHandler)
	mux.HandleFunc("GET /health", HealthHandler)
	return &http.Server{
		Addr:    addr,
		Handler: mux}
}
