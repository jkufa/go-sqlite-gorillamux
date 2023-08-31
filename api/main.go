package main

import (
	"log"
	"net/http"

	"kufa.io/sqlitego/api/routes"
)

func main() {
	r := routes.NewRouter()

	port := "8080"
	log.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
