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
	Iterations   int           `json:"iterations" binding:"required"`
	DurationS    int           `json:"durationS" binding:"required"`
	TickLengthMS int           `json:"tickLengthMS" binding:"required"`
	Player       PlayerRequest `json:"player" binding:"required"`
}

type SimulationResult struct {
	AverageDPS float64
	P0         *models.Encounter
	P100       *models.Encounter
	SimTimeMS  int64
}

type PlayerRequest struct {
	Job   string             `json:"job" binding:"required"`
	Stats PlayerStatsRequest `json:"stats" binding:"required"`
}

type PlayerStatsRequest struct {
	// Stats from Gear etc
	WeaponDamage    int `json:"weaponDamage" binding:"required"`
	Mainstat        int `json:"mainstat" binding:"required"`
	CriticalHitRate int `json:"criticalHitRate" binding:"required"`
	DirectHitRate   int `json:"directHitRate" binding:"required"`
	Determination   int `json:"determination" binding:"required"`
	Tenacity        int `json:"tenacity" binding:"required"`
	Piety           int `json:"piety" binding:"required"`
	Speed           int `json:"speed" binding:"required"`
	Vitality        int `json:"vitality" binding:"required"`
}

func simulate(c *gin.Context) {
	var request SimulationRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't allow more than 500 iterations
	if request.Iterations > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many iterations (max 500)"})
		return
	}

	// Don't allow more than 10 minutes of simulation time
	if request.DurationS > 600 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too long duration (max 10 minutes)"})
		return
	}

	// Build player Stats from request
	ps := models.NewPlayerStats(request.Player.Stats.WeaponDamage, request.Player.Stats.Mainstat, request.Player.Stats.CriticalHitRate, request.Player.Stats.DirectHitRate, request.Player.Stats.Determination, request.Player.Stats.Tenacity, request.Player.Stats.Piety, request.Player.Stats.Speed, 0, request.Player.Stats.Vitality)

	tickLength := time.Duration(request.TickLengthMS) * time.Millisecond
	duration := time.Duration(request.DurationS) * time.Second

	encounters := make([]*models.Encounter, request.Iterations)
	for i := 0; i < request.Iterations; i++ {
		encounters[i] = models.NewEncounter(i, 1, 1, duration, tickLength, ps)
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
	}

	// Build the result
	result := SimulationResult{
		AverageDPS: 0,
		P0:         encounters[0],
		P100:       encounters[0],
	}
	// Handle calculations for specific percentiles
	for _, encounter := range encounters {
		result.AverageDPS += encounter.DPS
		if encounter.DPS < result.P0.DPS {
			result.P0 = encounter
		}
		if encounter.DPS > result.P100.DPS {
			result.P100 = encounter
		}
	}
	result.AverageDPS /= float64(request.Iterations)
	// Round to 2 decimals
	result.AverageDPS = float64(int(result.AverageDPS*100)) / 100

	// Display Average and Calculated Crit/Direct Rates
	result.P0.CalculateCritDirectRate()
	result.P100.CalculateCritDirectRate()

	// Get the elapsed time
	elapsed := time.Since(start)
	fmt.Printf("Simulation took %d ms\n", elapsed.Milliseconds())
	result.SimTimeMS = elapsed.Milliseconds()

	c.JSON(http.StatusOK, result)
}
