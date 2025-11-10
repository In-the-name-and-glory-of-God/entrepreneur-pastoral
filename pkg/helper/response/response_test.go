package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	OK(w, "Success message", data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("Expected success to be true")
	}

	if resp.Message != "Success message" {
		t.Errorf("Expected message 'Success message', got '%s'", resp.Message)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"id": "123"}

	Created(w, "Resource created", data)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("Expected success to be true")
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()

	NoContent(w)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Error("Expected empty body for NoContent response")
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	details := map[string]string{
		"email": "Email is required",
	}

	BadRequest(w, "Validation failed", details)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}

	if resp.Error == nil {
		t.Fatal("Expected error to be present")
	}

	if resp.Error.Code != "BAD_REQUEST" {
		t.Errorf("Expected error code 'BAD_REQUEST', got '%s'", resp.Error.Code)
	}

	if resp.Error.Details["email"] != "Email is required" {
		t.Error("Expected validation details to be included")
	}
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()

	Unauthorized(w, "Invalid credentials")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}

	if resp.Error.Code != "UNAUTHORIZED" {
		t.Errorf("Expected error code 'UNAUTHORIZED', got '%s'", resp.Error.Code)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	NotFound(w, "Resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	details := map[string]string{
		"email": "Email already exists",
	}

	Conflict(w, "Resource already exists", details)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}

	if resp.Error.Code != "CONFLICT" {
		t.Errorf("Expected error code 'CONFLICT', got '%s'", resp.Error.Code)
	}
}

func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()

	InternalServerError(w, "Something went wrong")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Success {
		t.Error("Expected success to be false")
	}
}

func TestOKWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	data := []map[string]string{
		{"id": "1", "name": "Item 1"},
		{"id": "2", "name": "Item 2"},
	}
	meta := &Meta{
		Page:       1,
		PageSize:   10,
		TotalPages: 5,
		TotalCount: 50,
	}

	OKWithMeta(w, "Items retrieved", data, meta)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Error("Expected success to be true")
	}

	if resp.Meta == nil {
		t.Fatal("Expected meta to be present")
	}

	if resp.Meta.Page != 1 {
		t.Errorf("Expected page 1, got %d", resp.Meta.Page)
	}

	if resp.Meta.TotalCount != 50 {
		t.Errorf("Expected total count 50, got %d", resp.Meta.TotalCount)
	}
}

func TestContentTypeHeader(t *testing.T) {
	w := httptest.NewRecorder()

	OK(w, "Test", nil)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}
