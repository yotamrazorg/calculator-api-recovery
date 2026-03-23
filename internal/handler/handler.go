// Package handler provides HTTP handlers and route registration for the calculator API.
package handler

import (
	"net/http"

	"calculator-api/internal/calculator"
	"calculator-api/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes sets up all API routes on the given Gin engine.
// The db parameter is used by CRUD handlers (registered in Milestone 2).
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/health", HealthHandler)
	r.POST("/add", AddHandler)
	r.POST("/subtract", SubtractHandler)
	r.POST("/multiply", MultiplyHandler)
	r.POST("/divide", DivideHandler)
}

// abortWithError sends a JSON error response matching FastAPI's {"detail": "..."} shape.
func abortWithError(c *gin.Context, code int, msg string) {
	c.JSON(code, model.ErrorResponse{Detail: msg})
}

// HealthHandler returns the service health status.
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.HealthResponse{
		Status:  "ok",
		Version: "0.1.0",
	})
}

// AddHandler handles POST /add.
func AddHandler(c *gin.Context) {
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	result := calculator.Add(*req.A, *req.B)
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// SubtractHandler handles POST /subtract.
func SubtractHandler(c *gin.Context) {
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	result := calculator.Subtract(*req.A, *req.B)
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// MultiplyHandler handles POST /multiply.
func MultiplyHandler(c *gin.Context) {
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	result := calculator.Multiply(*req.A, *req.B)
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// DivideHandler handles POST /divide.
func DivideHandler(c *gin.Context) {
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := calculator.Divide(*req.A, *req.B)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}
