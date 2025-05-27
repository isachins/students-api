package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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

		lastID, err :=storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("User Created successfully", slog.String("id", fmt.Sprint(lastID)))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}


		response.WriteJSON(w, http.StatusCreated, map[string] int64{"id": lastID})
	}
}
