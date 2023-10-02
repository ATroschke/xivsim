package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ATroschke/xivsim/models"
)

func main() {
	// TODO: Get gear stats from Lodestone/Etro
	// HACK: Currently assuming 6.48 Warrior BIS from https://etro.gg/gearset/1103c082-1c80-4bf3-bb56-83734971d5ea
	playerStats := models.NewPlayerStats(132, 3330, 2450, 2576, 940, 2182, 529, 400, 0, 382)
	// Set a TickLength (how much simulated ms pass per tick)
	ticklength := 100 * time.Millisecond
	// Set an amount of encounters to simulate
	iterations := 50000
	// Set GOMAXPROCS to the number of CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Create a waitgroup to wait for all encounters to finish
	var wg sync.WaitGroup
	wg.Add(iterations)
	// Create encounters
	encounters := make([]*models.Encounter, iterations)
	for i := 0; i < iterations; i++ {
		encounters[i] = models.NewEncounter(i, 1, 1, 10*time.Second, ticklength, playerStats)
	}
	fmt.Printf("Starting %d encounters with a Tick of %dms\n", iterations, ticklength.Milliseconds())
	// Get the current time
	start := time.Now()
	// Start all encounters
	for _, encounter := range encounters {
		go encounter.Start(&wg)
	}
	fmt.Println("Waiting for encounters to finish")
	wg.Wait()
	// Calculate Average DPS, lowest DPS, and highest DPS
	var totalDPS float64
	var lowestDPS float64
	var highestDPS float64
	for _, encounter := range encounters {
		totalDPS += encounter.GetDPS()
		if encounter.GetDPS() < lowestDPS || lowestDPS == 0 {
			lowestDPS = encounter.GetDPS()
		}
		if encounter.GetDPS() > highestDPS {
			highestDPS = encounter.GetDPS()
		}
	}
	avgDPS := totalDPS / float64(iterations)
	fmt.Printf("Average DPS: %.2f\nHighest DPS: %.2f\nLowest DPS: %.2f\n", avgDPS, highestDPS, lowestDPS)
	// Get the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Simulation took %d ms\n", elapsed.Milliseconds())
}
