package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

// WriteJSON writes the data to the response writer with the given status code
// and content type.
// It returns an error if the data cannot be encoded to JSON.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
