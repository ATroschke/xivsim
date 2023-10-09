package models

import (
	"math"
)

// TODO: Implement all stats and formulas
// https://www.akhmorning.com/allagan-studies/how-to-be-a-math-wizard/shadowbringers/functions/#lv-80-fap

// Currently assuming all players are at level 90
const (
	LvMAIN = 390
	LvSUB  = 400
	LvDIV  = 1900
	LvHP   = 3000
)

// Speed holds GCDs and Cast Times in ms
type Speed struct {
	AA_DOT_MOD float64
	GCD1_5     int
	GCD2       int
	GCD2_5     int
	GCD2_8     int
	GCD3       int
	GCD3_5     int
	GCD4       int
}

type PlayerStats struct {
	// Stats from Gear etc
	weaponDamage    int
	mainstat        int
	criticalHitRate int
	directHitRate   int
	determination   int
	tenacity        int
	piety           int
	speed           int
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
	Speed              Speed
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
		speed:           skillSpeed,
		vitality:        vitality,
		TraitMod:        1,
	}

	playerStats.calculateAttackPower()
	playerStats.calculateCriticalHit()
	playerStats.calculateDirectHit()
	playerStats.calculateDeterminiation()
	playerStats.calculateTenacity()
	playerStats.calculateWeaponDamage()
	playerStats.calculateSpeed()

	return &playerStats
}

func (ps *PlayerStats) calculateAttackPower() {
	// f(ATK) depends on the player's class
	// TODO: Implement JobMods, Currently assuming Warrior
	// Mdps = 195
	// Mtank = 156
	ps.AttackPowerMod = math.Floor(156 * (float64(ps.mainstat) - float64(LvMAIN)) / float64(LvMAIN))
}

func (ps *PlayerStats) calculateCriticalHit() {
	// Critical Hit Rate is calculated as follows: p(CRIT) = ⌊ 200 × ( CRIT - Level Lv, SUB)/ Level Lv, DIV  + 50 ⌋ / 10
	ps.CriticalHitPercent = float64(200*(ps.criticalHitRate-LvSUB)/LvDIV+50) / 1000.0
	// Critical Hit Strength is calculated as follows: f(CRIT) = 1400 + ⌊ 200 × ( CRIT - Level Lv, SUB)/ Level Lv, DIV ⌋
	ps.CriticalHitMod = float64(1400+(200*(ps.criticalHitRate-LvSUB)/LvDIV)) / 1000.0
}

func (ps *PlayerStats) calculateDirectHit() {
	// Direct Hit Rate is calculated as follows: p(DH) = ⌊ 550 × ( DH - Level Lv, SUB)/ Level Lv, DIV ⌋ / 10
	ps.DirectHitPercent = float64(550*(ps.directHitRate-LvSUB)/LvDIV) / 1000.0
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
	ps.WeaponDamageMod = math.Floor(((float64(LvMAIN) * 110.0) / 1000.0) + float64(ps.weaponDamage))
}

func (ps *PlayerStats) calculateSpeed() {
	// AA and DOT Modifier is calculated as follows: f(SPD) = ( 1000 + ⌊ 130 × ( Speed - Level Lv, SUB )/ Level Lv, DIV ⌋ ) / 1000
	ps.Speed.AA_DOT_MOD = (1000 + math.Floor(130*(float64(ps.speed)-float64(LvSUB))/float64(LvDIV))) / 1000.0
	// GCD is calculated as follows: =(INT(GCD*(1000+CEILING(130*(400-Speed)/1900))/10000)/100)
	ps.Speed.GCD1_5 = int(1500*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
	ps.Speed.GCD2 = int(2000*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
	ps.Speed.GCD2_5 = int(2500*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
	ps.Speed.GCD2_8 = int(2800*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
	ps.Speed.GCD3 = int(3000*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
	ps.Speed.GCD3_5 = int(3500*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
	ps.Speed.GCD4 = int(4000*(1000+math.Ceil(130*(400-float64(ps.speed))/1900))/10000) * 10
}
