package models

import (
	"sync"
)

type Enemy struct {
	damageTaken int
	mutex       sync.Mutex
}

func (e *Enemy) takeDamage(damage int, crit bool, direct bool) (int, bool, bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.damageTaken += damage
	return damage, crit, direct
}

func (e *Enemy) getDamageTaken() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.damageTaken
}
