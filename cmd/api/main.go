package main

import (
	"TaskAPI2/internal/memory"
	"TaskAPI2/internal/task"
	"TaskAPI2/internal/taskhttp"
	"log"
	"net/http"
)

func main() {
	repo := memory.NewMemoryRepository()
	service := task.NewService(repo)
	handler := taskhttp.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.TasksHandler)
	mux.HandleFunc("/tasks/{id}", handler.TaskbyIdHandler)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
	log.Println("TodoAPI запущен на 8081")
}
