package skill

import (
	"github.com/ATroschke/xivsim/sim/enemy"
	xivmath "github.com/ATroschke/xivsim/xiv-math"
)

// type and enum for GCD
type GCD int

const (
	AA = iota
	GCD1_5
	GCD2
	GCD2_5
	GCD2_8
	GCD3
	GCD3_5
	GCD4
	OGCD
)

type Skill struct {
	Name                       string
	ID                         int
	Potency                    int
	ComboPotency               int
	CalculatedDamage           int
	CalculatedAutoCDHDamage    int
	CalculatedPotDamage        int
	CalculatedPotAutoCDHDamage int
	GCD                        GCD
	BreaksCombo                bool
	LockMS                     int
	CooldownMS                 int64
	NextCharge                 int64
	MaxCharges                 int
	Charges                    int
	NextCombo                  []*Skill
	CustomLogic                func(job any, target *enemy.Enemy, time int64)
	AutoCDH                    bool
	DamageDealt                int
	Uses                       int
	MOverride                  int // If this in't 0, the Mainstat Modifier will be overwritten by it (used for Esteems Damage Calculation...)
}

func (s *Skill) CalculateDamage(
	weaponDamage int,
	mainStat int,
	criticalHit int,
	directHit int,
	determination int,
	skillSpeed int,
	spellSpeed int,
	tenacity int,
	weaponDelay float64,
	M int,
	JobMod int,
) {
	// Calculate the MainStat of the Character using a Tincture (Hardcoded to Grade 8 Tinctures... +10%/Max 262)
	potionMS := mainStat + 262
	ms10pct := int(mainStat / 10)
	if mainStat+ms10pct < potionMS {
		potionMS = mainStat + ms10pct
	}
	// Check if the Mainstat Modifier is ovrwritten
	modMS := M
	if s.MOverride != 0 {
		modMS = s.MOverride
	}
	if s.GCD != AA {
		// Calculate the base damage of the skill
		s.CalculatedDamage = xivmath.DirectDamage(
			s.Potency,
			weaponDamage,
			modMS,
			mainStat,
			390,
			JobMod,
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			400,
			100,
		)
		s.CalculatedAutoCDHDamage = xivmath.DirectDamageAutoCDH(
			s.Potency,
			weaponDamage,
			modMS,
			mainStat,
			390,
			JobMod, // 105 war, drk, 100 pld, gnb
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			directHit,
			400,
			100,
		)
		s.CalculatedPotDamage = xivmath.DirectDamage(
			s.Potency,
			weaponDamage,
			modMS,
			potionMS,
			390,
			JobMod,
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			400,
			100,
		)
		s.CalculatedPotAutoCDHDamage = xivmath.DirectDamageAutoCDH(
			s.Potency,
			weaponDamage,
			modMS,
			potionMS,
			390,
			JobMod, // 105 war, drk, 100 pld, gnb
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			directHit,
			400,
			100,
		)
	} else {
		s.CalculatedDamage = xivmath.AutoAttack(
			s.Potency,
			weaponDamage,
			weaponDelay,
			modMS,
			mainStat,
			390,
			JobMod,
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			400,
			100,
		)
		s.CalculatedPotDamage = xivmath.AutoAttack(
			s.Potency,
			weaponDamage,
			weaponDelay,
			modMS,
			potionMS,
			390,
			JobMod,
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			400,
			100,
		)
	}
}
