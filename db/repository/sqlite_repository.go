package repository

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"kufa.io/sqlitego/db/models"
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
	db *sql.DB
}

// Constructor for Sqlite3 implementation
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
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

func (r *SQLiteRepository) Create(task models.Task) (*models.Task, error) {
	query := `
		INSERT INTO tasks (name, description, completed)
		VALUES (?, ?, ?);`
	res, err := r.db.Exec(query, task.Name, task.Description, false) // Exec is used for queries that don't return rows
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
	task.Id = id

	return &task, nil
}

func (r *SQLiteRepository) All() ([]models.Task, error) {
	query := `
		SELECT id, name, description, completed, created_date, completed_date
		FROM tasks;`
	rows, err := r.db.Query(query) // Query is used for rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	var completedDate sql.NullTime
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.Completed, &task.CreatedDate, &completedDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *SQLiteRepository) Get(id int64) (models.Task, error) {
	query := `
		SELECT id, name, description, completed, created_date, completed_date
		FROM tasks
		WHERE id = ?;`
	row := r.db.QueryRow(query, id) // QueryRow is used for a single row
	var task models.Task
	var completedDate sql.NullTime
	if err := row.Scan(&task.Id, &task.Name, &task.Description, &task.Completed, &task.CreatedDate, &completedDate); err != nil {
		return task, err
	}

	return task, nil
}

func (r *SQLiteRepository) GetByName(name string) ([]models.Task, error) {
	query := `
		SELECT id, name, description, completed, created_date, completed_date
		FROM tasks
		WHERE name = ?;`
	rows, err := r.db.Query(query, name) // Query is used for rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.Completed, &task.CreatedDate, &task.CompletedDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *SQLiteRepository) Update(task models.Task) (*models.Task, error) {
	query := `
		UPDATE tasks
		SET name = ?, description = ?, completed = ?, completed_date = ?
		WHERE id = ?;`
	res, err := r.db.Exec(query, task.Name, task.Description, task.Completed, task.CompletedDate, task.Id) // Exec is used for queries that don't return rows
	if err != nil {
		return nil, ErrUpdateFailed
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrNotExists
	}

	return &task, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	query := `
		DELETE FROM tasks
		WHERE id = ?;`
	_, err := r.db.Exec(query, id) // Exec is used for queries that don't return rows
	return err
}
