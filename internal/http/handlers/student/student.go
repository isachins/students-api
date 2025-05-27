package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/isachins/students-api/internal/types"
	"github.com/isachins/students-api/internal/utils/response"
)

func New() http.HandlerFunc {
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


		response.WriteJSON(w, http.StatusCreated, map[string] string{"success": "OK"})
	}
}
