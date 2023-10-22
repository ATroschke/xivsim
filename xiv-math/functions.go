package xivmath

import (
	"math"
	"math/rand"
)

// TODO:
/* Discord (Mahdi):
Something I will note.

There are 4 different formulae for damage in the game. Each are similar to each other with some differences.

Casters and physical jobs use a different formula from each other.

Casters have the potency multiplied into WD, and then inted as the first step.

Where physical will have potency multiplied into f(ap) and then trunced to 2 decimals.

And then they are similar after that where f(ap) and WD get multiplied into each other and INTed.

DoTs, aside from scaling with sps or sks, will have crit and dh happen after random().

AAs scale with SkS and have an altered WD used which is wd x delay / 3 and then INTed.

And lastly channel spells are different then dots, and currently only use the physical formulae even if casted by a caster. Otherwise they are the same as other Dots.

And then there is the +1 added for sub 100 potency.

And we/I need to do testing into buffs and debuffs ordering. But it appears that negative dmg downs will roundup to next int instead of rounding down.
*/

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

func F_AUTOCDH(DET int, DH int, LvMain int, LvSub int, LvDiv int) float64 {
	return math.Floor(140*(float64(DET-LvMain))/float64(LvDiv)+1000) + math.Floor(140*(float64(DH-LvSub))/float64(LvDiv)+1)
}

func F_TNC(TNC int, LvMain int, LvSub int, LvDiv int) float64 {
	return math.Floor(100*(float64(TNC-LvSub))/float64(LvDiv) + 1000)
}

func F_SPD(SPD int, LvSub int, LvDiv int) float64 {
	// From Excel: =INT(130*(O12-400)/1900+1000)
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

// AutoCDH Skills gain additional damage in the DET function.
func DirectDamageAutoCDH(Potency int, WD int, M int, MainStat int, LvMain int, JobMod int, DET int, LvDiv int, TNC int, SPD int, CRIT int, DH int, LvSub int, TraitMod int) int {
	D1 := math.Floor(math.Floor(math.Floor(float64(Potency)*F_AP(M, MainStat, LvMain)*F_AUTOCDH(DET, DH, LvMain, LvSub, LvDiv))/100) / 1000)
	D2 := math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(D1*F_TNC(TNC, LvMain, LvSub, LvDiv))/1000)*F_WD(WD, LvMain, JobMod))/100)*float64(TraitMod)) / 100)
	return int(D2)
}

// Calculates the damage of a Physical DoT Tick with the given Potency and Stats.
// This will be used to pre-bake the damage of a skill for a player without Crit, Direct Hit and Buffs to speed up the calculation.
func PhysicalDOTDamage(Potency int, WD int, M int, MainStat int, LvMain int, JobMod int, DET int, LvDiv int, TNC int, SPD int, CRIT int, LvSub int, TraitMod int) int {
	D1 := math.Floor(math.Floor(math.Floor(float64(Potency)*F_AP(M, MainStat, LvMain)*F_DET(DET, LvMain, LvDiv))/100) / 1000)
	D2 := math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(D1*F_TNC(TNC, LvMain, LvSub, LvDiv))/1000)*F_SPD(SPD, LvSub, LvDiv))/1000)*F_WD(WD, LvMain, JobMod))/100)*float64(TraitMod)/100) + 1
	return int(D2)
}

// Calculates the damage of a Magical DoT Tick with the given Potency and Stats.
// This will be used to pre-bake the damage of a skill for a player without Crit, Direct Hit and Buffs to speed up the calculation.
func MagicalDOTDamage(Potency int, WD int, M int, MainStat int, LvMain int, JobMod int, DET int, LvDiv int, TNC int, SPD int, CRIT int, LvSub int, TraitMod int) int {
	D1 := math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(float64(Potency)*F_WD(WD, LvMain, JobMod))/100)*F_AP(M, MainStat, LvMain))/100)*F_SPD(SPD, LvSub, LvDiv)) / 1000)
	D2 := math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(D1*F_DET(DET, LvMain, LvDiv))/1000)*F_TNC(TNC, LvMain, LvSub, LvDiv))/1000)*float64(TraitMod))/100) + 1
	return int(D2)
}

// Calculates the damage of an attack with the given Potency and Stats.
// This will be used to pre-bake the damage of a skill for a player without Crit, Direct Hit and Buffs to speed up the calculation.
func AutoAttack(Potency int, WD int, WDelay float64, M int, MainStat int, LvMain int, JobMod int, DET int, LvDiv int, TNC int, SPD int, CRIT int, LvSub int, TraitMod int) int {
	D1 := math.Floor(math.Floor(math.Floor(float64(Potency)*F_AP(M, MainStat, LvMain)*F_DET(DET, LvMain, LvDiv))/100) / 1000)
	D2 := math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(math.Floor(D1*F_TNC(TNC, LvMain, LvSub, LvDiv))/1000)*F_SPD(SPD, LvSub, LvDiv))/1000)*F_AUTO(LvMain, JobMod, WD, WDelay))/100)*float64(TraitMod)) / 100)
	return int(D2)
}

func RandMod() float64 {
	return 0.95 + (rand.Float64() * 0.1)
}

func RandBool(chance float64) bool {
	return rand.Float64() <= chance
}
