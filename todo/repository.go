package todo

type Repository interface {
	Migrate() error
	Create(todo Todo) (*Todo, error)
	All() ([]Todo, error)
	GetByName(name string) (*Todo, error)
	Update(id int64, updated Todo) (*Todo, error)
	Delete(id int64) error
}
