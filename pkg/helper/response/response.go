package response

import (
	"encoding/json"
	"errors"
	"net/http"
)

// General errors
var (
	ErrInternalServerError = errors.New("internal server error")
)

// Database errors
var (
	ErrDatabaseQuery     = errors.New("database query failed")
	ErrRecordNotFound    = errors.New("record not found")
	ErrDuplicateRecord   = errors.New("duplicate record")
	ErrTransactionFailed = errors.New("transaction failed")
)

// Response represents the standard API response structure
type Response struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorData `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

// ErrorData represents error details
type ErrorData struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Meta represents pagination and additional metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
	TotalCount int64 `json:"total_count,omitempty"`
}

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Success writes a successful JSON response
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSON(w, statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created writes a 201 Created response
func Created(w http.ResponseWriter, message string, data interface{}) {
	Success(w, http.StatusCreated, message, data)
}

// OK writes a 200 OK response
func OK(w http.ResponseWriter, message string, data interface{}) {
	Success(w, http.StatusOK, message, data)
}

// NoContent writes a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, statusCode int, code string, message string, details map[string]string) {
	JSON(w, statusCode, Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest writes a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string, details map[string]string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// Unauthorized writes a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Forbidden writes a 403 Forbidden response
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFound writes a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", message, nil)
}

// Conflict writes a 409 Conflict response
func Conflict(w http.ResponseWriter, message string, details map[string]string) {
	Error(w, http.StatusConflict, "CONFLICT", message, details)
}

// UnprocessableEntity writes a 422 Unprocessable Entity response
func UnprocessableEntity(w http.ResponseWriter, message string, details map[string]string) {
	Error(w, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, details)
}

// InternalServerError writes a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}

// NotImplemented writes a 501 Not Implemented response
func NotImplemented(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", message, nil)
}

// SuccessWithMeta writes a successful JSON response with metadata (useful for pagination)
func SuccessWithMeta(w http.ResponseWriter, statusCode int, message string, data interface{}, meta *Meta) {
	JSON(w, statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// OKWithMeta writes a 200 OK response with metadata
func OKWithMeta(w http.ResponseWriter, message string, data interface{}, meta *Meta) {
	SuccessWithMeta(w, http.StatusOK, message, data, meta)
}
