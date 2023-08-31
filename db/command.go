// TODO: move this to a separate package, something related to messaging

package main

type Command int64

const (
	Unknown Command = iota
	GetTasksByName
	GetAllTasks
	CreateTask
	UpdateTask
	DeleteTask
)
