package encounter

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/job"
	"github.com/ATroschke/xivsim/sim/player"
	"github.com/ATroschke/xivsim/sim/skill"
	"github.com/ATroschke/xivsim/sim/statistics"
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
	Prepull        int              // Prepull duration in MS
	Downtime       []Downtime       // Downtime of the Encounter
	results        EncounterResults
}

type Downtime struct {
	Start int
	End   int
}

type Iteration struct {
	Seed             int64      // Seed for the RNG, saved for reproducibility
	Rand             *rand.Rand // The Rand Instance of the iteration
	encounterTime    int        // Current Encounter Time
	maxEncounterTime int        // Max Encounter Time
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
	// Players should exist by now, so iterate through them and check what Prepull Duration we need
	prepull := 0
	for i := 0; i < len(basePlayers); i++ {
		if basePlayers[i].GetPrepullDuration() > prepull {
			prepull = basePlayers[i].GetPrepullDuration()
		}
	}
	return &Encounter{
		ID:           id,
		BasePlayers:  basePlayers,
		NrIterations: iterations,
		Duration:     duration,
		Ping:         ping,
		Tick:         tick,
		Status:       STATUS_NOT_STARTED,
		Downtime:     *downtimes,
		Prepull:      prepull,
	}
}

// Run starts the simulation of the Encounter.
func (e *Encounter) Run() {
	simulationStart := time.Now()
	e.Status = STATUS_RUNNING
	// Decide how many Iterations to run in parallel
	batchSize := 1000
	iterations := e.NrIterations
	if iterations < batchSize {
		batchSize = iterations
	}
	// Prepare the Results
	var results EncounterResults
	// Initialize XivSimStatistics to calculate the Confidence Interval
	xst := &statistics.XivSimStatistics{}
	// Add the Base Players to the Results
	for i := 0; i < len(e.BasePlayers); i++ {
		results.PlayerResults = append(results.PlayerResults, PlayerResults{
			Player: e.BasePlayers[i].Name,
		})
	}
	// Run the Iterations in batches
	for i := 0; i < iterations; i += batchSize {
		// Create the Iterations for the current batch
		for j := 0; j < batchSize; j++ {
			// Create a new Iteration
			iteration := NewIteration(int64(i+j), e.Duration, e.Prepull, &e.BasePlayers, &e.Downtime)
			// Add the Iteration to the Encounter
			e.Iterations = append(e.Iterations, iteration)
		}
		// Create a waitgroup for the batch
		var batchWG sync.WaitGroup
		batchWG.Add(batchSize)
		// Run the current batch
		for j := 0; j < batchSize; j++ {
			go e.Iterations[i+j].Run(&batchWG, e.Tick)
		}
		// Wait for the current batch to finish
		batchWG.Wait()
		// Loop through the Iterations of the current batch and analyze the results
		for j := 0; j < batchSize; j++ {
			itReferenced := false
			results.TotalDamage += e.Iterations[i+j].TotalDamage
			// Add to the XivSimStatistics
			xst.AddDamage(e.Iterations[i+j].TotalDamage)
			// Check if the current Iteration has the lowest Damage or the highest Damage
			if results.LowestDamage == 0 || results.LowestDamage > e.Iterations[i+j].TotalDamage {
				results.LowestDamage = e.Iterations[i+j].TotalDamage
			}
			if results.HighestDamage < e.Iterations[i+j].TotalDamage {
				results.HighestDamage = e.Iterations[i+j].TotalDamage
			}
			// Loop through the Players of the current Iteration
			for k := 0; k < len(e.Iterations[i+j].Players); k++ {
				// Add the Damage dealt by the Player to the total Damage dealt by the Player
				results.PlayerResults[k].AverageDamage += e.Iterations[i+j].Players[k].Job.DamageDealt
				pReferenced := false
				// Check if the current Player has the lowest Damage
				if results.PlayerResults[k].LowestDamage == 0 || results.PlayerResults[k].LowestDamage > e.Iterations[i+j].Players[k].Job.DamageDealt {
					results.PlayerResults[k].LowestDamage = e.Iterations[i+j].Players[k].Job.DamageDealt
					results.PlayerResults[k].LowestDamageLog = e.Iterations[i+j].Players[k].Job.SkillLog
					pReferenced = true
					itReferenced = true
				}
				// Check if the current Player has the highest Damage
				if results.PlayerResults[k].HighestDamage < e.Iterations[i+j].Players[k].Job.DamageDealt {
					results.PlayerResults[k].HighestDamage = e.Iterations[i+j].Players[k].Job.DamageDealt
					results.PlayerResults[k].HighestDamageLog = e.Iterations[i+j].Players[k].Job.SkillLog
					pReferenced = true
					itReferenced = true
				}
				// Free up memory if the Player is not referenced anymore
				if !pReferenced {
					e.Iterations[i+j].Players[k] = nil
				}
			}
			// Free up memory if the Iteration is not referenced anymore
			if !itReferenced {
				e.Iterations[i+j] = nil
			}
		}
		// Calculate the Margin of Error
		me := xst.MarginOfError(1.96)
		results.MarginOfError = me
		results.ConfidencePct = me / xst.Mean() * 100
		// Print the current progress every batchSize * 10 Iterations
		if i%batchSize*10 == 0 {
			// clear the console
			fmt.Print("\033[H\033[2J")
			fmt.Printf("Progress: %d/%d\n", i, e.NrIterations)
			fmt.Printf("Margin of Error: %.2f -  %.4f%%\n", results.MarginOfError, results.ConfidencePct)
		}
		// Trigger the GC if the OS memory usage is above 90%
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		if mem.Sys > 0 && mem.Alloc > 0 && mem.Sys-mem.Alloc > mem.Sys/10 {
			runtime.GC()
		}
	}
	// Calculate the Confidence Interval, Save the current number of Iterations done
	results.NrIterations = e.NrIterations
	results.ConfidenceLow, results.ConfidenceHigh = xst.ConfidenceInterval(1.96)
	// Save the Results
	e.results = results
	simulationEnd := time.Now()
	e.Status = STATUS_FINISHED
	e.SimulationTime = int(simulationEnd.Sub(simulationStart).Milliseconds())
}

type EncounterResults struct {
	Name           string
	SimulationTime int
	NrIterations   int
	Duration       int
	Ping           int
	Tick           int
	TotalDamage    int
	AverageDamage  int
	AverageDPS     int
	LowestDamage   int
	HighestDamage  int
	LowestDPS      int
	HighestDPS     int
	PlayerResults  []PlayerResults
	MarginOfError  float64
	ConfidenceLow  float64
	ConfidenceHigh float64
	ConfidencePct  float64
}

type PlayerResults struct {
	Player           string
	AverageDamage    int
	LowestDamage     int
	HighestDamage    int
	AverageDPS       int
	LowestDPS        int
	HighestDPS       int
	LowestDamageLog  []job.SkillLog
	HighestDamageLog []job.SkillLog
}

// Report prints a report of the Encounter.
func (e *Encounter) Report() EncounterResults {
	e.results.SimulationTime = e.SimulationTime
	e.results.Duration = e.Duration
	// Calculate the average Damage
	e.results.AverageDamage = e.results.TotalDamage / e.NrIterations
	// Calculate the average DPS
	e.results.AverageDPS = e.results.TotalDamage / (e.NrIterations * e.Duration)
	// Calculate the lowest DPS
	e.results.LowestDPS = e.results.LowestDamage / e.Duration
	// Calculate the highest DPS
	e.results.HighestDPS = e.results.HighestDamage / e.Duration
	// Calculate the average Damage for each Player
	for i := 0; i < len(e.results.PlayerResults); i++ {
		e.results.PlayerResults[i].AverageDamage = e.results.PlayerResults[i].AverageDamage / e.results.NrIterations
		// Calculate the average DPS for each Player
		e.results.PlayerResults[i].AverageDPS = e.results.PlayerResults[i].AverageDamage / e.Duration
		// Calculate the lowest DPS for each Player
		e.results.PlayerResults[i].LowestDPS = e.results.PlayerResults[i].LowestDamage / e.Duration
		// Calculate the highest DPS for each Player
		e.results.PlayerResults[i].HighestDPS = e.results.PlayerResults[i].HighestDamage / e.Duration
	}

	return e.results
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

func NewIteration(seed int64, duration int, prepull int, players *[]*player.Player, downtimes *[]Downtime) *Iteration {
	r := rand.New(rand.NewSource(seed))
	enemy := enemy.NewEnemy("Test Enemy", r)
	// Create a new Iteration
	iteration := &Iteration{
		Seed:             seed,
		Rand:             r,
		Enemy:            enemy,
		Players:          nil,
		encounterTime:    0 - prepull,
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
		sumPlayerDamage += i.Players[j].Job.DamageDealt
	}
	if sumPlayerDamage == i.TotalDamage {
		// The total Damage dealt by the Players is equal to the total Damage taken by the Enemies, so the simulation was successful
	} else {
		// The total Damage dealt by the Players is not equal to the total Damage taken by the Enemies, so the simulation was not successful
		fmt.Printf("Iteration %d encountered an error: P: %d != E: %d\n", i.Seed, sumPlayerDamage, i.TotalDamage)
	}
}

func (i *Iteration) Tick() {
	i.Enemy.Tick(int64(i.encounterTime))
	var playerWG sync.WaitGroup
	playerWG.Add(len(i.Players))
	// Run the current Tick for all the Players
	for j := 0; j < len(i.Players); j++ {
		go i.Players[j].Tick(&playerWG, i.encounterTime, i.Enemy, i.Players)
	}
	// Wait for all the Players to finish
	playerWG.Wait()
}
