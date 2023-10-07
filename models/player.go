package models

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	ID        int
	Name      string
	Stats     *PlayerStats
	Job       Job
	CDUntil   int64
	NextCombo []*Skill
	Buffs     []AppliedBuff
}

func (p *Player) Tick(encounterTime time.Duration, enemy *Enemy, wg *sync.WaitGroup, randGen *rand.Rand) {
	defer wg.Done()
	// Remove buffs that have expired
	for i, buff := range p.Buffs {
		if encounterTime.Milliseconds() > buff.AppliedUntil {
			p.Buffs = append(p.Buffs[:i], p.Buffs[i+1:]...)
		}
	}
	// Check if the player is on GCD
	if encounterTime.Milliseconds() < p.CDUntil {
		return
	}
	// Let the Job logic decide what skill to use
	desiredSkill := p.Job.SelectNextSkill(p, enemy, encounterTime.Milliseconds())
	// Log the skill
	fmt.Printf("%s:  \t %s (%s) used %s\n", encounterTime.String(), p.Name, p.Job.GetName(), desiredSkill.Name)
	// Execute the skill (Apply buffs, etc.)
	desiredSkill.Execute(p, enemy, encounterTime, randGen)
}

func (p *Player) calculateDamage(potency int, randGen *rand.Rand) int {
	// D1 = ⌊ Potency × f(ATK) × f(DET) ⌋ /100 ⌋ /1000 ⌋, multiplies the potency by the attack power and determination modifiers
	d1 := float64(potency) * p.Stats.AttackPowerMod * p.Stats.DeterminationMod / 100
	// D2 = ⌊ D1 × f(TNC) ⌋ /1000 ⌋ × f(WD) ⌋ /100 ⌋ × Trait ⌋ /100 ⌋, multiplies the result by the tenacity modifier and weapon damage modifier
	d2 := math.Floor(d1*p.Stats.TenacityMod) * p.Stats.WeaponDamageMod / 100 * p.Stats.TraitMod
	// D3 = ⌊ D2 × CRIT? ⌋ /1000 ⌋ × DH? ⌋ /100 ⌋, multiplies the result by the critical hit modifier and direct hit modifier
	// Roll for critical hit and direct hit
	crit := randGen.Float64() < p.Stats.CriticalHitPercent
	direct := randGen.Float64() < p.Stats.DirectHitPercent
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
	return int(d)
}
