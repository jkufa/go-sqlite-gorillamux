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
	commands := []Command{
		GetAllTasks,
		CreateTask,
		CreateTask,
		GetAllTasks,
		CreateTask,
		CreateTask,
		GetAllTasks,
		GetAllTasks,
	}

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

	writeChannel := make(chan Command, 1)
	readChannel := make(chan Command, 1)

	go func() {
		for c := range readChannel {
			fmt.Println("Read channel:", c)
			if c == GetAllTasks {
				GetAll(todoRepo)
			}
		}
	}()

	go func() {
		for c := range writeChannel {
			fmt.Println("Write channel:", c)
			if c == CreateTask {
				Create(todoRepo)
			}
		}
	}()

	// read/write to db depending on the channel
	for _, c := range commands {
		fmt.Println("Command:", c)
		if c == GetAllTasks {
			readChannel <- c
		} else if c == CreateTask {
			writeChannel <- c
		}
	}

	close(readChannel)
	close(writeChannel)

	for c := range readChannel {
		fmt.Println("Read channel:", c)
	}
	for c := range writeChannel {
		fmt.Println("Write channel:", c)
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
