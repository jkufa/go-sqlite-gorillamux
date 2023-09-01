package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"kufa.io/sqlitego/db/models"
	repo "kufa.io/sqlitego/db/repository"
	jsonmapper "kufa.io/sqlitego/utils"
)

type Handler struct {
	repository repo.SQLiteRepository
}

func NewHandler(repository *repo.SQLiteRepository) *Handler {
	return &Handler{repository: *repository}
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.repository.All()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "All tasks: %v", tasks)

	// Respond with tasks in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	// Decode the incoming Task json
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Fatal(err)
	}

	// Create the task
	newTask, err := h.repository.Create(task)
	if err != nil {
		log.Fatal(err)
	}

	// Respond with the user data in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Get the task id from the route parameters
	id := jsonmapper.StrToInt64(mux.Vars(r)["id"])

	// Get the task
	task, err := h.repository.Get(id)
	if err != nil {
		log.Fatal(err)
	}

	// Respond with the user data in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming Task json
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Fatal(err)
	}

	// Update the task
	updatedTask, err := h.repository.Update(task)
	if err != nil {
		log.Fatal(err)
	}

	// Respond with the task in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get the task id from the route parameters
	id := jsonmapper.StrToInt64(mux.Vars(r)["id"])

	// Delete the task
	err := h.repository.Delete(id)
	if err != nil {
		log.Fatal(err)
	}

	// Respond with the task in JSON format
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Task with id %d was deleted", id)
}
