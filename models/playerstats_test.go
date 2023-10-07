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

type SpdTier struct {
	spdAmount int
	spdMod    float64
	spd1_5    int
	spd2      int
	spd2_5    int
	spd2_8    int
	spd3      int
	spd3_5    int
	spd4      int
}

func SetupSpdTiers() []SpdTier {
	sts := []SpdTier{
		{400, 1.0, 1500, 2000, 2500, 2800, 3000, 3500, 4000},
		{415, 1.001, 1490, 1990, 2490, 2790, 2990, 3490, 3990},
		{941, 1.037, 1440, 1920, 2400, 2690, 2880, 3370, 3850},
		{2637, 1.153, 1270, 1690, 2110, 2370, 2540, 2960, 3380},
	}
	return sts
}

func TestCalculateSpeed(t *testing.T) {
	sts := SetupSpdTiers()

	for _, st := range sts {
		ps := NewPlayerStats(0, 0, 0, 0, 0, 0, st.spdAmount, st.spdAmount, 0, 0)

		if ps.Speed.AA_DOT_MOD != st.spdMod {
			t.Errorf("Expected %f, got %f", st.spdMod, ps.Speed.AA_DOT_MOD)
		}
		if ps.Speed.GCD1_5 != st.spd1_5 {
			t.Errorf("Expected %d, got %d", st.spd1_5, ps.Speed.GCD1_5)
		}
		if ps.Speed.GCD2 != st.spd2 {
			t.Errorf("Expected %d, got %d", st.spd2, ps.Speed.GCD2)
		}
		if ps.Speed.GCD2_5 != st.spd2_5 {
			t.Errorf("Expected %d, got %d", st.spd2_5, ps.Speed.GCD2_5)
		}
		if ps.Speed.GCD2_8 != st.spd2_8 {
			t.Errorf("Expected %d, got %d", st.spd2_8, ps.Speed.GCD2_8)
		}
		if ps.Speed.GCD3 != st.spd3 {
			t.Errorf("Expected %d, got %d", st.spd3, ps.Speed.GCD3)
		}
		if ps.Speed.GCD3_5 != st.spd3_5 {
			t.Errorf("Expected %d, got %d", st.spd3_5, ps.Speed.GCD3_5)
		}
		if ps.Speed.GCD4 != st.spd4 {
			t.Errorf("Expected %d, got %d", st.spd4, ps.Speed.GCD4)
		}
	}
}
