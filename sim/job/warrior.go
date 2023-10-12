package job

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/ATroschke/xivsim/sim/buff"
	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Warrior implements Job
type Warrior struct {
	GCDUntil           int64
	AnimationLockUntil int64
	AutoAttack         skill.Skill
	Skills             WarriorSkills
	Speed              *Speed
	DamageModifiers    DamageModifiers
	NextCombo          []*skill.Skill
	Buffs              WarriorBuffs
	BeastGauge         int
}

// WarriorSkills contains Warrior specific skills
type WarriorSkills struct {
	// GCDs
	HeavySwing skill.Skill
	Maim       skill.Skill
	StormsEye  skill.Skill
	StormsPath skill.Skill
	FellCleave skill.Skill
	InnerChaos skill.Skill
	PrimalRend skill.Skill
	// OGCDs
	Upheaval     skill.Skill
	Onslaught    skill.Skill
	Infuriate    skill.Skill
	InnerRelease skill.Skill
}

// WarriorBuffs contains all Warrior specific buffs
type WarriorBuffs struct {
	SurgingTempest  buff.Buff
	NascentChaos    buff.Buff
	InnerRelease    buff.Buff
	PrimalRendReady buff.Buff
}

func NewWarrior(speed *Speed) *Warrior {
	return &Warrior{
		AutoAttack: AutoAttack,
		Speed:      speed,
		Skills: WarriorSkills{
			HeavySwing:   HeavySwing,
			Maim:         Maim,
			StormsEye:    StormsEye,
			StormsPath:   StormsPath,
			FellCleave:   FellClave,
			Upheaval:     Upheaval,
			Onslaught:    Onslaught,
			Infuriate:    Infuriate,
			InnerChaos:   InnerChaos,
			InnerRelease: InnerRelease,
			PrimalRend:   PrimalRend,
		},
		Buffs: WarriorBuffs{
			SurgingTempest:  SurgingTempest,
			NascentChaos:    NascentChaos,
			InnerRelease:    InnerReleaseBuff,
			PrimalRendReady: PrimalRendReady,
		},
	}
}

func CopyWarrior(w *Warrior) *Warrior {
	warrior := &Warrior{
		AutoAttack:         w.AutoAttack,
		Speed:              w.Speed,
		Skills:             w.Skills,
		Buffs:              w.Buffs,
		DamageModifiers:    w.DamageModifiers,
		GCDUntil:           0,
		AnimationLockUntil: 0,
	}
	return warrior
}

func (w *Warrior) Report() {
	var report string
	// Report total AA Damage, Average AA Damage, AA Uses
	report += "Auto Attack: "
	report += fmt.Sprintf("Total Damage: %d - ", w.AutoAttack.DamageDealt)
	report += fmt.Sprintf("Average Damage: %d - ", w.AutoAttack.DamageDealt/w.AutoAttack.Uses)
	report += fmt.Sprintf("Uses: %d\n", w.AutoAttack.Uses)
	// Report total Skill Damage, Average Skill Damage, Skill Uses
	report += "\nSkills:\n"
	t := reflect.TypeOf(w.Skills)
	for i := 0; i < t.NumField(); i++ {
		skill := reflect.ValueOf(&w.Skills).Elem().Field(i).Addr().Interface().(*skill.Skill)
		report += fmt.Sprintf("%s: - ", skill.Name)
		report += fmt.Sprintf("Total Damage: %d - ", skill.DamageDealt)
		report += fmt.Sprintf("Average Damage: %d - ", skill.DamageDealt/skill.Uses)
		report += fmt.Sprintf("Uses: %d\n", skill.Uses)
	}
	fmt.Print(report)
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
	weaponDelay := float64(w.Speed.AA) / 1000
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
			weaponDelay,
		)
	}
	w.AutoAttack.CalculateDamage(
		weaponDamage,
		mainStat,
		criticalHit,
		directHit,
		determination,
		skillSpeed,
		spellSpeed,
		tenacity,
		weaponDelay,
	)
}

func (w *Warrior) Update(encounterTime int64) {
	// Check if any buffs need to be removed
	bt := reflect.TypeOf(w.Buffs)
	for i := 0; i < bt.NumField(); i++ {
		buff := reflect.ValueOf(&w.Buffs).Elem().Field(i).Addr().Interface().(*buff.Buff)
		if buff.AppliedUntil <= encounterTime {
			buff.AppliedUntil = 0
		}
	}
	// Check if any skills need to be recharged
	st := reflect.TypeOf(w.Skills)
	for i := 0; i < st.NumField(); i++ {
		skill := reflect.ValueOf(&w.Skills).Elem().Field(i).Addr().Interface().(*skill.Skill)
		if skill.NextCharge > 0 {
			if skill.NextCharge <= encounterTime {
				skill.Charges++
				if skill.Charges >= skill.MaxCharges {
					skill.NextCharge = 0
				} else {
					skill.NextCharge += skill.CooldownMS
				}
			}
		}
	}
}

func (w *Warrior) IncreaseBeastGauge(amount int) {
	w.BeastGauge += amount
	if w.BeastGauge > 100 {
		w.BeastGauge = 100
	}
}

func (w *Warrior) DecreaseBeastGauge(amount int) {
	w.BeastGauge -= amount
	if w.BeastGauge < 0 {
		w.BeastGauge = 0
	}
}

// TODO: Tick is called every time the encounter time is increased
func (w *Warrior) Tick(enemy *enemy.Enemy, encounterTime int64) (int, int) {
	// Update the Warrior
	w.Update(encounterTime)
	// Auto Attack
	aa := 0
	s := 0
	if w.AutoAttack.NextCharge <= encounterTime || w.AutoAttack.NextCharge == 0 {
		aa = w.UseSkill(enemy, &w.AutoAttack, encounterTime)
	}
	// Select the next skill
	nextSkill := w.SelectNextSkill(encounterTime)
	if nextSkill != nil {
		// Use the skill
		s = w.UseSkill(enemy, nextSkill, encounterTime)
	}
	return aa, s
}

func (w *Warrior) UseSkill(target *enemy.Enemy, usedSkill *skill.Skill, encounterTime int64) int {
	usedSkill.Uses += 1
	// We always need to set the AnimationLockUntil
	if usedSkill.LockMS != 0 {
		w.AnimationLockUntil = encounterTime + int64(usedSkill.LockMS)
	}
	// If the skill is a GCD, set the GCDUntil
	if usedSkill.GCD != skill.OGCD && usedSkill.GCD != skill.AA {
		w.GCDUntil = encounterTime + int64(w.Speed.GetGCD(usedSkill.GCD))
	}
	// If the Skill is an AA, set its NextCharge
	if usedSkill.GCD == skill.AA {
		usedSkill.NextCharge = encounterTime + int64(w.Speed.AA)
	}

	// If the skill has charges, decrease the charges
	if usedSkill.MaxCharges > 0 {
		usedSkill.Charges--
		// If the Skill now has less than MaxCharges left, and NextCharge is 0, set NextCharge
		if usedSkill.Charges < usedSkill.MaxCharges && usedSkill.NextCharge == 0 {
			usedSkill.NextCharge = encounterTime + usedSkill.CooldownMS
		}
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
	// Roll for Crit and Direct Hit, if the skill is not AutoCDH
	autoCDH := usedSkill.AutoCDH
	var isCrit, isDirect bool
	// Special case for Fell Cleave, since it is AutoCDH if Inner Release is active
	if usedSkill.ID == FellClave.ID && w.Buffs.InnerRelease.Stacks > 0 {
		autoCDH = true
	} else if !autoCDH {
		randCrit := rand.Float64()
		randDirect := rand.Float64()
		// Check if the skill is a crit
		isCrit = randCrit <= w.DamageModifiers.CritRate
		// Check if the skill is a direct hit
		isDirect = randDirect <= w.DamageModifiers.DirectRate
	}
	// Calculate the damage modifiers
	damageModifiers := randDamage
	if isCrit || autoCDH {
		damageModifiers *= w.DamageModifiers.CritDamage
	}
	if isDirect || autoCDH {
		damageModifiers *= 1.25
	}
	// Apply Surgin Tempest damage modifier if the buff is active
	if w.Buffs.SurgingTempest.AppliedUntil > encounterTime {
		damageModifiers *= w.Buffs.SurgingTempest.DamageMod
	}
	// Calculate the damage
	rolledDamage := float64(usedSkill.CalculatedDamage) * damageModifiers
	damage := target.TakeDamage(int(rolledDamage))
	usedSkill.DamageDealt += damage

	// Apply custom logic if the skill has any
	if usedSkill.CustomLogic != nil {
		usedSkill.CustomLogic(w, encounterTime)
	}

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
	// Check if the default animation lock (700ms) would cut into the GCD
	if encounterTime+700 >= w.GCDUntil {
		// The default animation lock would cut into the GCD, so we can't use any skill
		return nil
	}
	// GCD is not ready and wouldn't be cut into, so we need to select an OGCD
	return w.SelectNextOGCD(encounterTime)
}

func (w *Warrior) SelectNextGCD(encounterTime int64) *skill.Skill {
	// TODO: Check if we are in a group, and if so, check if we need to use/hold
	if w.Buffs.NascentChaos.AppliedUntil-encounterTime > 7000 && w.Buffs.SurgingTempest.AppliedUntil-encounterTime > 7000 {
		return &w.Skills.InnerChaos
	} else if w.Buffs.NascentChaos.AppliedUntil > 0 && w.Buffs.NascentChaos.AppliedUntil-encounterTime < 3000 {
		return &w.Skills.InnerChaos
	}
	if w.Buffs.PrimalRendReady.Stacks != 0 {
		return &w.Skills.PrimalRend
	}
	// Check if the player has enough Beast Gauge for Fell Cleave, and if Surging Tempest is active
	if (w.Buffs.InnerRelease.Stacks > 0 || w.BeastGauge >= 50) && w.Buffs.SurgingTempest.AppliedUntil-encounterTime > 7000 {
		return &w.Skills.FellCleave
	}
	// Check if the player is in a combo
	if w.NextCombo != nil {
		// The player is in a combo, so we need to select the next skill in the combo
		if w.NextCombo[0].Name == StormsEye.Name || w.NextCombo[0].Name == StormsPath.Name {
			if w.Buffs.SurgingTempest.AppliedUntil-encounterTime > 7000 {
				return &w.Skills.StormsPath
			}
			return &w.Skills.StormsEye
		}
		if w.NextCombo[0].Name == Maim.Name {
			return &w.Skills.Maim
		}
	}
	// The player is not in a combo, so we will select our first filler GCD
	return &w.Skills.HeavySwing
}

func (w *Warrior) SelectNextOGCD(encounterTime int64) *skill.Skill {
	if w.Skills.Infuriate.Charges > 0 && w.BeastGauge <= 50 && w.Buffs.NascentChaos.AppliedUntil == 0 {
		return &w.Skills.Infuriate
	}
	if w.Buffs.SurgingTempest.AppliedUntil-encounterTime > 1500 {
		if w.Skills.InnerRelease.Charges > 0 {
			return &w.Skills.InnerRelease
		}
		if w.Skills.Upheaval.Charges > 0 {
			return &w.Skills.Upheaval
		}
		if w.Skills.Onslaught.Charges > 0 {
			return &w.Skills.Onslaught
		}
	}
	return nil
}

// Warrior Skills
var (
	// Auto Attack
	AutoAttack = skill.Skill{
		Name:        "Attack",
		ID:          1,
		Potency:     90,
		GCD:         skill.AA,
		Charges:     1,
		MaxCharges:  1,
		BreaksCombo: false,
		LockMS:      0,
	}
	// GCDs
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
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Increase the Beast Gauge by 10
			w.IncreaseBeastGauge(10)
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
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Increase the Beast Gauge by 10
			w.IncreaseBeastGauge(10)
			// Apply the Surging Tempest buff
			w.Buffs.SurgingTempest.AppliedUntil = time + w.Buffs.SurgingTempest.DurationMS
		},
	}
	StormsPath = skill.Skill{
		Name:        "Storm's Path",
		ID:          42,
		Potency:     440,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Increase the Beast Gauge by 20
			w.IncreaseBeastGauge(20)
		},
	}
	FellClave = skill.Skill{
		Name:        "Fell Cleave",
		ID:          3549,
		Potency:     520,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			if w.Buffs.InnerRelease.Stacks > 0 {
				w.Buffs.InnerRelease.Stacks--
			} else {
				// Decrease the Beast Gauge by 50
				w.DecreaseBeastGauge(50)
			}
			// Reduce the cooldown of Infuriate by 5 seconds
			w.Skills.Infuriate.NextCharge -= 5000
		},
	}
	InnerChaos = skill.Skill{
		Name:        "Inner Chaos",
		ID:          16462,
		Potency:     660,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		AutoCDH:     true,
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Remove the Nascent Chaos buff
			w.Buffs.NascentChaos.AppliedUntil = 0
			// Decrease the Beast Gauge by 50
			w.DecreaseBeastGauge(50)
			// Reduce the cooldown of Infuriate by 5 seconds
			w.Skills.Infuriate.NextCharge -= 5000
		},
	}
	PrimalRend = skill.Skill{
		Name:        "Primal Rend",
		ID:          25753,
		Potency:     700,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      1000,
		AutoCDH:     true,
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Remove the Primal Rend Ready buff
			w.Buffs.PrimalRendReady.Stacks--
		},
	}
	// OGCDs
	Upheaval = skill.Skill{
		Name:        "Upheaval",
		ID:          7387,
		Potency:     400,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  1,
		Charges:     1,
	}
	Onslaught = skill.Skill{
		Name:        "Onslaught",
		ID:          7386,
		Potency:     150,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  3,
		Charges:     3,
	}
	Infuriate = skill.Skill{
		Name:        "Infuriate",
		ID:          52,
		Potency:     0,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  2,
		Charges:     2,
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Apply the Nascent Chaos buff
			w.Buffs.NascentChaos.AppliedUntil = time + w.Buffs.NascentChaos.DurationMS
			// Increase the Beast Gauge by 50
			w.IncreaseBeastGauge(50)
		},
	}
	InnerRelease = skill.Skill{
		Name:        "Inner Release",
		ID:          7389,
		Potency:     0,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(v any, time int64) {
			// Cast v to *Warrior
			w := v.(*Warrior)
			// Apply the Inner Release buff
			w.Buffs.InnerRelease.Stacks = w.Buffs.InnerRelease.MaxStacks
			// Apply the Primal Rend Ready buff
			w.Buffs.PrimalRendReady.Stacks = w.Buffs.PrimalRendReady.MaxStacks
			// Extend the duration of the Surging Tempest buff by 10 seconds
			w.Buffs.SurgingTempest.AppliedUntil += 10000
		},
	}
	// Buffs
	SurgingTempest = buff.Buff{
		Name:         "Surging Tempest",
		ID:           2677,
		DamageMod:    1.1,
		DurationMS:   30000,
		AppliedUntil: 0,
	}
	NascentChaos = buff.Buff{
		Name:       "Nascent Chaos",
		ID:         1897,
		DurationMS: 30000,
	}
	InnerReleaseBuff = buff.Buff{
		Name:      "Inner Release",
		ID:        1177,
		Stacks:    0,
		MaxStacks: 3,
	}
	PrimalRendReady = buff.Buff{
		Name:      "Primal Rend Ready",
		ID:        2624,
		Stacks:    0,
		MaxStacks: 1,
	}
)
