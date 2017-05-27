package errors

import "net/http"

type validationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// InternalServerError creates a new api error representing an internal server error (HTTP 500)
func InternalServerError(err error) *APIError {
	return NewAPIError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", Params{"error": err.Error()})
}

// BadRequest creates a new api error representing a bad request (HTTP 403)
func BadRequest(err string) *APIError {
	return NewAPIError(http.StatusBadRequest, "BAD_REQUEST", Params{"error": err})
}
