package routes

import (
	"github.com/gorilla/mux"
	h "kufa.io/sqlitego/api/handler"
)

func NewRouter(h *h.Handler) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// GET
	r.HandleFunc("/tasks", h.GetAllTasks).Methods("GET")
	r.HandleFunc("/task/{id}", h.GetTask).Methods("GET")

	// POST
	r.HandleFunc("/task/", h.CreateTask).Methods("POST")

	// PUT
	r.HandleFunc("/task", h.UpdateTask).Methods("PUT")

	// DELETE
	r.HandleFunc("/task/{id}", h.DeleteTask).Methods("DELETE")

	return r
}
