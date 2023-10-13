package job

import "github.com/ATroschke/xivsim/sim/skill"

// Paladin implements Job
type Paladin struct {
	GCDUntil           int64
	AnimationLockUntil int64
	AutoAttack         skill.Skill
	Skills             PaladinSkills
	Speed              *Speed
	DamageModifiers    DamageModifiers
	NextCombo          []*skill.Skill
	Buffs              PaladinBuffs
	BeastGauge         int
	SkillLog           []SkillLog
}

type PaladinSkills struct {
	// GCDs
	// OGCDs
}

type PaladinBuffs struct {
}

func NewPaladin(speed *Speed) *Paladin {
	return &Paladin{
		AutoAttack: AutoAttack,
		Speed:      speed,
	}
}

// CopyPaladin

// GetRateAverages

// Report

// CalculateSkills

// Update

// Tick
