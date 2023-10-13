package encounter

import (
	"fmt"
	"sync"
	"time"

	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/player"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Encounter is a struct that contains all the information about an encounter and
// and handles the simulation of the encounter.

const (
	STATUS_NOT_STARTED = iota
	STATUS_RUNNING
	STATUS_FINISHED
	STATUS_ERROR
)

type Encounter struct {
	ID             string           // TODO: Will later be a GUID of the Encounter
	BasePlayers    []*player.Player // Holds the Base Players of the Encounter, which are then copied for each iteration of the Simulation
	Iterations     []*Iteration     // Holds all the Iterations of the Encounter
	NrIterations   int              // Number of Iterations of the Encounter
	Duration       int              // Duration of the Encounter in Seconds
	Ping           int              // Simulated Ping in Milliseconds
	Tick           int              // Simulated Tick in Milliseconds (How often the Simulation is updated)
	Status         int              // Status of the Encounter
	SimulationTime int              // Time the Simulation took in Milliseconds
	Downtime       []Downtime       // Downtime of the Encounter
}

type Downtime struct {
	Start int
	End   int
}

type Iteration struct {
	Seed             int64 // Seed for the RNG, saved for reproducibility
	encounterTime    int   // Current Encounter Time
	maxEncounterTime int   // Max Encounter Time
	Enemy            *enemy.Enemy
	Players          []*player.Player
	TotalDamage      int // Total Damage taken by the Enemies, used to sort the Iterations
	Downtime         []Downtime
}

type PlayerDamage struct {
	Player        *player.Player
	AverageDamage int
	LowestDamage  *FightResult
	HighestDamage *FightResult
}

type FightResult struct {
	Player   *player.Player
	Damage   int
	CritRate float64
	DHitRate float64
	CDHRate  float64
}

// NewEncounter creates a new Encounter with the given ID, Base Enemies and Base Players.
func NewEncounter(id string, basePlayers []*player.Player, duration int, ping int, tick int, iterations int, downtimes *[]Downtime) *Encounter {
	return &Encounter{
		ID:           id,
		BasePlayers:  basePlayers,
		Iterations:   make([]*Iteration, iterations),
		NrIterations: iterations,
		Duration:     duration,
		Ping:         ping,
		Tick:         tick,
		Status:       STATUS_NOT_STARTED,
		Downtime:     *downtimes,
	}
}

// Run starts the simulation of the Encounter.
func (e *Encounter) Run() {
	simulationStart := time.Now()
	e.Status = STATUS_RUNNING
	// Create all the Iterations
	for i := 0; i < e.NrIterations; i++ {
		// Create a new Iteration
		e.Iterations[i] = NewIteration(int64(i), e.Duration, &e.BasePlayers, &e.Downtime)
	}
	// Create a waitgroup for the Iterations
	var wg sync.WaitGroup
	wg.Add(len(e.Iterations))
	// Run all the Iterations
	for i := 0; i < len(e.Iterations); i++ {
		go e.Iterations[i].Run(&wg, e.Tick)
	}
	// Wait for all the Iterations to finish
	wg.Wait()
	// TODO: Sort the Iterations by Total Damage
	simulationEnd := time.Now()
	e.Status = STATUS_FINISHED
	e.SimulationTime = int(simulationEnd.Sub(simulationStart).Milliseconds())
}

// Report prints a report of the Encounter.
func (e *Encounter) Report() {
	fmt.Printf("Encounter: %s\nSimulation Time: %dms\nIterations: %d - Ping: %dms - Duration: %ds - Tick: %dms\n\n", e.ID, e.SimulationTime, len(e.Iterations), e.Ping, e.Duration, e.Tick)
	TotalDamage := 0
	LowestDamage := 0
	HighestDamage := 0
	for i := 0; i < len(e.Iterations); i++ {
		TotalDamage += e.Iterations[i].TotalDamage
		if e.Iterations[i].TotalDamage < LowestDamage || LowestDamage == 0 {
			LowestDamage = e.Iterations[i].TotalDamage
		}
		if e.Iterations[i].TotalDamage > HighestDamage {
			HighestDamage = e.Iterations[i].TotalDamage
		}
	}
	// Calculate the Average Encounter Damage
	AverageDamage := TotalDamage / len(e.Iterations)
	AverageDPS := AverageDamage / e.Duration
	LowestDPS := LowestDamage / e.Duration
	HighestDPS := HighestDamage / e.Duration
	// Calculate the Damage and DPS per Player
	var playerDamage []PlayerDamage
	// Build the PlayerDamage slice (based on the Base Players)
	for i := 0; i < len(e.BasePlayers); i++ {
		playerDamage = append(playerDamage, PlayerDamage{
			Player:        e.BasePlayers[i],
			AverageDamage: 0,
		})
	}
	// Add the Damage of each Player to the PlayerDamage slice
	for i := 0; i < len(e.Iterations); i++ {
		for j := 0; j < len(e.Iterations[i].Players); j++ {
			for k := 0; k < len(playerDamage); k++ {
				if playerDamage[k].Player.Name == e.Iterations[i].Players[j].Name {
					playerDamage[k].AverageDamage += e.Iterations[i].Players[j].DamageDealt
					if playerDamage[k].LowestDamage == nil || playerDamage[k].LowestDamage.Damage > e.Iterations[i].Players[j].DamageDealt {
						playerDamage[k].LowestDamage = &FightResult{
							Damage: e.Iterations[i].Players[j].DamageDealt,
							Player: e.Iterations[i].Players[j],
						}
					}
					if playerDamage[k].HighestDamage == nil || playerDamage[k].HighestDamage.Damage < e.Iterations[i].Players[j].DamageDealt {
						playerDamage[k].HighestDamage = &FightResult{
							Damage: e.Iterations[i].Players[j].DamageDealt,
							Player: e.Iterations[i].Players[j],
						}
					}
				}
			}
		}
	}
	// Calculate the Average Damage per Player
	for i := 0; i < len(playerDamage); i++ {
		playerDamage[i].AverageDamage /= len(e.Iterations)
	}
	// Calculate the Rate Averages for the Lowest an Highest Damage Runs
	for i := 0; i < len(playerDamage); i++ {
		playerDamage[i].LowestDamage.CritRate, playerDamage[i].LowestDamage.DHitRate, playerDamage[i].LowestDamage.CDHRate = playerDamage[i].LowestDamage.Player.Job.GetRateAverages()
		playerDamage[i].HighestDamage.CritRate, playerDamage[i].HighestDamage.DHitRate, playerDamage[i].HighestDamage.CDHRate = playerDamage[i].HighestDamage.Player.Job.GetRateAverages()
	}
	fmt.Printf("Total Damage (All Encounters): %d\nAverage Damage: %d (%d DPS)\nLowest Damage: %d (%d DPS)\nHighest Damage: %d (%d DPS)\n\n",
		TotalDamage, AverageDamage, AverageDPS, LowestDamage, LowestDPS, HighestDamage, HighestDPS)
	for i := 0; i < len(playerDamage); i++ {
		fmt.Printf("Player: %s - Average Damage: %d (%d DPS) - Lowest Damage: %d (%d DPS, c:%f%%,d:%f%%,cdh:%f%%) - Highest Damage: %d (%d DPS, c:%f%%,d:%f%%,cd:%f%%)\n",
			playerDamage[i].Player.Name, playerDamage[i].AverageDamage, playerDamage[i].AverageDamage/e.Duration,
			playerDamage[i].LowestDamage.Damage, playerDamage[i].LowestDamage.Damage/e.Duration, playerDamage[i].LowestDamage.CritRate, playerDamage[i].LowestDamage.DHitRate, playerDamage[i].LowestDamage.CDHRate,
			playerDamage[i].HighestDamage.Damage, playerDamage[i].HighestDamage.Damage/e.Duration, playerDamage[i].HighestDamage.CritRate, playerDamage[i].HighestDamage.DHitRate, playerDamage[i].HighestDamage.CDHRate)
		playerDamage[i].HighestDamage.Player.Job.Report()
	}
}

func PrintSkillDamage(s skill.Skill, iterations int, duration int, totalDamage int) {
	fmt.Printf("%s: %d Average Uses - %d Damage/Iteration - %d Damage/Use - %d DPS - %f%%\n",
		s.Name,
		s.Uses/iterations,
		s.DamageDealt/iterations,
		s.DamageDealt/s.Uses,
		s.DamageDealt/(iterations*duration),
		float64(s.DamageDealt)/float64(totalDamage)*100)
}

func NewIteration(seed int64, duration int, players *[]*player.Player, downtimes *[]Downtime) *Iteration {
	enemy := &enemy.Enemy{
		Name: fmt.Sprintf("Test Enemy ITERATION %d", seed),
	}
	// Create a new Iteration
	iteration := &Iteration{
		Seed:             seed,
		Enemy:            enemy,
		Players:          nil,
		maxEncounterTime: duration * 1000,
		Downtime:         *downtimes,
	}
	// Copy the Base Players
	for j := 0; j < len(*players); j++ {
		iteration.Players = append(iteration.Players, player.CopyPlayer((*players)[j]))
	}
	return iteration
}

func (i *Iteration) Run(encounterWG *sync.WaitGroup, tick int) {
	defer encounterWG.Done()
	// Loop until the encounter is over (duration is reached)
	for i.encounterTime < i.maxEncounterTime {
		// Check there are downtimes specified
		if len(i.Downtime) == 0 {
			// No downtimes specified, so we just run the current Tick
			i.Tick()
			// Increase the encounter time by the duration of the Tick
			i.encounterTime += tick
			continue
		}
		for j := 0; j < len(i.Downtime); j++ {
			if i.encounterTime >= i.Downtime[j].Start && i.encounterTime < i.Downtime[j].End {
				// We are in a Downtime, so we skip this Tick
				i.encounterTime += tick
				continue
			}
			// Run the current Tick
			i.Tick()
			// Increase the encounter time by the duration of the Tick
			i.encounterTime += tick
		}
	}
	// Calculate the total Damage taken by the Enemies
	i.TotalDamage = i.Enemy.DamageTaken
	// Calculate the total Damage dealt by the Players and compare it to the total Damage taken by the Enemies
	sumPlayerDamage := 0
	for j := 0; j < len(i.Players); j++ {
		sumPlayerDamage += i.Players[j].DamageDealt
	}
	if sumPlayerDamage == i.TotalDamage {
		// The total Damage dealt by the Players is equal to the total Damage taken by the Enemies, so the simulation was successful
	} else {
		// The total Damage dealt by the Players is not equal to the total Damage taken by the Enemies, so the simulation was not successful
		fmt.Printf("Iteration %d encountered an error: %d != %d\n", i.Seed, sumPlayerDamage, i.TotalDamage)
	}
}

func (i *Iteration) Tick() {
	var wg sync.WaitGroup
	wg.Add(len(i.Players))
	// Run the current Tick for all the Players
	for j := 0; j < len(i.Players); j++ {
		go i.Players[j].Tick(&wg, i.encounterTime, i.Enemy, i.Players)
	}
	// Wait for all the Players to finish
	wg.Wait()
}
