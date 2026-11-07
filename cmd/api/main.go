package main

import (
	"TaskAPI2/internal/config"
	"TaskAPI2/internal/httpserver"
	"TaskAPI2/internal/memory"
	"TaskAPI2/internal/postgres"
	"TaskAPI2/internal/task"
	"TaskAPI2/internal/taskhttp"
	"context"
	"log"
	"time"
)

func main() {
	cnfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	cntxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := postgres.Open(cntxt, cnfg.DatabaseURL)
	cancel()
	if err != nil {
		log.Printf("No db connection %v", err)
		return
	}
	defer db.Close()
	repo := memory.NewMemoryRepository()
	service := task.NewService(repo)
	handler := taskhttp.NewHandler(service)
	server := httpserver.New(cnfg.HTTPAddr, handler)
	log.Printf("TodoAPI запущен на %s", cnfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
