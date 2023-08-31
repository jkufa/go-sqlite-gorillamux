package routes

import (
	"github.com/gorilla/mux"
	"kufa.io/sqlitego/api/handler" // Import your handler package
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/tasks", handler.GetAllTasks).Methods("GET")

	return r
}
