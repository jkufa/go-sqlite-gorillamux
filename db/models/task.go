package models

import (
	"time"
)

type Task struct {
	Id            int64
	Name          string
	Description   string
	Completed     bool
	CreatedDate   time.Time
	CompletedDate time.Time
}
