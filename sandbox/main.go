package main

import (
	"TaskAPI2/internal/task"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var tasks []task.Task
var nextId = 1
var mu sync.Mutex

func writeJSON(w http.ResponseWriter, status int, value any) { //JSON encode и проверки полей
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func searchID(s string) (int, bool) {
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	return id, true
}
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	t := task.Task{
		ID:          nextId,
		Title:       "Hello world",
		Description: "Hello",
		Done:        true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	writeJSON(w, http.StatusOK, t)
	//w.Header().Set("Content-Type", "application/json")
	//encoder := json.NewEncoder(w)
	//if err := encoder.Encode(tasks); err != nil {
	// log.Println(err)
	//}
	//w.Write([]byte("Hello world"))
}
func PutHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := searchID(r.PathValue("id"))
	mu.Lock()
	defer mu.Unlock()
	for i, t := range tasks {
		if t.ID == id {
			var input task.Input
			err := json.NewDecoder(r.Body).Decode(&input) //JSON decode и проверки полей
			if err != nil {
				log.Println("Error decode JSON")
				writeError(w, http.StatusBadRequest, "Error decode JSON")
				return
			}
			if len(input.Title) < 3 || len(input.Description) < 3 {
				writeError(w, http.StatusBadRequest, "Title or description must be >3")
				return
			}

			t = task.Task{
				ID:          t.ID,
				Title:       input.Title,
				Description: input.Description,
				Done:        input.Done,
				CreatedAt:   t.CreatedAt,
				UpdatedAt:   time.Now(),
			}
			tasks[i] = t
			log.Println("Task Updated")
			writeJSON(w, http.StatusOK, t)
			return
		}
	}
	writeError(w, http.StatusNotFound, task.ErrNotFound.Error())
	log.Println(task.ErrNotFound.Error())
}
func GetbyID(w http.ResponseWriter, r *http.Request) {
	id, _ := searchID(r.PathValue("id"))
	mu.Lock()
	defer mu.Unlock()
	for _, t := range tasks {
		if t.ID == id {
			writeJSON(w, http.StatusOK, t)
			return
		}
	}
	writeError(w, http.StatusNotFound, task.ErrNotFound.Error())
	log.Println(task.ErrNotFound.Error())

}
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var input task.Input
	err := json.NewDecoder(r.Body).Decode(&input) //JSON decode и проверки полей
	if err != nil {
		log.Println("Error decode JSON")
		writeError(w, http.StatusBadRequest, "Error decode JSON")
		return
	}
	if len(input.Title) < 3 || len(input.Description) < 3 {
		writeError(w, http.StatusBadRequest, "Title or description must be >3")
		return
	}
	mu.Lock()
	defer mu.Unlock()
	t := task.Task{
		ID:          nextId,
		Title:       input.Title,
		Description: input.Description,
		Done:        input.Done,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	nextId += 1
	tasks = append(tasks, t)
	log.Println("Task Added")
	writeJSON(w, http.StatusOK, t)
}
func GetTasks(w http.ResponseWriter, _ *http.Request) {
	if tasks == nil {
		b := make([]task.Task, 0)
		writeJSON(w, http.StatusOK, b)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := searchID(r.PathValue("id"))
	mu.Lock()
	defer mu.Unlock()
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writeError(w, http.StatusNotFound, task.ErrNotFound.Error())
	log.Println(task.ErrNotFound.Error())
}
func TasksHandler(w http.ResponseWriter, r *http.Request) { // методы разбираются
	switch r.Method {
	case http.MethodGet:
		GetTasks(w, r)
	case http.MethodPost:
		CreateHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
func TaskbyIdHandler(w http.ResponseWriter, r *http.Request) {
	if id, ok := searchID(r.PathValue("id")); !ok || id <= 0 {
		writeError(w, http.StatusBadRequest, "Not valid id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		GetbyID(w, r)
	case http.MethodPut:
		PutHandler(w, r)
	case http.MethodDelete:
		DeleteHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
func main() {
	port := ":8081"
	mux := http.NewServeMux()
	//mux.HandleFunc("GET /", HelloHandler)
	mux.HandleFunc("/tasks", TasksHandler)
	mux.HandleFunc("/tasks/", TaskbyIdHandler)
	log.Printf("TodoAPI запущен на%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}

}
