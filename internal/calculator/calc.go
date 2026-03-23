// Package calculator provides pure arithmetic functions.
package calculator

import "errors"

// ErrDivideByZero is returned when attempting to divide by zero.
var ErrDivideByZero = errors.New("Cannot divide by zero")

// Add returns the sum of a and b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns a minus b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of a and b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns a divided by b. Returns an error if b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
