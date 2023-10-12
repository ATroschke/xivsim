package skill

import (
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
	Name             string
	ID               int
	Potency          int
	ComboPotency     int
	CalculatedDamage int
	GCD              GCD
	BreaksCombo      bool
	LockMS           int
	CooldownMS       int64
	NextCharge       int64
	MaxCharges       int
	Charges          int
	NextCombo        []*Skill
	CustomLogic      func(v any, time int64)
	AutoCDH          bool
	DamageDealt      int
	Uses             int
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
) {
	if s.GCD != AA {
		// Calculate the base damage of the skill
		s.CalculatedDamage = xivmath.DirectDamage(
			s.Potency,
			weaponDamage,
			156,
			mainStat,
			390,
			110,
			determination,
			1900,
			tenacity,
			skillSpeed,
			criticalHit,
			400,
			100,
		)
	} else {
		s.CalculatedDamage = xivmath.AutoAttack(
			s.Potency,
			weaponDamage,
			weaponDelay,
			156,
			mainStat,
			390,
			110,
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
