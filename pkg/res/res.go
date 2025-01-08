package res

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error" example:"message"`
}

// BadRequestError represents a 400 Bad Request error
type BadRequestError struct {
	Error string `json:"error" example:"Invalid request: Please check your input"`
}

// UnauthorizedError represents a 401 Unauthorized error
type UnauthorizedError struct {
	Error string `json:"error" example:"Unauthorized: Authentication required"`
}

// PaymentRequiredError represents a 402 Payment Required error
type PaymentRequiredError struct {
	Error string `json:"error" example:"Payment required to access this resource"`
}

// ForbiddenError represents a 403 Forbidden error
type ForbiddenError struct {
	Error string `json:"error" example:"Forbidden: You don't have permission to access this resource"`
}

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	Error string `json:"error" example:"Resource not found"`
}

// MethodNotAllowedError represents a 405 Method Not Allowed error
type MethodNotAllowedError struct {
	Error string `json:"error" example:"Method not allowed for this endpoint"`
}

// ConflictError represents a 409 Conflict error
type ConflictError struct {
	Error string `json:"error" example:"Resource conflict: The request could not be completed due to a conflict"`
}

// InternalServerError represents a 500 Internal Server Error
type InternalServerError struct {
	Error string `json:"error" example:"Internal server error occurred"`
}

func Error(w http.ResponseWriter, error string, code int) {
	Success(w, ErrorResponse{Error: error}, code)
}

func Success(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
