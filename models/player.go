package models

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	ID                   int
	Name                 string
	Ping                 int           `json:"-"`
	Stats                *PlayerStats  `json:"-"`
	CalculatedCritRate   float64       `json:"calculatedCritRate"`
	CalculatedDirectRate float64       `json:"calculatedDirectRate"`
	ActualCritRate       float64       `json:"actualCritRate"`
	ActualDirectRate     float64       `json:"actualDirectRate"`
	Job                  Job           `json:"-"`
	CDUntil              int64         `json:"-"`
	AnimLockUntil        int64         `json:"-"`
	NextCombo            []*Skill      `json:"-"`
	Buffs                []AppliedBuff `json:"-"`
	SkillLog             []SkillLog    //`json:"-"`
}

type SkillLog struct {
	SkillID       int
	SkillName     string
	EncounterTime int64
	Damage        int
	Crit          bool
	Direct        bool
}

func (p *Player) Tick(encounterTime time.Duration, enemy *Enemy, wg *sync.WaitGroup, randGen *rand.Rand) {
	defer wg.Done()
	// Remove buffs that have expired
	for i, buff := range p.Buffs {
		if encounterTime.Milliseconds() > buff.AppliedUntil && buff.Buff.Duration > 0 {
			p.Buffs = append(p.Buffs[:i], p.Buffs[i+1:]...)
		}
	}
	// Update the charges
	p.Job.UpdateCharges(encounterTime.Milliseconds())
	// Check if the player is animation locked
	if encounterTime.Milliseconds() < p.AnimLockUntil {
		return
	}
	var desiredSkill *Skill
	// Check if the player is on GCD
	if encounterTime.Milliseconds() < p.CDUntil {
		// Check if we can use an OGCD
		desiredSkill = p.Job.SelectNextOGCD(p, enemy, encounterTime.Milliseconds(), p.CDUntil)
	} else {
		// Let the Job logic decide what skill to use
		desiredSkill = p.Job.SelectNextGCD(p, enemy, encounterTime.Milliseconds())
	}
	if desiredSkill != nil {
		desiredSkill.Execute(p, enemy, encounterTime, randGen)
		// Apply the animation lock
		p.AnimLockUntil = encounterTime.Milliseconds() + int64(desiredSkill.LockMS) + int64(p.Ping)
	}
}

func (p *Player) calculateDamage(potency int, randGen *rand.Rand, autocrit bool, autodirect bool) (int, bool, bool) {
	// TODO: Handle auto crit and auto direct modifiers

	// D1 = ⌊ Potency × f(ATK) × f(DET) ⌋ /100 ⌋ /1000 ⌋, multiplies the potency by the attack power and determination modifiers
	d1 := float64(potency) * p.Stats.AttackPowerMod * p.Stats.DeterminationMod / 100
	// D2 = ⌊ D1 × f(TNC) ⌋ /1000 ⌋ × f(WD) ⌋ /100 ⌋ × Trait ⌋ /100 ⌋, multiplies the result by the tenacity modifier and weapon damage modifier
	d2 := math.Floor(d1*p.Stats.TenacityMod) * p.Stats.WeaponDamageMod / 100 * p.Stats.TraitMod
	// D3 = ⌊ D2 × CRIT? ⌋ /1000 ⌋ × DH? ⌋ /100 ⌋, multiplies the result by the critical hit modifier and direct hit modifier
	// Roll for critical hit and direct hit
	crit := randGen.Float64() < p.Stats.CriticalHitPercent || autocrit
	direct := randGen.Float64() < p.Stats.DirectHitPercent || autodirect
	if crit {
		d2 *= p.Stats.CriticalHitMod
	}
	if direct {
		d2 *= 1.25
	}
	// D = ⌊ D3 × rand[95,105] ⌋ /100 ⌋, multiplies the result by a random number between 95 and 105 (damage variance)
	d := math.Floor(d2 * float64(randGen.Intn(10)+95) / 100)
	// Buffs
	for _, buff := range p.Buffs {
		if buff.Buff.DamageMod != 0 {
			d *= buff.Buff.DamageMod
		}
	}
	return int(d), crit, direct
}

func (p *Player) LogSkill(skill *Skill, damage int, encounterTime time.Duration, crit bool, direct bool) {
	// Create the SkillLog if it doesn't exist
	if p.SkillLog == nil {
		p.SkillLog = make([]SkillLog, 0)
	}
	p.SkillLog = append(p.SkillLog, SkillLog{
		SkillID:       skill.ID,
		SkillName:     skill.Name,
		EncounterTime: encounterTime.Milliseconds(),
		Damage:        damage,
		Crit:          crit,
		Direct:        direct,
	})
	//fmt.Println(p.Name, "used", skill.Name, "for", damage, "damage")
}
