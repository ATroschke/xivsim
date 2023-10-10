package enemy

import (
	"sync"
)

type Enemy struct {
	Name        string
	DamageTaken int
	// Debuffs
	// Mutex for Locking the enemy when taking damage
	mutex sync.Mutex
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
