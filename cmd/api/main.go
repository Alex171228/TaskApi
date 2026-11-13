package main

import (
	"TaskAPI2/internal/config"
	"TaskAPI2/internal/httpserver"
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
		log.Printf("No db connection")
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db %v", err)
		}
	}()
	migrcontext, cancelmigr := context.WithTimeout(context.Background(), 10*time.Second)
	err = postgres.Migrate(migrcontext, db)
	cancelmigr()
	if err != nil {
		log.Printf("Migrate error %v", err)
		return
	}
	//repo := memory.NewMemoryRepository()
	repodb := postgres.NewTaskRepository(db)
	service := task.NewService(repodb)
	handler := taskhttp.NewHandler(service)
	server := httpserver.New(cnfg.HTTPAddr, handler)
	log.Printf("TodoAPI запущен на %s", cnfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
