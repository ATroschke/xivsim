package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	models "github.com/ATroschke/xivsim/Models"
	"github.com/gin-gonic/gin"
)

// A Go Gin based webserver that receives a request for a simulation and returns the result

func main() {
	e := gin.Default()
	e.POST("/simulate", simulate)
	e.Run()
}

type SimulationRequest struct {
	Iterations   int
	DurationMS   int
	TickLengthMS int
	Player       PlayerRequest
}

type SimulationResult struct {
	AverageDPS float64
	P0         *models.Encounter
	P25        *models.Encounter
	P50        *models.Encounter
	P75        *models.Encounter
	P100       *models.Encounter
}

type PlayerRequest struct {
	Job   string
	Stats models.PlayerStats
}

func simulate(c *gin.Context) {
	var request SimulationRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encounters := make([]*models.Encounter, request.Iterations)
	for i := 0; i < request.Iterations; i++ {
		encounters[i] = models.NewEncounter(i, 1, 1, time.Duration(request.DurationMS), time.Duration(request.TickLengthMS), &request.Player.Stats)
	}

	// Get the current time
	start := time.Now()
	// Start encounters in batches of 10% of the total amount of encounters (if possible)
	batchSize := int(float64(request.Iterations) * 0.1)
	if batchSize == 0 {
		batchSize = 1
	}

	for i := 0; i < request.Iterations; i += batchSize {
		var wg sync.WaitGroup
		for j := i; j < i+batchSize && j < request.Iterations; j++ {
			wg.Add(1)
			go encounters[j].Start(&wg)
		}
		wg.Wait()
		fmt.Printf("Finished %d encounters\n", i+batchSize)
	}

	// Build the result
	result := SimulationResult{
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
		if encounter.GetDPS() < result.P25.GetDPS() {
			result.P25 = encounter
		}
		if encounter.GetDPS() < result.P50.GetDPS() {
			result.P50 = encounter
		}
		if encounter.GetDPS() < result.P75.GetDPS() {
			result.P75 = encounter
		}
		if encounter.GetDPS() < result.P100.GetDPS() {
			result.P100 = encounter
		}
	}

	// Get the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Simulation took %d ms\n", elapsed.Milliseconds())

	c.JSON(http.StatusOK, result)
}
