package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)


type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK = "OK"
	StatusCreate = "Create"
	StatusBadRequets = "Bad Request"
	StatusUnathorized = "Unathorized"
	StatusForbidden = "Forbidden"
	StatusNotFound = "Not Found"
	StatusInternalServerError = "Internal Server Error"
	
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Create() Response {
	return Response{
		Status: StatusCreate,
	}
}

func BadRequest() Response {
	return Response{
		Status: StatusBadRequets,
	}
}

func Forbidden() Response {
	return Response{
		Status: StatusForbidden,
	}
}

func NotFound() Response {
	return Response{
		Status: StatusNotFound,
	}
}

func InternalServerError() Response {
	return Response{
		Status: StatusInternalServerError,
	}
}

func Error(save string) Response {
	return Response{
		Error: StatusNotFound,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusOK,
		Error:  strings.Join(errMsgs, ", "),
	}
}
