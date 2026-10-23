package taskhttp

import (
	"TaskAPI2/internal/httpx"
	"TaskAPI2/internal/task"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *task.Service
}

func NewHandler(s *task.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) searchID(s string) (int, bool) {
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	return id, true
}

func (h *Handler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTasks(w, r)
	case http.MethodPost:
		h.CreateHandler(w, r)
	default:
		httpx.MethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}
func (h *Handler) TaskbyIdHandler(w http.ResponseWriter, r *http.Request) {
	if id, ok := h.searchID(r.PathValue("id")); !ok || id <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "Not valid id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.GetbyID(w, r)
	case http.MethodPut:
		h.PutHandler(w, r)
	case http.MethodDelete:
		h.DeleteHandler(w, r)
	default:
		httpx.MethodNotAllowed(w, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := h.searchID(r.PathValue("id"))
	if !ok || id <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "Not valid id")
		return
	}

	var body TaskRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	bodyconv := ConvertRequest(body)
	t, err := h.service.Update(r.Context(), id, bodyconv)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, task.ErrNotFound.Error())
		case errors.Is(err, task.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, task.ErrInvalidInput.Error())
		default:
			log.Println(err)
			httpx.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	taskansw := ConvertAnswer(t)
	httpx.WriteJSON(w, http.StatusOK, taskansw)
}
func (h *Handler) GetbyID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.searchID(r.PathValue("id"))
	if !ok || id <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "Not valid id")
		return
	}
	t, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	taskansw := ConvertAnswer(t)
	httpx.WriteJSON(w, http.StatusOK, taskansw)

}
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var body TaskRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	bodyconv := ConvertRequest(body)
	t, err := h.service.Create(r.Context(), bodyconv)
	if err != nil {
		if errors.Is(err, task.ErrInvalidInput) {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Println(err)
		httpx.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	taskansw := ConvertAnswer(t)
	httpx.WriteJSON(w, http.StatusCreated, taskansw)
}
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	t, err := h.service.List(r.Context())
	if err != nil {
		log.Println(err)
		httpx.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	answtasks := make([]TaskAnswer, len(t))
	for i, tasks := range t {
		answtasks[i] = ConvertAnswer(tasks)
	}
	httpx.WriteJSON(w, http.StatusOK, answtasks)
}
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := h.searchID(r.PathValue("id"))
	if !ok || id <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "Not valid id")
		return
	}
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, task.ErrNotFound.Error())
			return
		}
		log.Println(err)
		httpx.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
