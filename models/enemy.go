package models

import "sync"

type Enemy struct {
	damageTaken int
	mutex       sync.Mutex
}

func (e *Enemy) takeDamage(damage int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.damageTaken += damage
}

func (e *Enemy) getDamageTaken() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.damageTaken
}
