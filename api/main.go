package main

import (
	"database/sql"
	"log"
	"net/http"

	"kufa.io/sqlitego/api/handler"
	"kufa.io/sqlitego/api/routes"
	repo "kufa.io/sqlitego/db/repository"
)

func main() {
	// Connect to the database
	fileName := "../sqlite.db"
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	repository := repo.NewSQLiteRepository(db)

	// initialize handler
	h := handler.NewHandler(repository)

	// initialize the router
	r := routes.NewRouter(h)

	port := "8080"
	log.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
