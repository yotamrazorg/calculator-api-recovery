// Package model defines database models and API request/response DTOs.
package model

import "time"

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
	ID        int       `json:"id"`
	Operation string    `json:"operation"`
	A         float64   `json:"a"`
	B         float64   `json:"b"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a Calculation model to a CalculationResponse DTO.
func (c *Calculation) ToResponse() CalculationResponse {
	return CalculationResponse{
		ID:        c.ID,
		Operation: c.Operation,
		A:         c.A,
		B:         c.B,
		Result:    c.Result,
		CreatedAt: c.CreatedAt,
	}
}
