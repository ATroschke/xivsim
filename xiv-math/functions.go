package xivmath

import (
	"math"
)

// Contains functions for calculating specific modifiers/values for the game.
// Most formulas are taken from Allagan Studies

func F_WD(WD int, LvMain int, JobMod int) float64 {
	return math.Floor((float64(LvMain) * float64(JobMod) / 1000.0) + float64(WD))
}

func F_AP(M int, MainStat int, LvMain int) float64 {
	return math.Floor(float64(M)*(float64(MainStat)-float64(LvMain))/float64(LvMain)) + 100
}

func F_DET(DET int, LvMain int, LvDiv int) float64 {
	return math.Floor(140*(float64(DET-LvMain))/float64(LvDiv) + 1000)
}

func F_TNC(TNC int, LvMain int, LvSub int, LvDiv int) float64 {
	return math.Floor(100*(float64(TNC-LvSub))/float64(LvDiv) + 1000)
}

func F_SPD(SPD int, LvSub int, LvDiv int) float64 {
	return math.Floor(130*(float64(SPD-LvSub))/float64(LvDiv) + 1000)
}

func F_CRIT(CRIT int, LvSub int, LvDiv int) float64 {
	return math.Floor(200*(float64(CRIT-LvSub))/float64(LvDiv)+1400) / 1000
}

func F_AUTO(LvMain int, JobMod int, WD int, WDelay float64) float64 {
	return math.Floor((math.Floor((float64(LvMain) * float64(JobMod) / 1000) + float64(WD))) * (WDelay / 3.0))
}

// Calculates the damage of an attack with the given Potency and Stats.
// This will be used to pre-bake the damage of a skill for a player without Crit, Direct Hit and Buffs to speed up the calculation.
func DirectDamage(Potency int, WD int, M int, MainStat int, LvMain int, JobMod int, DET int, LvDiv int, TNC int, SPD int, CRIT int, LvSub int, TraitMod int) int {
	D1 := math.Floor(math.Floor(math.Floor(float64(Potency)*F_AP(M, MainStat, LvMain)*F_DET(DET, LvMain, LvDiv))/100) / 1000)
	D2 := math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(D1*F_TNC(TNC, LvMain, LvSub, LvDiv))/1000)*F_WD(WD, LvMain, JobMod))/100)*float64(TraitMod)) / 100)
	return int(D2)
}
