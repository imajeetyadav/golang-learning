package sqlite

import (
	"database/sql"
	"demo/internal/config"
	"demo/internal/types"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type SQLiteStorage struct {
	Db *sql.DB
}

func (s *SQLiteStorage) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DeleteStudent implements storage.Storage.
func (s *SQLiteStorage) DeleteStudent(id int64) error {
	_, err := s.Db.Exec("DELETE FROM students WHERE id = ?", id)
	return err
}

// GetAllStudents implements storage.Storage.
func (s *SQLiteStorage) GetAllStudents() ([]types.Student, error) {
	rows, err := s.Db.Query("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// GetStudentByID implements storage.Storage.
func (s *SQLiteStorage) GetStudentByID(id int64) (*types.Student, error) {
	var student types.Student
	err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No student found
		}
		return nil, err
	}
	return &student, nil
}

// UpdateStudent implements storage.Storage.
func (s *SQLiteStorage) UpdateStudent(id int64, name string, email string, age int) error {
	_, err := s.Db.Exec("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?", name, email, age, id)
	return err
}

func New(cfg *config.Config) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	// Create tables if they do not exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL CHECK(age >= 0 AND age <= 130)
	);
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{Db: db}, nil
}
