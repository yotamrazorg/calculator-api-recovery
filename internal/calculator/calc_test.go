package calculator

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{1, 2, 3},
		{-1, 1, 0},
		{0, 0, 0},
		{1.5, 2.5, 4},
		{-3.5, -2.5, -6},
	}
	for _, tt := range tests {
		got := Add(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{5, 3, 2},
		{3, 5, -2},
		{0, 0, 0},
		{1.5, 0.5, 1},
	}
	for _, tt := range tests {
		got := Subtract(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{2, 3, 6},
		{-2, 3, -6},
		{0, 100, 0},
		{1.5, 2, 3},
	}
	for _, tt := range tests {
		got := Multiply(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{6, 3, 2},
		{7, 2, 3.5},
		{-6, 3, -2},
		{0, 5, 0},
	}
	for _, tt := range tests {
		got, err := Divide(tt.a, tt.b)
		if err != nil {
			t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
			continue
		}
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(5, 0)
	if err == nil {
		t.Fatal("Divide(5, 0) expected error, got nil")
	}
	if err != ErrDivideByZero {
		t.Errorf("Divide(5, 0) error = %v, want ErrDivideByZero", err)
	}
	if err.Error() != "Cannot divide by zero" {
		t.Errorf("error message = %q, want %q", err.Error(), "Cannot divide by zero")
	}
}
