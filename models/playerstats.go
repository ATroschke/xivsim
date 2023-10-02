package models

import "math"

// TODO: Implement all stats and formulas
// https://www.akhmorning.com/allagan-studies/how-to-be-a-math-wizard/shadowbringers/functions/#lv-80-fap

// Currently assuming all players are at level 90
const (
	LvMAIN = 390
	LvSUB  = 400
	LvDIV  = 1900
	LvHP   = 3000
)

type PlayerStats struct {
	// Stats from Gear etc
	weaponDamage    int
	mainstat        int
	criticalHitRate int
	directHitRate   int
	determination   int
	tenacity        int
	piety           int
	skillSpeed      int
	spellSpeed      int
	vitality        int
	// Calculated stats
	AttackPowerMod     float64
	WeaponDamageMod    float64
	CriticalHitPercent float64
	CriticalHitMod     float64
	DirectHitPercent   float64
	DeterminationMod   float64
	TenacityMod        float64
	TraitMod           float64 // Depends on class, defaulting to 1 (for now)
}

func NewPlayerStats(weaponDamage int, mainstat int, criticalHitRate int, directHitRate int, determination int, tenacity int, piety int, skillSpeed int, spellSpeed int, vitality int) *PlayerStats {
	playerStats := PlayerStats{
		weaponDamage:    weaponDamage,
		mainstat:        mainstat,
		criticalHitRate: criticalHitRate,
		directHitRate:   directHitRate,
		determination:   determination,
		tenacity:        tenacity,
		piety:           piety,
		skillSpeed:      skillSpeed,
		spellSpeed:      spellSpeed,
		vitality:        vitality,
		TraitMod:        1,
	}

	playerStats.calculateAttackPower()
	playerStats.calculateCriticalHit()
	playerStats.calculateDeterminiation()
	playerStats.calculateTenacity()
	playerStats.calculateWeaponDamage()

	return &playerStats
}

func (ps *PlayerStats) calculateAttackPower() {
	// f(ATK) depends on the player's class
	ps.AttackPowerMod = math.Floor(165 * (float64(ps.mainstat) - float64(LvMAIN)) / float64(LvMAIN))
}

func (ps *PlayerStats) calculateCriticalHit() {
	// Critical Hit Rate is calculated as follows: p(CRIT) = ⌊ 200 × ( CRIT - Level Lv, SUB)/ Level Lv, DIV  + 50 ⌋ / 10
	ps.CriticalHitPercent = float64(200*(ps.criticalHitRate-LvSUB)/LvDIV+50) / 1000.0
	// Critical Hit Strength is calculated as follows: f(CRIT) = 1400 + ⌊ 200 × ( CRIT - Level Lv, SUB)/ Level Lv, DIV ⌋
	ps.CriticalHitMod = float64(1400+(200*(ps.criticalHitRate-LvSUB)/LvDIV)) / 1000.0
}

func (ps *PlayerStats) calculateDeterminiation() {
	// Determination is calculated as follows: f(DET) = ⌊ 140 × ( DET - Level Lv, MAIN )/ Level Lv, DIV + 1000 ⌋
	ps.DeterminationMod = float64(140*(ps.determination-LvMAIN)/LvDIV+1000) / 1000.0
}

func (ps *PlayerStats) calculateTenacity() {
	// Tenacity is calculated as follows: f(TNC) = 1000 + ⌊ 100 × ( TNC - Level Lv, SUB )/ Level Lv, DIV ⌋
	ps.TenacityMod = float64(1000+(100*(ps.tenacity-LvSUB)/LvDIV)) / 1000.0
}

func (ps *PlayerStats) calculateWeaponDamage() {
	// Weapon damage is calculated as follows: f(WD) = ⌊ ( LevelModLv, MAIN · JobModJob, Attribute / 1000 ) + WD ⌋
	// TODO: Implement JobMods, Currently assuming Warrior
	ps.WeaponDamageMod = math.Floor((float64(LvMAIN)/105.0)/1000) + float64(ps.weaponDamage)
}
