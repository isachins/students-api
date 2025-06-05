package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/isachins/students-api/internal/storage"
	"github.com/isachins/students-api/internal/types"
	"github.com/isachins/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// TODO: Request Validation

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastID, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("User Created successfully", slog.String("id", fmt.Sprint(lastID)))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusCreated, response.GeneralSuccess("Student Created Successfully", map[string]int64{"id": lastID}))
	}

}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("error parsing id", slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentByID(intID)
		if err != nil {
			slog.Error("error getting student", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, response.GeneralSuccess("Student List Fetched Successfully", student))
	}
}

func GetStudentsList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("error getting students", slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, response.GeneralSuccess("Student List Fetched Successfully", students))
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("updating student")

		// 1. Get the student ID from path parameter
		idStr := r.PathValue("id")
		if idStr == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format")))
			return
		}

		// 2. Get the existing student to update
		existingStudent, err := storage.GetStudentByID(id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student not found")))
			return
		}

		// 3. Parse the request body into a map to handle partial updates
		var updateData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
			return
		}

		// 4. Update only the fields that are present in the request
		if name, ok := updateData["name"].(string); ok {
			existingStudent.Name = name
		}
		if email, ok := updateData["email"].(string); ok {
			existingStudent.Email = email
		}
		if age, ok := updateData["age"].(float64); ok {
			existingStudent.Age = int(age)
		}

		// 5. Validate the updated student
		if err := validator.New().Struct(existingStudent); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// 6. Update the student
		updatedStudent, err := storage.UpdateStudent(
			id,
			existingStudent.Name,
			existingStudent.Email,
			existingStudent.Age,
		)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 7. Return the updated student
		response.WriteJSON(w, http.StatusOK, response.GeneralSuccess("Student Updated Successfully", updatedStudent))
	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")

		slog.Info("deleting a student: ", slog.String("id", strId))

		intID, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			slog.Error("error parsing id", slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.DeleteStudent(intID)
		if err != nil {
			slog.Error("error deleting student", slog.String("id", strId))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		var result = map[string]int64{"id": id}
		response.WriteJSON(w, http.StatusOK, response.GeneralSuccess("Successfully Deleted user", result))
	}
}
