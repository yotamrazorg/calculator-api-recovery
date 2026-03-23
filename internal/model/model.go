// Package model defines database models and API request/response DTOs.
package model

import (
	"fmt"
	"strings"
	"time"
)

// PythonTime wraps time.Time to serialize JSON in Python's datetime format
// (ISO 8601 without trailing Z, with microsecond precision).
// Example: "2024-01-15T10:30:00.123456"
type PythonTime time.Time

// MarshalJSON formats time to match Python's datetime.isoformat() output.
func (pt PythonTime) MarshalJSON() ([]byte, error) {
	t := time.Time(pt)
	// Format: 2006-01-02T15:04:05.000000 (microsecond precision, no timezone)
	s := t.UTC().Format("2006-01-02T15:04:05.000000")
	// Trim trailing zeros after the decimal point, but keep at least one digit
	// Python trims to microsecond precision but keeps trailing zeros
	// Actually Python keeps all 6 digits: datetime(2024,1,15,10,30,0,0).isoformat() => "2024-01-15T10:30:00"
	// datetime(2024,1,15,10,30,0,123456).isoformat() => "2024-01-15T10:30:00.123456"
	// datetime(2024,1,15,10,30,0,100000).isoformat() => "2024-01-15T10:30:00.100000"
	// So Python omits the fractional part entirely if microseconds == 0, otherwise shows 6 digits.
	if t.Nanosecond() == 0 {
		s = t.UTC().Format("2006-01-02T15:04:05")
	}
	return []byte(fmt.Sprintf("%q", s)), nil
}

// UnmarshalJSON parses Python-style datetime strings.
func (pt *PythonTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	// Try with fractional seconds first, then without
	for _, layout := range []string{
		"2006-01-02T15:04:05.000000",
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
		time.RFC3339,
	} {
		if t, err := time.Parse(layout, s); err == nil {
			*pt = PythonTime(t)
			return nil
		}
	}
	// Fallback: try parsing with time.Parse for any ISO-like format
	t, err := time.Parse("2006-01-02T15:04:05.999999999", s)
	if err != nil {
		// Try stripping Z suffix
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", s)
		if err != nil {
			return fmt.Errorf("cannot parse %q as PythonTime", s)
		}
	}
	*pt = PythonTime(t)
	return nil
}

// --- API Request/Response DTOs ---

// CalculationRequest is the request body for arithmetic endpoints (POST /add, /subtract, /multiply, /divide).
type CalculationRequest struct {
	A *float64 `json:"a" binding:"required"`
	B *float64 `json:"b" binding:"required"`
}

// ResultResponse is the response body for arithmetic endpoints.
type ResultResponse struct {
	Result float64 `json:"result"`
}

// HealthResponse is the response body for the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ErrorResponse matches FastAPI's HTTPException JSON shape: {"detail": "..."}.
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// ValidationErrorDetail represents a single validation error in FastAPI/Pydantic v2 format.
type ValidationErrorDetail struct {
	Loc   []interface{} `json:"loc"`
	Msg   string        `json:"msg"`
	Type  string        `json:"type"`
	Input interface{}   `json:"input"`
}

// ValidationErrorResponse matches FastAPI's 422 response: {"detail": [...]}.
type ValidationErrorResponse struct {
	Detail []ValidationErrorDetail `json:"detail"`
}

// --- Database Model ---

// Calculation is the GORM model representing a stored calculation.
type Calculation struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Operation string    `json:"operation" gorm:"not null"`
	A         float64   `json:"a" gorm:"not null"`
	B         float64   `json:"b" gorm:"not null"`
	Result    float64   `json:"result" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// CalculationCreate is the request body for creating a calculation (POST /calculations).
type CalculationCreate struct {
	Operation string   `json:"operation" binding:"required"`
	A         *float64 `json:"a" binding:"required"`
	B         *float64 `json:"b" binding:"required"`
}

// CalculationResponse is the response body for calculation CRUD endpoints.
type CalculationResponse struct {
	ID        int        `json:"id"`
	Operation string     `json:"operation"`
	A         float64    `json:"a"`
	B         float64    `json:"b"`
	Result    float64    `json:"result"`
	CreatedAt PythonTime `json:"created_at"`
}

// ToResponse converts a Calculation model to a CalculationResponse DTO.
func (c *Calculation) ToResponse() CalculationResponse {
	return CalculationResponse{
		ID:        c.ID,
		Operation: c.Operation,
		A:         c.A,
		B:         c.B,
		Result:    c.Result,
		CreatedAt: PythonTime(c.CreatedAt),
	}
}
