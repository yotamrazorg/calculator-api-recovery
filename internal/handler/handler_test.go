package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
