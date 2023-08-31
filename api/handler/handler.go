package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"kufa.io/sqlitego/db/repository"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API!")
}

/*
  - TODO: Remove code pertaining to opening the database connection as this should
    be handled by the repository layer
*/
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	fileName, _ := filepath.Abs("../db/sqlite.db")
	fmt.Println(fileName)
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewSQLiteRepository(db)
	tasks, err := repo.All()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "All tasks: %v", tasks)

	// Respond with the user data in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
