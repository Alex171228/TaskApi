package internal

import (
	"TaskAPI2/internal/task"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type MemoryRepository struct {
	tasks []task.Task
	id    int
	mu    sync.Mutex
}

func writeJSON(w http.ResponseWriter, status int, value any) {
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
func Put(w http.ResponseWriter, r *http.Request) {
	id, _ := searchID(r.PathValue("id"))
	mu.Lock()
	defer mu.Unlock()
	for i, t := range tasks {
		if t.ID == id {
			var input task.Input
			err := json.NewDecoder(r.Body).Decode(&input)
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
func Create(w http.ResponseWriter, r *http.Request) {
	var input task.Input
	err := json.NewDecoder(r.Body).Decode(&input)
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
func Get(w http.ResponseWriter, _ *http.Request) {
	if tasks == nil {
		b := make([]task.Task, 0)
		writeJSON(w, http.StatusOK, b)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}
func Delete(w http.ResponseWriter, r *http.Request) {
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
