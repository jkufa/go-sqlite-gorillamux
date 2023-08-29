package todo

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

// All errors that can be returned by the repository
var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

// Struct for Sqlite3 implementation
type SQLiteRepository struct {
	Repository
	db *sql.DB
}

// Constructor for Sqlite3 implementation
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			completed_date TIMESTAMP
			);`

	_, err := r.db.Exec(query) // Exec is used for queries that don't return rows
	return err
}

func (r *SQLiteRepository) Create(todo Todo) (*Todo, error) {
	query := `
		INSERT INTO todos (name, description, completed)
		VALUES (?, ?, ?);`
	res, err := r.db.Exec(query, todo.Name, todo.Description, todo.Completed) // Exec is used for queries that don't return rows
	if err != nil {
		var sqliteErr *sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	todo.Id = id

	return &todo, nil
}

func (r *SQLiteRepository) All() ([]Todo, error) {
	query := `
		SELECT id, name, description, completed, created_date, completed_date
		FROM todos;`
	rows, err := r.db.Query(query) // Query is used for rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Completed, &todo.CreatedDate, &todo.CompletedDate); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *SQLiteRepository) GetByName(name string) ([]Todo, error) {
	query := `
		SELECT id, name, description, completed, created_date, completed_date
		FROM todos
		WHERE name = ?;`
	rows, err := r.db.Query(query, name) // Query is used for rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Completed, &todo.CreatedDate, &todo.CompletedDate); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
