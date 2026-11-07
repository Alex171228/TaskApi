package main

import (
	"TaskAPI2/internal/config"
	"TaskAPI2/internal/httpserver"
	"TaskAPI2/internal/memory"
	"TaskAPI2/internal/task"
	"TaskAPI2/internal/taskhttp"
	"log"
)

func main() {
	cnfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	repo := memory.NewMemoryRepository()
	service := task.NewService(repo)
	handler := taskhttp.NewHandler(service)
	server := httpserver.New(cnfg.HTTPAddr, handler)
	log.Printf("TodoAPI запущен на %s", cnfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
