package xivmath

import "math"

// Contains functions for calculating specific values for the game.
// Most formulas are taken from Allagan Studies

func F_WD(WD int, LvMain int, JobMod int) float64 {
	return math.Floor((float64(LvMain) * float64(JobMod) / 1000.0) + float64(WD))
}

func F_AP(M int, MainStat int, LvMain int) float64 {
	return math.Floor(float64(M)*(float64(MainStat)-float64(LvMain))/float64(LvMain)) + 100
}

func F_DET(DET int, LvMain int, LvDiv int) float64 {
	return math.Floor(130*(float64(DET-LvMain))/float64(LvDiv) + 1000)
}

func F_TNC(TNC int, LvMain int, LvSub int, LvDiv int) float64 {
	return math.Floor(100*(float64(TNC-LvSub))/float64(LvDiv) + 1000)
}

func F_SPD(SPD int, LvSub int, LvDiv int) float64 {
	return math.Floor(130*(float64(SPD-LvSub))/float64(LvDiv) + 1000)
}

func F_CRIT(CRIT int, LvSub int, LvDiv int) float64 {
	return math.Floor(200*(float64(CRIT-LvSub))/float64(LvDiv) + 1400)
}

func F_AUTO(LvMain int, JobMod int, WD int, WDelay float64) float64 {
	return math.Floor((math.Floor((float64(LvMain) * float64(JobMod) / 1000) + float64(WD))) * (WDelay / 3.0))
}
