package models

import (
	"fmt"
	"sync"
	"time"
)

// Encounter is a struct that holds the players and enemiess, and manages the encounter loop.
type Encounter struct {
	id      int
	Players []*Player
	Enemies []*Enemy
	// Arbitrary time value for how long the simulated encounter has been going (needed for GCDs, Cooldowns, etc.) in ms
	encounterTime time.Duration
	// Desired length of the encounter
	encounterLength time.Duration
	// How much time passes per tick
	tickLength time.Duration
	// How much DPS the encounter did
	dps float64
}

// NewEncounter creates a new encounter with the given amount of players and enemiess.
func NewEncounter(id int, numPlayer int, numEnemies int, length time.Duration, tickLength time.Duration, ps *PlayerStats) *Encounter {
	players := make([]*Player, numPlayer)
	enemies := make([]*Enemy, numEnemies)

	for i := 0; i < numPlayer; i++ {
		players[i] = &Player{
			Name:  fmt.Sprintf("Player %d", i+1),
			Stats: ps,
		}
	}

	for i := 0; i < numEnemies; i++ {
		enemies[i] = &Enemy{}
	}

	return &Encounter{
		id:              id,
		Players:         players,
		Enemies:         enemies,
		encounterLength: length,
		encounterTime:   0,
		tickLength:      tickLength,
	}
}

// GetDPS returns the DPS of the encounter.
func (e *Encounter) GetDPS() float64 {
	return e.dps
}

// Start starts the encounter loop.
func (e *Encounter) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var wg sync.WaitGroup
		wg.Add(len(e.Players))

		for _, player := range e.Players {
			go player.Tick(e.Enemies[0], &wg)
		}

		// Waits for all players to finish their tick
		wg.Wait()
		// Increment encounter time
		e.encounterTime += e.tickLength

		switch {
		case e.encounterTime >= e.encounterLength:
			// Get how much damage the enemy took
			damageTaken := e.Enemies[0].getDamageTaken()
			// Calculate DPS
			e.dps = float64(damageTaken) / (float64(e.encounterTime.Seconds()))
			return
		default:
			// Continue
		}
	}
}
