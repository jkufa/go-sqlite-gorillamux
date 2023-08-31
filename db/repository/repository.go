package repository

import (
	"kufa.io/sqlitego/db/models"
)

type Repository interface {
	Migrate() error
	Create(todo models.Task) (*models.Task, error)
	All() ([]models.Task, error)
	GetByName(name string) (*models.Task, error)
	Update(id int64, updated models.Task) (*models.Task, error)
	Delete(id int64) error
}
