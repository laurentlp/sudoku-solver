package errors

import "net/http"

type validationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// BadRequest creates a new api error representing a bad request (HTTP 400)
func BadRequest(err string) *APIError {
	return NewAPIError(http.StatusBadRequest, "BAD_REQUEST", Params{"error": err})
}
