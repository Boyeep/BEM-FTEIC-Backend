package apperr

import (
	"errors"
	"net/http"
)

type Error struct {
	Code       string
	Message    string
	HTTPStatus int
	Cause      error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(code, message string, status int) *Error {
	return &Error{Code: code, Message: message, HTTPStatus: status}
}

func Wrap(cause error, code, message string, status int) *Error {
	return &Error{Code: code, Message: message, HTTPStatus: status, Cause: cause}
}

func NotFound(resource string) *Error {
	return New("NOT_FOUND", resource+" not found", http.StatusNotFound)
}

func Validation(message string) *Error {
	return New("VALIDATION_ERROR", message, http.StatusBadRequest)
}

func Internal(cause error) *Error {
	return Wrap(cause, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
}

func As(err error) *Error {
	var appError *Error
	if errors.As(err, &appError) {
		return appError
	}
	return Internal(err)
}
