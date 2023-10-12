package job

import (
	"math/rand"
	"reflect"

	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Warrior implements Job
type Warrior struct {
	GCDUntil           int64
	AnimationLockUntil int64
	Skills             WarriorSkills
	Speed              *Speed
	DamageModifiers    DamageModifiers
	NextCombo          []*skill.Skill
	//Buffs              []Buff
}

type WarriorSkills struct {
	HeavySwing skill.Skill
	Maim       skill.Skill
	StormsEye  skill.Skill
	StormsPath skill.Skill
}

func NewWarrior(speed *Speed) *Warrior {
	return &Warrior{
		Speed: speed,
		Skills: WarriorSkills{
			HeavySwing: HeavySwing,
			Maim:       Maim,
			StormsEye:  StormsEye,
			StormsPath: StormsPath,
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
	t := reflect.TypeOf(w.Skills)
	for i := 0; i < t.NumField(); i++ {
		skill := reflect.ValueOf(&w.Skills).Elem().Field(i).Addr().Interface().(*skill.Skill)
		skill.CalculateDamage(
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

func (w *Warrior) UseSkill(target *enemy.Enemy, usedSkill *skill.Skill, encounterTime int64) int {
	// We always need to set the AnimationLockUntil
	w.AnimationLockUntil = encounterTime + int64(usedSkill.LockMS)
	// If the skill is a GCD, set the GCDUntil
	if usedSkill.GCD != skill.OGCD {
		w.GCDUntil = encounterTime + int64(w.Speed.GetGCD(usedSkill.GCD))
	}

	// If the skill breaks the combo, reset the NextCombo
	if usedSkill.BreaksCombo {
		w.NextCombo = nil
	}

	// If the skill has a combo, set the NextCombo
	if usedSkill.NextCombo != nil {
		w.NextCombo = usedSkill.NextCombo
	}

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
	rolledDamage := float64(usedSkill.CalculatedDamage) * damageModifiers
	damage := target.TakeDamage(int(rolledDamage))
	return damage
}

func (w *Warrior) SelectNextSkill(encounterTime int64) *skill.Skill {
	// Check if the player is animation locked
	if encounterTime < w.AnimationLockUntil {
		// The player is animation locked, so we can't use any skill
		return nil
	}
	// Check if the GCD is ready
	if encounterTime >= w.GCDUntil {
		// GCD is ready, so we need to select a GCD
		return w.SelectNextGCD(encounterTime)
	}
	// Return nil if no skill can be used
	return nil
}

func (w *Warrior) SelectNextGCD(encounterTime int64) *skill.Skill {
	// Check if the player is in a combo
	if w.NextCombo != nil {
		// The player is in a combo, so we need to select the next skill in the combo
		if w.NextCombo[0].Name == StormsEye.Name {
			return &w.Skills.StormsEye
		}
		if w.NextCombo[0].Name == Maim.Name {
			return &w.Skills.Maim
		}
	}
	// The player is not in a combo, so we will select our first filler GCD
	return &w.Skills.HeavySwing
}

// Warrior Skills

// GCDs

var (
	HeavySwing = skill.Skill{
		Name:        "Heavy Swing",
		ID:          31,
		Potency:     200,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&Maim},
	}
	Maim = skill.Skill{
		Name:        "Maim",
		ID:          37,
		Potency:     300,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&StormsEye},
		CustomLogic: func(v any) {
			// Cast v to *Warrior
			w := v.(*Warrior)
		},
	}
	StormsEye = skill.Skill{
		Name:        "Storm's Eye",
		ID:          45,
		Potency:     440,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
	}
	StormsPath = skill.Skill{
		Name:        "Storm's Path",
		ID:          42,
		Potency:     440,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
	}
)
