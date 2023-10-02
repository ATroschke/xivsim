package models

import (
	"math"
	"math/rand"
	"sync"
)

type Player struct {
	Name  string
	Stats *PlayerStats
}

func (p *Player) Tick(enemy *Enemy, wg *sync.WaitGroup) {
	defer wg.Done()
	// Random damage
	damage := p.calculateDamage(100)
	enemy.takeDamage(damage)
}

func (p *Player) calculateDamage(potency int) int {
	// D1 = ⌊ Potency × f(ATK) × f(DET) ⌋ /100 ⌋ /1000 ⌋, multiplies the potency by the attack power and determination modifiers
	d1 := float64(potency) * p.Stats.AttackPowerMod * p.Stats.DeterminationMod / 100
	// D2 = ⌊ D1 × f(TNC) ⌋ /1000 ⌋ × f(WD) ⌋ /100 ⌋ × Trait ⌋ /100 ⌋, multiplies the result by the tenacity modifier and weapon damage modifier
	d2 := math.Floor(d1*p.Stats.TenacityMod) * p.Stats.WeaponDamageMod / 100 * p.Stats.TraitMod
	// D3 = ⌊ D2 × CRIT? ⌋ /1000 ⌋ × DH? ⌋ /100 ⌋, multiplies the result by the critical hit modifier and direct hit modifier
	// Roll for critical hit and direct hit
	crit := rand.Float64() < p.Stats.CriticalHitPercent
	direct := rand.Float64() < p.Stats.DirectHitPercent
	if crit {
		d2 *= p.Stats.CriticalHitMod
	}
	if direct {
		d2 *= 1.25
	}
	// D = ⌊ D3 × rand[95,105] ⌋ /100 ⌋, multiplies the result by a random number between 95 and 105 (damage variance)
	d := math.Floor(d2 * float64(rand.Intn(10)+95) / 100)
	// TODO: Buffs

	return int(d)
}