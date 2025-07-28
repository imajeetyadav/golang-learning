package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
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
	var errMessages []string

	for _, err := range errs {

		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, "Field "+fmt.Sprintf("%s", err.Field())+" is required")
		case "email":
			errMessages = append(errMessages, "Field "+fmt.Sprintf("%s", err.Field())+" must be a valid email address")
		case "gte":
			errMessages = append(errMessages, "Field "+fmt.Sprintf("%s", err.Field())+" must be greater than or equal to "+err.Param())
		case "lte":
			errMessages = append(errMessages, "Field "+fmt.Sprintf("%s", err.Field())+" must be less than or equal to "+err.Param())
		default:
			errMessages = append(errMessages, "Field "+fmt.Sprintf("%s", err.Field())+" failed validation for tag "+err.ActualTag())
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMessages, ", "),
	}
}
