package skill

import (
	"fmt"

	xivmath "github.com/ATroschke/xivsim/xiv-math"
)

// type and enum for GCD
type GCD int

const (
	GCD1_5 = iota
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
	CooldownMS       int
	MaxCharges       int
	//AppliesBuffs     []Buff
	NextCombo   []*Skill
	CustomLogic func(v any)
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
) {
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
	fmt.Printf("%s: %d\n", s.Name, s.CalculatedDamage)
}
