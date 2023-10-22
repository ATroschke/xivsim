package enemy

import (
	"math"
	"math/rand"
	"sync"

	"github.com/ATroschke/xivsim/sim/buff"
	xivmath "github.com/ATroschke/xivsim/xiv-math"
)

type Enemy struct {
	Name        string
	DamageTaken int
	NextDOTTick int64
	Dots        []*buff.DOT
	mutex       sync.Mutex // Mutex for Locking the enemy when taking damage
}

func NewEnemy(Name string, r *rand.Rand) *Enemy {
	// Generate a random first DOTTick between 0 and 3000ms into the fight
	firstTick := int64(math.Floor(r.Float64() * 3000))
	return &Enemy{
		Name:        Name,
		NextDOTTick: firstTick,
	}
}

func (e *Enemy) Tick(encounterTime int64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	// Tick all Dots
	if e.NextDOTTick <= encounterTime {
		for _, dot := range e.Dots {
			// If the Dot has expired, remove it from the list
			if dot.AppliedUntil < encounterTime {
				e.Dots = e.Dots[1:]
				continue
			}

			// Get a random damage modifier
			randMod := xivmath.RandMod()
			// Roll for Crit and Direct Hit
			crit := xivmath.RandBool(dot.CritRate)
			directHit := xivmath.RandBool(dot.DirectRate)
			// Calculate the Damage
			var damageDealt int
			if crit && directHit {
				damageDealt = int(float64(dot.CritDirectDamage) * randMod)
			} else if crit {
				damageDealt = int(float64(dot.CritDamage) * randMod)
			} else if directHit {
				damageDealt = int(float64(dot.DirectDamage) * randMod)
			} else {
				damageDealt = int(float64(dot.Damage) * randMod)
			}
			// Apply the Buff Modifier
			damageDealt = int(float64(damageDealt) * dot.BuffMod)
			e.DamageTaken += damageDealt
			// Attribute the Damage to the Source
			*dot.Source += damageDealt
		}
		// Set the next tick time
		e.NextDOTTick = encounterTime + 3000
	}
}

func (e *Enemy) TakeDamage(damage int) int {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	// TODO: Calculate additional Damage from Debuffs and add it to the total Damage (for DPS)
	// TODO: Calculate additional Damage Mod from Debuffs an their sources so that rDPS can be correctly calculated
	// Add Damage taken to the total Damage taken, which acts as a control value for the simulation
	// Total Player Damage (DPS or rDPS) shouldn't be higher than this value
	e.DamageTaken += damage
	// Return the Damage taken
	return damage
}

func (e *Enemy) ApplyDot(dot *buff.DOT, source *int, time int64, buffmod float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	// Add the DOT to the list of Dots
	dot.AppliedUntil = time + dot.DurationMS
	dot.Source = source
	dot.BuffMod = buffmod
	e.Dots = append(e.Dots, dot)
}
