package xivmath

import "math"

// Contains functions for calculating specific values for the game.
// Most formulas are taken from Allagan Studies

func P_DHR(DHR int, LvSub int, LvDiv int) float64 {
	return math.Floor(550*(float64(DHR-LvSub))/float64(LvDiv)) / 1000
}

func P_CHR(CHR int, LvSub int, LvDiv int) float64 {
	return math.Floor(200*(float64(CHR-LvSub))/float64(LvDiv)+50) / 1000
}

// A Character regenerates 200 (2%) MP per Tick in Combat.
// Healers can increase this value with Piety.
func P_MPT(PIE int, LvMain int, LvDiv int) float64 {
	return math.Floor(150*(float64(PIE-LvMain))/float64(LvDiv) + 200)
}
