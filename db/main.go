package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"kufa.io/sqlitego/db/models"

	"kufa.io/sqlitego/db/repository"

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

	r := repository.NewSQLiteRepository(db)
	if err := r.Migrate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	writeChannel := make(chan Command, 1)
	readChannel := make(chan Command, 1)

	go func() {
		for c := range readChannel {
			fmt.Println("Read channel:", c)
			if c == GetAllTasks {
				GetAll(r)
			}
		}
	}()

	go func() {
		for c := range writeChannel {
			fmt.Println("Write channel:", c)
			if c == CreateTask {
				Create(r)
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

func GetAll(r *repository.SQLiteRepository) {
	// Get all tasks
	tasks, err := r.All()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All tasks:", tasks)
}

func Create(r *repository.SQLiteRepository) *models.Task {
	// Create a new task
	task1 := models.Task{
		Name:        "Learn Go",
		Description: "Learn Golang by solving algorithms and writing code!",
		Completed:   false,
	}
	created, err := r.Create(task1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created task:", created)
	return created
}
