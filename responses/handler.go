package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
	Data       any
}

func NewAPIError(statusCode int, message string, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
		Data:       err.Error(),
	}
}

func (e APIError) Error() string {
	return fmt.Sprintf("APIError: %d", e.StatusCode)
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "Invalid request data",
		Data:       errors,
	}
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, "Invalid JSON", fmt.Errorf("invalid JSON"))
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeResponse(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				WriteJson(w, apiErr.StatusCode, apiErr)
			} else {
				errResponse := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"message":    err.Error(),
				}
				WriteJson(w, http.StatusInternalServerError, errResponse)
			}
		}
	}
}

func WriteJson(w http.ResponseWriter, statusCode int, value any) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}
