package models

import "testing"

type CritTier struct {
	critAmount int
	critRate   float64
	critMod    float64
}

func SetupCritTiers() []CritTier {
	cts := []CritTier{
		{400, 0.05, 1.4},
		{410, 0.051, 1.401},
		{419, 0.052, 1.402},
		{429, 0.053, 1.403},
		{1065, 0.12, 1.47},
	}
	return cts
}

func TestCalculateCriticalHit(t *testing.T) {
	cts := SetupCritTiers()

	for _, ct := range cts {
		ps := NewPlayerStats(0, 0, ct.critAmount, 0, 0, 0, 0, 0, 0, 0)

		if ps.CriticalHitPercent != ct.critRate {
			t.Errorf("Expected %f, got %f", ct.critRate, ps.CriticalHitPercent)
		}
		if ps.CriticalHitMod != ct.critMod {
			t.Errorf("Expected %f, got %f", ct.critMod, ps.CriticalHitMod)
		}
	}
}

type DetTier struct {
	detAmount int
	detMod    float64
}

func SetupDetTiers() []DetTier {
	dts := []DetTier{
		{390, 1.0},
		{404, 1.001},
		{418, 1.002},
		{1571, 1.087},
	}
	return dts
}

func TestCalculateDetermination(t *testing.T) {
	dts := SetupDetTiers()

	for _, dt := range dts {
		ps := NewPlayerStats(0, 0, 0, 0, dt.detAmount, 0, 0, 0, 0, 0)

		if ps.DeterminationMod != dt.detMod {
			t.Errorf("Expected %f, got %f", dt.detMod, ps.DeterminationMod)
		}
	}
}

type TenTier struct {
	tenAmount int
	tenMod    float64
}

func SetupTenTiers() []TenTier {
	tts := []TenTier{
		{400, 1.0},
		{419, 1.001},
		{438, 1.002},
		{2319, 1.101},
	}
	return tts
}

func TestCalculateTenacity(t *testing.T) {
	tts := SetupTenTiers()

	for _, tt := range tts {
		ps := NewPlayerStats(0, 0, 0, 0, 0, tt.tenAmount, 0, 0, 0, 0)

		if ps.TenacityMod != tt.tenMod {
			t.Errorf("Expected %f, got %f", tt.tenMod, ps.TenacityMod)
		}
	}
}
