package sqlite

import (
	"database/sql"
	"demo/internal/config"
	"demo/internal/types"
	"fmt"
	"log/slog"

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

func (s *SQLiteStorage) UpdateStudent(id int64, name string, email string, age int) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("no student found with ID %d", id)
	}
	slog.Info("Student updated successfully", slog.Int64("id", id))
	return nil
}

func (s *SQLiteStorage) DeleteStudent(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("no student found with ID %d", id)
	}
	slog.Info("Student deleted successfully", slog.Int64("id", id))
	return nil
}

func (s *SQLiteStorage) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to get all students: %w", err)
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *SQLiteStorage) GetStudentByID(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with ID %d not found", id)
		}
		return types.Student{}, fmt.Errorf("failed to get student by ID: %w", err)
	}
	return student, nil

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
