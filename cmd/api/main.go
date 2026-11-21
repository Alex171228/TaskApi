package main

import (
	"TaskAPI2/internal/config"
	"TaskAPI2/internal/httpserver"
	"TaskAPI2/internal/postgres"
	"TaskAPI2/internal/task"
	"TaskAPI2/internal/taskhttp"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()
	stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed with error")
		}
		return
	case <-stopContext.Done():
		log.Println("Shotdown signal recived")

	}
	contextshot, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(contextshot); err != nil {
		log.Printf("Server shutdown with %v", err)
	}
}
