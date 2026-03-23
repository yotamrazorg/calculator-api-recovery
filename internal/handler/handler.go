// Package handler provides HTTP handlers and route registration for the calculator API.
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"calculator-api/internal/calculator"
	"calculator-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// readRawBody reads and restores the request body, returning it as a map.
func readRawBody(c *gin.Context) map[string]interface{} {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var raw map[string]interface{}
	_ = json.Unmarshal(body, &raw)
	return raw
}

// operations maps operation names to calculator functions.
var operations = map[string]func(float64, float64) (float64, error){
	"add": func(a, b float64) (float64, error) { return calculator.Add(a, b), nil },
	"sub": func(a, b float64) (float64, error) { return calculator.Subtract(a, b), nil },
	"mul": func(a, b float64) (float64, error) { return calculator.Multiply(a, b), nil },
	"div": calculator.Divide,
}

// RegisterRoutes sets up all API routes on the given Gin engine.
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/health", HealthHandler)
	r.POST("/add", AddHandler)
	r.POST("/subtract", SubtractHandler)
	r.POST("/multiply", MultiplyHandler)
	r.POST("/divide", DivideHandler)
	r.POST("/calculations", CreateCalculationHandler(db))
	r.GET("/calculations", ListCalculationsHandler(db))
	r.GET("/calculations/:id", GetCalculationByIDHandler(db))
	r.DELETE("/calculations/:id", DeleteCalculationByIDHandler(db))
}

// abortWithError sends a JSON error response matching FastAPI's {"detail": "..."} shape.
func abortWithError(c *gin.Context, code int, msg string) {
	c.JSON(code, model.ErrorResponse{Detail: msg})
}

// abortWithValidationError sends a 422 response matching FastAPI/Pydantic v2 validation error format.
func abortWithValidationError(c *gin.Context, err error, rawBody map[string]interface{}) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make([]model.ValidationErrorDetail, 0, len(ve))
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			details = append(details, model.ValidationErrorDetail{
				Loc:   []interface{}{"body", field},
				Msg:   "Field required",
				Type:  "missing",
				Input: rawBody,
			})
		}
		c.JSON(http.StatusUnprocessableEntity, model.ValidationErrorResponse{Detail: details})
		return
	}

	// JSON unmarshal / type errors
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		var inputVal interface{}
		if rawBody != nil {
			inputVal = rawBody[ute.Field]
		}
		details := []model.ValidationErrorDetail{
			{
				Loc:   []interface{}{"body", ute.Field},
				Msg:   "Input should be a valid number, unable to parse string as a number",
				Type:  "float_parsing",
				Input: inputVal,
			},
		}
		c.JSON(http.StatusUnprocessableEntity, model.ValidationErrorResponse{Detail: details})
		return
	}

	// Fallback
	details := []model.ValidationErrorDetail{
		{
			Loc:   []interface{}{"body"},
			Msg:   err.Error(),
			Type:  "value_error",
			Input: rawBody,
		},
	}
	c.JSON(http.StatusUnprocessableEntity, model.ValidationErrorResponse{Detail: details})
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
	rawBody := readRawBody(c)
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, err, rawBody)
		return
	}
	result := calculator.Add(*req.A, *req.B)
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// SubtractHandler handles POST /subtract.
func SubtractHandler(c *gin.Context) {
	rawBody := readRawBody(c)
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, err, rawBody)
		return
	}
	result := calculator.Subtract(*req.A, *req.B)
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// MultiplyHandler handles POST /multiply.
func MultiplyHandler(c *gin.Context) {
	rawBody := readRawBody(c)
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, err, rawBody)
		return
	}
	result := calculator.Multiply(*req.A, *req.B)
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// DivideHandler handles POST /divide.
func DivideHandler(c *gin.Context) {
	rawBody := readRawBody(c)
	var req model.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithValidationError(c, err, rawBody)
		return
	}
	result, err := calculator.Divide(*req.A, *req.B)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, model.ResultResponse{Result: result})
}

// CreateCalculationHandler handles POST /calculations.
func CreateCalculationHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawBody := readRawBody(c)
		var req model.CalculationCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			abortWithValidationError(c, err, rawBody)
			return
		}

		opFunc, ok := operations[req.Operation]
		if !ok {
			msg := fmt.Sprintf("Unknown operation: %s. Use: ['add', 'sub', 'mul', 'div']", req.Operation)
			abortWithError(c, http.StatusBadRequest, msg)
			return
		}

		result, err := opFunc(*req.A, *req.B)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, err.Error())
			return
		}

		calc := model.Calculation{
			Operation: req.Operation,
			A:         *req.A,
			B:         *req.B,
			Result:    result,
			CreatedAt: time.Now().UTC(),
		}
		if err := db.Create(&calc).Error; err != nil {
			abortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, calc.ToResponse())
	}
}

// ListCalculationsHandler handles GET /calculations.
func ListCalculationsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var calcs []model.Calculation
		if err := db.Order("created_at desc").Find(&calcs).Error; err != nil {
			abortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		responses := make([]model.CalculationResponse, len(calcs))
		for i, calc := range calcs {
			responses[i] = calc.ToResponse()
		}
		c.JSON(http.StatusOK, responses)
	}
}

// GetCalculationByIDHandler handles GET /calculations/:id.
func GetCalculationByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		var calc model.Calculation
		if err := db.First(&calc, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				abortWithError(c, http.StatusNotFound, "Calculation not found")
			} else {
				abortWithError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.JSON(http.StatusOK, calc.ToResponse())
	}
}

// DeleteCalculationByIDHandler handles DELETE /calculations/:id.
func DeleteCalculationByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		var calc model.Calculation
		if err := db.First(&calc, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				abortWithError(c, http.StatusNotFound, "Calculation not found")
			} else {
				abortWithError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		if err := db.Delete(&calc).Error; err != nil {
			abortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.Status(http.StatusNoContent)
	}
}
