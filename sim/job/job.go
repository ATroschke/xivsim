package job

import (
	"fmt"
	"math"

	"github.com/ATroschke/xivsim/sim/buff"
	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
	xivmath "github.com/ATroschke/xivsim/xiv-math"
)

type IJob interface {
	GetSkills() []*skill.Skill
	GetBuffs() []*buff.Buff
	NextSkill(j *Job, encounterTime int64) *skill.Skill
	ValidateAutoCDH(s *skill.Skill) bool
	GetBuffModifiers(encounterTime int64) float64
}

type Job struct {
	Name               string
	GCDUntil           int64
	AnimationLockUntil int64
	AutoAttack         skill.Skill
	Speed              *Speed
	DamageModifiers    DamageModifiers
	NextCombo          []*skill.Skill
	SkillLog           []SkillLog
	JobImpl            IJob
}

type DamageModifiers struct {
	CritRate   float64
	CritDamage float64
	DirectRate float64
}

type SkillLog struct {
	Name      string
	Damage    int
	Crit      bool
	DirectHit bool
	Time      int64
}

// NewJob creates a new Job with the given name and Speed.
func NewJob(name string, speed *Speed) *Job {
	var jobImpl IJob
	switch name {
	case "WAR":
		jobImpl = NewWarrior(speed)
	}
	return &Job{
		Name:       name,
		AutoAttack: AutoAttack,
		Speed:      speed,
		NextCombo:  nil,
		JobImpl:    jobImpl,
	}
}

// Create a new instance of a job with the same stats as the given job.
func (j *Job) CopyJob() *Job {
	var jobImpl IJob
	switch j.Name {
	case "WAR":
		jobImpl = CopyWarrior(j)
	}
	return &Job{
		Name:               j.Name,
		GCDUntil:           j.GCDUntil,
		AnimationLockUntil: j.AnimationLockUntil,
		AutoAttack:         j.AutoAttack,
		Speed:              j.Speed,
		DamageModifiers:    j.DamageModifiers,
		NextCombo:          j.NextCombo,
		SkillLog:           j.SkillLog,
		JobImpl:            jobImpl,
	}
}

func CalculateDamageModifiers(criticalHit int, directHit int, LvSub int, LvDiv int) DamageModifiers {
	return DamageModifiers{
		CritRate:   xivmath.P_CHR(criticalHit, LvSub, LvDiv),
		CritDamage: xivmath.F_CRIT(criticalHit, LvSub, LvDiv),
		DirectRate: xivmath.P_DHR(directHit, LvSub, LvDiv),
	}
}

func (j *Job) Tick(enemy *enemy.Enemy, encounterTime int64) (int, int) {
	// Update the Player
	j.Update(encounterTime)
	// Keep track of the damage dealt this tick
	var aaDamage, skillDamage int
	// Check if the Player can Auto Attack
	if j.AutoAttack.NextCharge <= encounterTime || j.AutoAttack.NextCharge == 0 {
		aaDamage = j.UseSkill(enemy, &j.AutoAttack, encounterTime)
	}
	// Select the next skill to use
	skill := j.JobImpl.NextSkill(j, encounterTime)
	if skill != nil {
		// Use the skill
		skillDamage = j.UseSkill(enemy, skill, encounterTime)
	}
	return aaDamage, skillDamage
}

func (j *Job) Update(encounterTime int64) {
	// Get the JobImpl's Buffs
	buffs := j.JobImpl.GetBuffs()
	// Check if any Buff needs to be removed
	for i := range buffs {
		if buffs[i].AppliedUntil <= encounterTime {
			buffs[i].AppliedUntil = 0
		}
	}
	// Check if any skills need to be recharged
	skills := j.JobImpl.GetSkills()
	for i := range skills {
		if skills[i].NextCharge <= encounterTime && skills[i].NextCharge != 0 {
			skills[i].Charges++
			if skills[i].Charges >= skills[i].MaxCharges {
				skills[i].NextCharge = 0
			} else {
				skills[i].NextCharge += skills[i].CooldownMS
			}
		}
	}
}

func (j *Job) UseSkill(enemy *enemy.Enemy, s *skill.Skill, encounterTime int64) int {
	s.Uses++
	if s.GCD == skill.OGCD {
		fmt.Printf("%s: %d\n", s.Name, encounterTime)
	}
	// We always need to set the Animation Lock
	j.AnimationLockUntil = encounterTime + int64(s.LockMS)
	// If the skill is a GCD, set GCDUntil
	if s.GCD != skill.AA && s.GCD != skill.OGCD {
		j.GCDUntil = encounterTime + int64(j.Speed.GetGCD(s.GCD))
	}
	// If the Skill is an AA, set its NextCharge
	if s.GCD == skill.AA {
		s.NextCharge = encounterTime + int64(j.Speed.AA)
	} else {
		// If the skill has charges, reduce the charges
		if s.MaxCharges > 0 {
			s.Charges--
			// If the Skill now has less than MaxCharges, set its NextCharge
			if s.Charges < s.MaxCharges {
				s.NextCharge = encounterTime + s.CooldownMS
			}
		}
	}

	// If the skill breaks the combo, reset the combo
	if s.BreaksCombo {
		j.NextCombo = nil
	}

	// If the skill has a combo, set the next combo
	if s.NextCombo != nil {
		j.NextCombo = s.NextCombo
	}

	// If the skill deals damage, calculate the damage
	damage := 0
	var crit, directHit, autoCDH bool
	if s.Potency > 0 {
		// Instantiate a Random Number Generator
		randMod := xivmath.RandMod()
		autoCDH = j.JobImpl.ValidateAutoCDH(s)
		if !autoCDH {
			crit = xivmath.RandBool(j.DamageModifiers.CritRate)
			directHit = xivmath.RandBool(j.DamageModifiers.DirectRate)
		}
		if crit || autoCDH {
			randMod *= j.DamageModifiers.CritDamage
		}
		if directHit || autoCDH {
			randMod *= 1.25
		}
		buffMod := j.JobImpl.GetBuffModifiers(encounterTime)
		if autoCDH {
			damage = int(math.Floor(float64(s.CalculatedAutoCDHDamage) * randMod * buffMod))
		} else {
			damage = int(math.Floor(float64(s.CalculatedDamage) * randMod * buffMod))
		}
		damage = enemy.TakeDamage(damage)
		s.DamageDealt += damage
	}

	// Apply the skill's custom logic
	if s.CustomLogic != nil {
		s.CustomLogic(j, encounterTime)
	}

	j.SkillLog = append(j.SkillLog, SkillLog{
		Name:      s.Name,
		Damage:    damage,
		Crit:      crit || autoCDH,
		DirectHit: directHit || autoCDH,
		Time:      encounterTime,
	})

	return damage
}

func (j *Job) GetRateAverages() (float64, float64, float64) {
	var critRate, directRate, critDirectRate float64
	for _, s := range j.SkillLog {
		if s.Crit {
			critRate += 1
		}
		if s.DirectHit {
			directRate += 1
		}
		if s.Crit && s.DirectHit {
			critDirectRate += 1
		}
	}
	return critRate / float64(len(j.SkillLog)) * 100, directRate / float64(len(j.SkillLog)) * 100, critDirectRate / float64(len(j.SkillLog)) * 100
}

func (j *Job) Report() {
	// Print the Skill Log
	fmt.Println("Skill Log:")
	for _, s := range j.SkillLog {
		// Time: Skill - Damage
		fmt.Printf("%d: %s - %d", s.Time, s.Name, s.Damage)
		// If it's a crit + directhit, append a !!
		if s.Crit && s.DirectHit {
			fmt.Printf("!!")
		} else if s.Crit {
			fmt.Printf("!")
		} else if s.DirectHit {
			fmt.Printf("*")
		}
		// TODO: Write job specific stuff
		// Append a newline
		fmt.Printf("\n")
	}
}

// CalculateSkills calculates the base damage of all skills (without crit, direct hit, buffs, etc.)
func (j *Job) CalculateSkills(
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
	j.DamageModifiers = CalculateDamageModifiers(criticalHit, directHit, 400, 1900)
	weaponDelay := float64(j.Speed.AA) / 1000
	// Get the skills to calculate from the JobImpl
	skills := j.JobImpl.GetSkills()
	// Iterate over all skills
	for i := range skills {
		// Calculate the base damage of the skill
		skills[i].CalculateDamage(
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
	j.AutoAttack.CalculateDamage(
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

// Default Skills
var // Auto Attack
AutoAttack = skill.Skill{
	Name:        "Attack",
	ID:          7,
	Potency:     90, // TODO: Should be 80 for Casters
	GCD:         skill.AA,
	Charges:     1,
	MaxCharges:  1,
	BreaksCombo: false,
	LockMS:      0,
}
