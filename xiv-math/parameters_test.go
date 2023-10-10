package xivmath

import (
	"math"
	"testing"
)

// TODO: Refactor where LvMain... values come from
// LvMain, LvSub and LvDiv are already in functions_test.go

func TestP_DHR(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		DHR      int
		expected float64
	}{
		{
			name:     "P_DHR Case 1",
			DHR:      940,
			expected: 15.6,
		},
		{
			name:     "P_DHR Case 2",
			DHR:      400,
			expected: 0,
		},
		{
			name:     "P_DHR Case 3",
			DHR:      2000,
			expected: 46.3,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := P_DHR(tt.DHR, LvSub, LvDiv)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
func TestP_CHR(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		CHR      int
		expected float64
	}{
		{
			name:     "P_CHR Case 1",
			CHR:      2576,
			expected: 27.9,
		},
		{
			name:     "P_CHR Case 2",
			CHR:      400,
			expected: 5,
		},
		{
			name:     "P_CHR Case 3",
			CHR:      1250,
			expected: 13.9,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := P_CHR(tt.CHR, LvSub, LvDiv)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestP_MPT(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		PIE      int
		expected float64
	}{
		{
			name:     "P_MPT Case 1",
			PIE:      390,
			expected: 200,
		},
		{
			name:     "P_MPT Case 2",
			PIE:      700,
			expected: 224,
		},
		{
			name:     "P_MPT Case 3",
			PIE:      1200,
			expected: 263,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := P_MPT(tt.PIE, LvMain, LvDiv)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
