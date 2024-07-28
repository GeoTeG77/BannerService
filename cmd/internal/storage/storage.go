package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists = errors.New("url exists")
	ErrTokenIsNotValid = errors.New("invalid jwt token")

	
	StatusOK    = "OK"
	StatusCreate = "Create"
	StatusBadRequets = "Bad Request"
	StatusForbidden = "Forbidden"
	StatusNotFound = "Not Found"
	//StatusError = "Error"
	StatusInternalServerError = "Internal Server Error"

)
