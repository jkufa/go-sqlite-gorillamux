package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"kufa.io/sqlitego/todo"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library (driver for sqlite) without using it explicitly
)

const fileName = "sqlite.db"

func main() {
	// check for parameter to delete the database file
	if len(os.Args) > 1 && os.Args[1] == "--restart" {
		os.Remove(fileName)
		fmt.Println("Deleted database file")
	}

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	todoRepo := todo.NewSQLiteRepository(db)
	if err := todoRepo.Migrate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	writeChannel := make(chan Command)
	readChannel := make(chan Command)

	go func() {
		readChannel <- GetAllTasks
	}()

	go func() {
		writeChannel <- CreateTask
	}()

	// read/write depending on the channel
	select {
	case read := <-readChannel:
		GetAll(todoRepo)
		fmt.Println("Read:", read)
	case write := <-writeChannel:
		fmt.Println("Write:", write)
		Create(todoRepo)
	}
}

func GetAll(r *todo.SQLiteRepository) {
	// Get all todos
	todos, err := r.All()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All todos:", todos)
}

func Create(r *todo.SQLiteRepository) *todo.Todo {
	// Create a new todo
	todo1 := todo.Todo{
		Name:        "Learn Go",
		Description: "Learn Golang by solving algorithms and writing code!",
		Completed:   false,
	}
	created, err := r.Create(todo1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created todo:", created)
	return created
}
