package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ATroschke/xivsim/models"
)

type SimulationResult struct {
	AverageDPS float64
	P0         *models.Encounter
	P25        *models.Encounter
	P50        *models.Encounter
	P75        *models.Encounter
	P100       *models.Encounter
}

func main() {
	// TODO: Get gear stats from Lodestone/Etro
	fmt.Println("Setting up player stats")
	// HACK: Currently assuming 6.48 Warrior BIS from https://etro.gg/gearset/1103c082-1c80-4bf3-bb56-83734971d5ea
	playerStats := models.NewPlayerStats(132, 3330, 2450, 2576, 940, 2182, 529, 400, 0, 382)

	fmt.Println("Setting up encounters")
	// Set a TickLength (how much simulated ms pass per tick)
	ticklength := 10 * time.Millisecond
	// Set an amount of encounters to simulate
	iterations := 1
	// Set the length of the encounters
	encounterLength := 3 * time.Minute
	// Set GOMAXPROCS to the number of CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Create encounters
	encounters := make([]*models.Encounter, iterations)
	for i := 0; i < iterations; i++ {
		encounters[i] = models.NewEncounter(i, 1, 1, encounterLength, ticklength, playerStats)
	}
	fmt.Printf("Simulating %d encounter(s) with a duration of %.2f minutes and a tickrate of %dms\n", iterations, encounterLength.Minutes(), ticklength.Milliseconds())
	// Get the current time
	start := time.Now()
	// Start encounters in batches of 10% of the total amount of encounters (if possible)
	batchSize := int(float64(iterations) * 0.1)
	if batchSize == 0 {
		batchSize = 1
	}

	for i := 0; i < iterations; i += batchSize {
		var wg sync.WaitGroup
		for j := i; j < i+batchSize && j < iterations; j++ {
			wg.Add(1)
			go encounters[j].Start(&wg)
		}
		wg.Wait()
		fmt.Printf("Finished %d encounters\n", i+batchSize)
	}
	// Build the result
	result := EncounterResult{
		AverageDPS: 0,
		P0:         encounters[0],
		P25:        encounters[0],
		P50:        encounters[0],
		P75:        encounters[0],
		P100:       encounters[0],
	}
	for _, encounter := range encounters {
		result.AverageDPS += encounter.GetDPS()
		if encounter.GetDPS() < result.P0.GetDPS() {
			result.P0 = encounter
		}
		if encounter.GetDPS() > result.P100.GetDPS() {
			result.P100 = encounter
		}
	}
	fmt.Printf("Average DPS: %.2f\n", result.AverageDPS/float64(iterations))
	fmt.Printf("Best DPS: %.2f Seed: %d\n", result.P100.GetDPS(), result.P100.GetSeed())
	fmt.Printf("Worst DPS: %.2f Seed: %d\n", result.P0.GetDPS(), result.P0.GetSeed())
	// Get the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Simulation took %d ms\n", elapsed.Milliseconds())
}
