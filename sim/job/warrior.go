package job

import (
	"github.com/ATroschke/xivsim/sim/buff"
	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Warrior implements Job
type Warrior struct {
	Skills     WarriorSkills
	Buffs      WarriorBuffs
	BeastGauge int
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
	Tincture        buff.Buff
}

func NewWarrior() *Warrior {
	return &Warrior{
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
			Tincture:        TinctureBuff,
		},
	}
}

func CopyWarrior(j *Job) *Warrior {
	w := j.JobImpl.(*Warrior)
	return &Warrior{
		Skills: WarriorSkills{
			HeavySwing:   w.Skills.HeavySwing,
			Maim:         w.Skills.Maim,
			StormsEye:    w.Skills.StormsEye,
			StormsPath:   w.Skills.StormsPath,
			FellCleave:   w.Skills.FellCleave,
			Upheaval:     w.Skills.Upheaval,
			Onslaught:    w.Skills.Onslaught,
			Infuriate:    w.Skills.Infuriate,
			InnerChaos:   w.Skills.InnerChaos,
			InnerRelease: w.Skills.InnerRelease,
			PrimalRend:   w.Skills.PrimalRend,
		},
		Buffs: WarriorBuffs{
			SurgingTempest:  w.Buffs.SurgingTempest,
			NascentChaos:    w.Buffs.NascentChaos,
			InnerRelease:    w.Buffs.InnerRelease,
			PrimalRendReady: w.Buffs.PrimalRendReady,
			Tincture:        w.Buffs.Tincture,
		},
	}
}

func (w *Warrior) GetSkills() []*skill.Skill {
	return []*skill.Skill{
		&w.Skills.HeavySwing,
		&w.Skills.Maim,
		&w.Skills.StormsEye,
		&w.Skills.StormsPath,
		&w.Skills.FellCleave,
		&w.Skills.Upheaval,
		&w.Skills.Onslaught,
		&w.Skills.Infuriate,
		&w.Skills.InnerChaos,
		&w.Skills.InnerRelease,
		&w.Skills.PrimalRend,
	}
}

func (w *Warrior) GetBuffs() []*buff.Buff {
	return []*buff.Buff{
		&w.Buffs.SurgingTempest,
		&w.Buffs.NascentChaos,
		&w.Buffs.InnerRelease,
		&w.Buffs.PrimalRendReady,
		&w.Buffs.Tincture,
	}
}

func (w *Warrior) GetDots() []*buff.DOT {
	return []*buff.DOT{}
}

func (w *Warrior) GetJobMod() (int, int) {
	return 156, 105
}

func (w *Warrior) NextSkill(job *Job, encounterTime int64) *skill.Skill {
	// Normal Rotation
	if encounterTime >= 0 {
		// Check if the player is animation locked
		if encounterTime < job.AnimationLockUntil {
			// The player is animation locked, so we can't use any skill
			return nil
		}
		// Check if the GCD is ready
		if encounterTime >= job.GCDUntil {
			// GCD is ready, so we need to select a GCD
			return w.SelectNextGCD(job, encounterTime)
		}
		// Check if the default animation lock (700ms) would cut into the GCD
		if encounterTime+700 >= job.GCDUntil {
			// The default animation lock would cut into the GCD, so we can't use any skill
			return nil
		}
		// GCD is not ready and wouldn't be cut into, so we need to select an OGCD
		return w.SelectNextOGCD(encounterTime)
	} else {
		// Warrior has no Prepull
		return nil
	}
}

func (w *Warrior) SelectNextGCD(job *Job, encounterTime int64) *skill.Skill {
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
	if job.NextCombo != nil {
		// The player is in a combo, so we need to select the next skill in the combo
		if job.NextCombo[0].Name == StormsEye.Name || job.NextCombo[0].Name == StormsPath.Name {
			if w.Buffs.SurgingTempest.AppliedUntil-encounterTime > 7000 {
				return &w.Skills.StormsPath
			}
			return &w.Skills.StormsEye
		}
		if job.NextCombo[0].Name == Maim.Name {
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

func (w *Warrior) ValidateAutoCDH(s *skill.Skill) bool {
	if s.AutoCDH {
		return true
	}
	if s.Name == w.Skills.FellCleave.Name && w.Buffs.InnerRelease.Stacks > 0 {
		return true
	}
	return false
}

func (w *Warrior) GetBuffModifiers(encounterTime int64) float64 {
	if w.Buffs.SurgingTempest.AppliedUntil > encounterTime {
		return w.Buffs.SurgingTempest.DamageMod
	}
	return 1
}

func (w *Warrior) CheckPotionStatus(time int64) bool {
	//return w.Buffs.Tincture.AppliedUntil >= time
	return false
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

func (w *Warrior) GetDesiredPrepullDuration() int {
	return 0
}

// Warrior Skills
var (
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
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
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			// Cast v to *Warrior
			w := job.(*Job).JobImpl.(*Warrior)
			// Apply the Inner Release buff
			w.Buffs.InnerRelease.Stacks = w.Buffs.InnerRelease.MaxStacks
			// Apply the Primal Rend Ready buff
			w.Buffs.PrimalRendReady.Stacks = w.Buffs.PrimalRendReady.MaxStacks
			// Extend the duration of the Surging Tempest buff by 10 seconds
			// Does IR apply Surging Tempest if appliedUntil < time+10000?
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
