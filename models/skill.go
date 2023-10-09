package models

import (
	"math/rand"
	"time"
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
	Name         string
	ID           int
	Potency      int
	ComboPotency int
	GCD          GCD
	BreaksCombo  bool
	LockMS       int
	CooldownMS   int
	MaxCharges   int
	AppliesBuffs []Buff
	NextCombo    []*Skill
	CustomLogic  func(player *Player, job Job) (bool, bool)
}

type Cooldown struct {
	SkillID       int
	CooldownUntil int64
}

func (s *Skill) Execute(player *Player, enemy *Enemy, encounterTime time.Duration, randGen *rand.Rand) {
	// Set the CDUntil to the current time + the GCD of the skill according to the player's speed
	var playerGCD int
	switch s.GCD {
	case GCD1_5:
		playerGCD = player.Stats.Speed.GCD1_5
	case GCD2:
		playerGCD = player.Stats.Speed.GCD2
	case GCD2_5:
		playerGCD = player.Stats.Speed.GCD2_5
	case GCD2_8:
		playerGCD = player.Stats.Speed.GCD2_8
	case GCD3:
		playerGCD = player.Stats.Speed.GCD3
	case GCD3_5:
		playerGCD = player.Stats.Speed.GCD3_5
	case GCD4:
		playerGCD = player.Stats.Speed.GCD4
	}
	// If the skill is an OGCD, don't apply the GCD
	if s.GCD != OGCD {
		player.CDUntil = encounterTime.Milliseconds() + int64(playerGCD)
	}
	// Apply custom logic (if the skill has any), this should also return if the skill should auto-crit or auto-direct
	autocrit, autodirect := false, false
	if s.CustomLogic != nil {
		autocrit, autodirect = s.CustomLogic(player, player.Job)
	}
	// Apply damage if the skill has a potency
	skillDamage := 0
	crit := false
	direct := false
	if s.Potency > 0 {
		skillDamage, crit, direct = enemy.takeDamage(player.calculateDamage(s.Potency, randGen, autocrit, autodirect))
	}
	// Apply buffs
	for _, buff := range s.AppliesBuffs {
		// Check if the buff is already applied
		var buffAlreadyApplied bool
		for i, appliedBuff := range player.Buffs {
			if appliedBuff.Buff.ID == buff.ID {
				// If the buff is already applied, check if it should be extended
				if buff.ApplyType == ApplyTypeSelfExtend {
					// If the buff should be extended, set the new duration
					// new duration = time until the buff expires + the buff duration
					player.Buffs[i].AppliedUntil = player.Buffs[i].AppliedUntil + buff.Duration
					buffAlreadyApplied = true
					break
				}
				// If the buff is already applied and it shouldn't be extended, reapply the buff
				player.Buffs[i].AppliedUntil = encounterTime.Milliseconds() + buff.Duration
				buffAlreadyApplied = true
				break
			}
		}
		// If the buff is not already applied, apply it
		if !buffAlreadyApplied {
			player.Buffs = append(player.Buffs, AppliedBuff{
				Buff:         buff,
				SourceID:     player.ID,
				AppliedUntil: encounterTime.Milliseconds() + buff.Duration,
				Stacks:       buff.Stacks,
			})
		}
	}
	// Set the next combo (if the skill has one)
	if s.NextCombo != nil {
		player.NextCombo = s.NextCombo
	} else {
		// Remove the combo if the combo is finished
		// HACK/TODO: Removes the Combo if the Skill Breaks Combo, but has no NextCombo. This is a hack to make the simulation work for now.
		if s.BreaksCombo {
			player.NextCombo = nil
		}
	}
	// Log the skill
	player.LogSkill(s, skillDamage, encounterTime, crit, direct)
}
