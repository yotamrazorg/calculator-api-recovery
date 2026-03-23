package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"calculator-api/internal/database"
	"calculator-api/internal/model"

	"github.com/gin-gonic/gin"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Use a temporary SQLite database for tests.
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp db: %v", err)
	}
	tmpFile.Close()
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	db, err := database.NewDB(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to init test db: %v", err)
	}

	r := gin.New()
	RegisterRoutes(r, db)
	return r
}

func TestHealthEndpoint(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status = %q, want %q", resp.Status, "ok")
	}
	if resp.Version != "0.1.0" {
		t.Errorf("version = %q, want %q", resp.Version, "0.1.0")
	}
}

func TestAddEndpoint(t *testing.T) {
	r := setupRouter(t)

	body := `{"a": 3, "b": 5}`
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /add status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Result != 8 {
		t.Errorf("result = %v, want %v", resp.Result, 8.0)
	}
}

func TestSubtractEndpoint(t *testing.T) {
	r := setupRouter(t)

	body := `{"a": 10, "b": 4}`
	req := httptest.NewRequest(http.MethodPost, "/subtract", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /subtract status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Result != 6 {
		t.Errorf("result = %v, want %v", resp.Result, 6.0)
	}
}

func TestMultiplyEndpoint(t *testing.T) {
	r := setupRouter(t)

	body := `{"a": 3, "b": 7}`
	req := httptest.NewRequest(http.MethodPost, "/multiply", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /multiply status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Result != 21 {
		t.Errorf("result = %v, want %v", resp.Result, 21.0)
	}
}

func TestDivideEndpoint(t *testing.T) {
	r := setupRouter(t)

	body := `{"a": 15, "b": 3}`
	req := httptest.NewRequest(http.MethodPost, "/divide", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /divide status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Result != 5 {
		t.Errorf("result = %v, want %v", resp.Result, 5.0)
	}
}

func TestDivideByZeroEndpoint(t *testing.T) {
	r := setupRouter(t)

	body := `{"a": 5, "b": 0}`
	req := httptest.NewRequest(http.MethodPost, "/divide", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /divide (by zero) status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}
	if resp.Detail != "Cannot divide by zero" {
		t.Errorf("detail = %q, want %q", resp.Detail, "Cannot divide by zero")
	}
}

func TestAddEndpointWithFloats(t *testing.T) {
	r := setupRouter(t)

	body := `{"a": 1.5, "b": 2.5}`
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /add status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp model.ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp.Result != 4 {
		t.Errorf("result = %v, want %v", resp.Result, 4.0)
	}
}

func TestAddEndpointInvalidBody(t *testing.T) {
	r := setupRouter(t)

	body := `{"invalid": "json"}`
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// FastAPI returns 422 for validation errors; our Go service matches this behavior.
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("POST /add (invalid body) status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}

	var resp model.ValidationErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse validation error response: %v", err)
	}
	if len(resp.Detail) == 0 {
		t.Error("expected non-empty validation error details")
	}
}

// --- CRUD Endpoint Tests ---

// helper to create a calculation via POST /calculations and return the raw JSON map.
func createCalculation(t *testing.T, r *gin.Engine, operation string, a, b float64) map[string]interface{} {
	t.Helper()
	body := fmt.Sprintf(`{"operation": %q, "a": %v, "b": %v}`, operation, a, b)
	req := httptest.NewRequest(http.MethodPost, "/calculations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /calculations status = %d, want %d; body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return resp
}

func TestCreateCalculation(t *testing.T) {
	r := setupRouter(t)

	resp := createCalculation(t, r, "add", 5, 3)

	// Verify all expected fields are present.
	if resp["operation"] != "add" {
		t.Errorf("operation = %v, want %q", resp["operation"], "add")
	}
	if resp["a"] != 5.0 {
		t.Errorf("a = %v, want %v", resp["a"], 5.0)
	}
	if resp["b"] != 3.0 {
		t.Errorf("b = %v, want %v", resp["b"], 3.0)
	}
	if resp["result"] != 8.0 {
		t.Errorf("result = %v, want %v", resp["result"], 8.0)
	}
	if resp["id"] == nil || resp["id"].(float64) < 1 {
		t.Errorf("expected valid id, got %v", resp["id"])
	}
	// Verify created_at is in Python datetime format (no trailing Z).
	createdAt, ok := resp["created_at"].(string)
	if !ok || createdAt == "" {
		t.Fatalf("expected non-empty created_at string, got %v", resp["created_at"])
	}
	// Python format: "2024-01-15T10:30:00.123456" or "2024-01-15T10:30:00"
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?$`, createdAt)
	if !matched {
		t.Errorf("created_at = %q, does not match Python datetime format", createdAt)
	}
}

func TestCreateCalculationAllOperations(t *testing.T) {
	r := setupRouter(t)

	tests := []struct {
		op     string
		a, b   float64
		result float64
	}{
		{"add", 10, 5, 15},
		{"sub", 10, 5, 5},
		{"mul", 10, 5, 50},
		{"div", 10, 5, 2},
	}
	for _, tt := range tests {
		resp := createCalculation(t, r, tt.op, tt.a, tt.b)
		if resp["result"] != tt.result {
			t.Errorf("POST /calculations op=%s: result = %v, want %v", tt.op, resp["result"], tt.result)
		}
	}
}

func TestCreateCalculationUnknownOperation(t *testing.T) {
	r := setupRouter(t)

	body := `{"operation": "foo", "a": 5, "b": 3}`
	req := httptest.NewRequest(http.MethodPost, "/calculations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /calculations (unknown op) status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}
	expected := "Unknown operation: foo. Use: ['add', 'sub', 'mul', 'div']"
	if resp.Detail != expected {
		t.Errorf("detail = %q, want %q", resp.Detail, expected)
	}
}

func TestCreateCalculationDivideByZero(t *testing.T) {
	r := setupRouter(t)

	body := `{"operation": "div", "a": 1, "b": 0}`
	req := httptest.NewRequest(http.MethodPost, "/calculations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /calculations (div by zero) status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}
	if resp.Detail != "Cannot divide by zero" {
		t.Errorf("detail = %q, want %q", resp.Detail, "Cannot divide by zero")
	}
}

func TestListCalculations(t *testing.T) {
	r := setupRouter(t)

	// Create two calculations.
	createCalculation(t, r, "add", 1, 2)
	createCalculation(t, r, "sub", 10, 3)

	// List calculations.
	req := httptest.NewRequest(http.MethodGet, "/calculations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /calculations status = %d, want %d", w.Code, http.StatusOK)
	}

	var list []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 calculations, got %d", len(list))
	}

	// Verify ordered by created_at descending (newest first).
	// The second created calculation ("sub") should appear first.
	if list[0]["operation"] != "sub" {
		t.Errorf("first item operation = %v, want %q (newest first)", list[0]["operation"], "sub")
	}
	if list[1]["operation"] != "add" {
		t.Errorf("second item operation = %v, want %q", list[1]["operation"], "add")
	}
}

func TestListCalculationsEmpty(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/calculations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /calculations status = %d, want %d", w.Code, http.StatusOK)
	}

	// Should return an empty JSON array, not null.
	body := w.Body.String()
	if body != "[]" {
		t.Errorf("GET /calculations (empty) body = %q, want %q", body, "[]")
	}
}

func TestGetCalculationByID(t *testing.T) {
	r := setupRouter(t)

	created := createCalculation(t, r, "mul", 4, 5)
	id := int(created["id"].(float64))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/calculations/%d", id), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /calculations/%d status = %d, want %d", id, w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if resp["operation"] != "mul" {
		t.Errorf("operation = %v, want %q", resp["operation"], "mul")
	}
	if resp["result"] != 20.0 {
		t.Errorf("result = %v, want %v", resp["result"], 20.0)
	}
}

func TestGetCalculationByIDNotFound(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/calculations/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GET /calculations/9999 status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}
	if resp.Detail != "Calculation not found" {
		t.Errorf("detail = %q, want %q", resp.Detail, "Calculation not found")
	}
}

func TestDeleteCalculation(t *testing.T) {
	r := setupRouter(t)

	created := createCalculation(t, r, "add", 1, 1)
	id := int(created["id"].(float64))

	// Delete the calculation.
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/calculations/%d", id), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("DELETE /calculations/%d status = %d, want %d", id, w.Code, http.StatusNoContent)
	}
	if w.Body.Len() != 0 {
		t.Errorf("DELETE response body should be empty, got %q", w.Body.String())
	}

	// Verify the calculation is gone.
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/calculations/%d", id), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GET /calculations/%d after delete: status = %d, want %d", id, w.Code, http.StatusNotFound)
	}
}

func TestCreateCalculationMissingFields(t *testing.T) {
	r := setupRouter(t)

	body := `{"operation": "add"}`
	req := httptest.NewRequest(http.MethodPost, "/calculations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// FastAPI returns 422 for validation errors; our Go service matches this behavior.
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("POST /calculations (missing fields) status = %d, want %d; body = %s", w.Code, http.StatusUnprocessableEntity, w.Body.String())
	}

	var resp model.ValidationErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse validation error response: %v", err)
	}
	if len(resp.Detail) == 0 {
		t.Error("expected non-empty validation error details")
	}
}

func TestDeleteCalculationNotFound(t *testing.T) {
	r := setupRouter(t)

	req := httptest.NewRequest(http.MethodDelete, "/calculations/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("DELETE /calculations/9999 status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var resp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}
	if resp.Detail != "Calculation not found" {
		t.Errorf("detail = %q, want %q", resp.Detail, "Calculation not found")
	}
}
