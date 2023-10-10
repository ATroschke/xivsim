package job

import (
	"math/rand"

	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Warrior implements Job
type Warrior struct {
	GCDUntil           int64
	AnimationLockUntil int64
	Skills             []skill.Skill
	Speed              *Speed
	DamageModifiers    DamageModifiers
	//Buffs              []Buff
}

func NewWarrior(speed *Speed) *Warrior {
	return &Warrior{
		Speed: speed,
		Skills: []skill.Skill{
			HeavySwing,
		},
	}
}

func CopyWarrior(w *Warrior) *Warrior {
	warrior := &Warrior{
		Speed:              w.Speed,
		Skills:             w.Skills,
		DamageModifiers:    w.DamageModifiers,
		GCDUntil:           0,
		AnimationLockUntil: 0,
	}
	return warrior
}

// CalculateSkills calculates the base damage of all skills (without crit, direct hit, buffs, etc.)
func (w *Warrior) CalculateSkills(
	weaponDamage int,
	mainStat int,
	criticalHit int,
	directHit int,
	determination int,
	skillSpeed int,
	spellSpeed int,
	tenacity int,
) {
	// Calculate variable damage modifiers
	w.DamageModifiers = CalculateDamageModifiers(criticalHit, directHit, 400, 1900)
	// TODO: Pass down all needed values from the Players stats
	for i := range w.Skills {
		w.Skills[i].CalculateDamage(
			weaponDamage,
			mainStat,
			criticalHit,
			directHit,
			determination,
			skillSpeed,
			spellSpeed,
			tenacity,
		)
	}
}

// TODO: AddBuff adds a buff to the player
func (w *Warrior) AddBuff() {

}

// TODO: Tick is called every time the encounter time is increased
func (w *Warrior) Tick(enemy *enemy.Enemy, encounterTime int64) int {
	// Select the next skill
	nextSkill := w.SelectNextSkill(encounterTime)
	if nextSkill != nil {
		// Use the skill
		return w.UseSkill(enemy, nextSkill, encounterTime)
	} else {
		// The GCD is not ready, so we can't use any skill
		return 0
	}
}

func (w *Warrior) UseSkill(target *enemy.Enemy, skill *skill.Skill, encounterTime int64) int {
	// Set the GCD until the next skill can be used
	w.GCDUntil = encounterTime + int64(w.Speed.GetGCD(skill.GCD))
	// Roll for Damage Range (0.95-1.05)
	randDamage := rand.Float64()*0.1 + 0.95
	// Roll for Crit and Direct Hit
	randCrit := rand.Float64()
	randDirect := rand.Float64()
	// Check if the skill is a crit
	isCrit := randCrit <= w.DamageModifiers.CritRate
	// Check if the skill is a direct hit
	isDirect := randDirect <= w.DamageModifiers.DirectRate
	// Calculate the damage modifiers
	damageModifiers := randDamage
	if isCrit {
		damageModifiers *= w.DamageModifiers.CritDamage
	}
	if isDirect {
		damageModifiers *= 1.25
	}
	// Calculate the damage
	rolledDamage := float64(skill.CalculatedDamage) * damageModifiers

	damage := target.TakeDamage(int(rolledDamage))
	return damage
}

func (w *Warrior) SelectNextSkill(encounterTime int64) *skill.Skill {
	// Check if the GCD is ready
	if encounterTime >= w.GCDUntil {
		// The GCD is ready, so we can use the next skill
		return &w.Skills[0]
	}
	// The GCD is not ready, so we can't use any skill
	return nil
}

// Warrior Skills
var (
	HeavySwing = skill.Skill{
		Name:        "Heavy Swing",
		ID:          31,
		Potency:     200,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
	}
)
