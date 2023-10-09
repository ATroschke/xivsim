package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Encounter is a struct that holds the players and enemiess, and manages the encounter loop.
type Encounter struct {
	id      int
	Players []*Player
	Enemies []*Enemy `json:"-"`
	// Arbitrary time value for how long the simulated encounter has been going (needed for GCDs, Cooldowns, etc.) in ms
	encounterTime time.Duration
	// Desired length of the encounter
	encounterLength time.Duration
	// How much time passes per tick
	tickLength time.Duration
	// How much DPS the encounter did
	DPS float64
	// The encounters seed (to reproduce the same encounter)
	Seed int64
}

// NewEncounter creates a new encounter with the given amount of players and enemiess.
func NewEncounter(id int, numPlayer int, numEnemies int, length time.Duration, tickLength time.Duration, ps *PlayerStats) *Encounter {
	players := make([]*Player, numPlayer)
	enemies := make([]*Enemy, numEnemies)

	for i := 0; i < numPlayer; i++ {
		players[i] = &Player{
			ID:    i,
			Name:  fmt.Sprintf("Player %d", i+1),
			Stats: ps,
			Job:   NewWarrior(),
			Ping:  10,
		}
	}

	for i := 0; i < numEnemies; i++ {
		enemies[i] = &Enemy{}
	}

	// Generate a random seed (since encounters are deterministic, this is needed to reproduce the same encounter, the encounters are created in batches of 10% of the total amount of encounters, so the seed is the current time + the encounter id)
	seed := time.Now().UnixNano() + int64(id)

	return &Encounter{
		id:              id,
		Players:         players,
		Enemies:         enemies,
		encounterLength: length,
		encounterTime:   0,
		tickLength:      tickLength,
		Seed:            seed,
	}
}

// Start starts the encounter loop.
func (e *Encounter) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	// Set the seed
	randGen := rand.New(rand.NewSource(e.Seed))
	for {
		var wg sync.WaitGroup
		wg.Add(len(e.Players))

		for _, player := range e.Players {
			go player.Tick(e.encounterTime, e.Enemies[0], &wg, randGen)
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
			e.DPS = float64(damageTaken) / (float64(e.encounterTime.Seconds()))
			// Round to 2 decimals
			e.DPS = float64(int(e.DPS*100)) / 100
			return
		default:
			// Continue
		}
	}
}

func (e *Encounter) CalculateCritDirectRate() {
	for _, player := range e.Players {
		player.CalculatedCritRate = player.Stats.CriticalHitPercent * 100
		player.CalculatedDirectRate = player.Stats.DirectHitPercent * 100
		player.ActualCritRate = 0
		player.ActualDirectRate = 0
		for _, log := range player.SkillLog {
			if log.Crit {
				player.ActualCritRate++
			}
			if log.Direct {
				player.ActualDirectRate++
			}
		}
		player.ActualCritRate /= float64(len(player.SkillLog))
		player.ActualCritRate = float64(int(player.ActualCritRate*100)) / 100
		player.ActualCritRate *= 100
		player.ActualDirectRate /= float64(len(player.SkillLog))
		player.ActualDirectRate = float64(int(player.ActualDirectRate*100)) / 100
		player.ActualDirectRate *= 100
	}
}
