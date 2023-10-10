package xivmath

import (
	"testing"
)

const (
	// LvMain is the amount of Base Mainstat at a specific level (390 for level 90)
	LvMain = 390
	// LvSub is the amount of Base Substat at a specific level (400 for level 90)
	LvSub = 400
	// LvDiv is the divisor for the substat formula (1900 for level 90)
	LvDiv = 1900
)

func TestF_WD(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		WD       int
		JobMod   int
		expected float64
	}{
		{
			name:     "F_WD Case 1",
			WD:       132,
			JobMod:   110, // 110 for tanks
			expected: 174,
		},
		{
			name:     "F_WD Case 2",
			WD:       131,
			JobMod:   110, // 110 for tanks
			expected: 173,
		},
		{
			name:     "F_WD Case 3",
			WD:       132,
			JobMod:   115,
			expected: 176,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := F_WD(tt.WD, LvMain, tt.JobMod)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestF_AP(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		M        int
		MainStat int
		expected float64
	}{
		{
			name:     "F_AP Case 1",
			M:        156,
			MainStat: 3330,
			expected: 1276,
		},
		{
			name:     "F_AP Case 2",
			M:        156,
			MainStat: 1000,
			expected: 344,
		},
		{
			name:     "F_AP Case 3",
			M:        156,
			MainStat: 3500,
			expected: 1344,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := F_AP(tt.M, tt.MainStat, LvMain)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
func TestF_DET(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		DET      int
		expected float64
	}{
		{
			name:     "F_DET Case 1",
			DET:      2182,
			expected: 1122,
		},
		{
			name:     "F_DET Case 2",
			DET:      900,
			expected: 1034,
		},
		{
			name:     "F_DET Case 3",
			DET:      7432,
			expected: 1481,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := F_DET(tt.DET, LvMain, LvDiv)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
func TestF_TNC(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		TNC      int
		expected float64
	}{
		{
			name:     "F_TNC Case 1",
			TNC:      529,
			expected: 1006,
		},
		{
			name:     "F_TNC Case 2",
			TNC:      400,
			expected: 1000,
		},
		{
			name:     "F_TNC Case 3",
			TNC:      2180,
			expected: 1093,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := F_TNC(tt.TNC, LvMain, LvSub, LvDiv)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
func TestF_SPD(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		SPD      int
		expected float64
	}{
		{
			name:     "F_SPD Case 1",
			SPD:      400,
			expected: 1000,
		},
		{
			name:     "F_SPD Case 2",
			SPD:      1182,
			expected: 1053,
		},
		{
			name:     "F_SPD Case 3",
			SPD:      2200,
			expected: 1123,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := F_SPD(tt.SPD, LvSub, LvDiv)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
func TestF_CRIT(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		CRIT     int
		expected float64
	}{
		{
			name:     "F_CRIT Case 1",
			CRIT:     2576,
			expected: 1629,
		},
		{
			name:     "F_CRIT Case 2",
			CRIT:     1000,
			expected: 1463,
		},
		{
			name:     "F_CRIT Case 3",
			CRIT:     400,
			expected: 1400,
		},
	}

	// Test each case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := F_CRIT(tt.CRIT, LvSub, LvDiv)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
