package students

import (
	"demo/internal/storage"
	"demo/internal/types"
	"demo/internal/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Received request", slog.String("method", r.Method), slog.String("url", r.URL.String()))

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(&student); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Error("Failed to create student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Student created successfully", slog.Int64("id", lastId))
		response.WriteJson(w, http.StatusCreated, map[string]int64{
			"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Received request", slog.String("method", r.Method), slog.String("url", r.URL.String()))

		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return

		}

		IntId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format: %s", id)))
			return
		}
		students, err := storage.GetStudentByID(IntId)
		if err != nil {
			slog.Error("Failed to get student by ID", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student with ID %s not found", id)))
			return
		}

		slog.Info("Student retrieved successfully", slog.Int64("id", students.Id))
		response.WriteJson(w, http.StatusOK, students)
	}
}
