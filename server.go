package main

import (
	"runtime"
	"time"

	"github.com/ATroschke/xivsim/api"
	"github.com/ATroschke/xivsim/sim/encounter"
	"github.com/ATroschke/xivsim/sim/player"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type SimRequest struct {
	Players  []PlayerRequest `json:"players" binding:"required"`
	Duration int             `json:"duration" binding:"required"`
	Ping     int             `json:"ping" binding:"required"`
	TickRate int             `json:"tickRate" binding:"required"`
	Iters    int             `json:"iters" binding:"required"`
}

type PlayerRequest struct {
	Name          string `json:"name" binding:"required"`
	Job           string `json:"job" binding:"required"`
	WeaponDelayMS int    `json:"weaponDelayMS" binding:"required"`
	EtroID        string `json:"etroID" binding:"required"`
}

func main() {
	// Set GOMAXPROCS to the number of CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())

	g := gin.Default()
	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "http://168.119.56.82:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))
	p := ginprometheus.NewPrometheus("xivSim")
	p.Use(g)

	g.POST("/sim", func(c *gin.Context) {
		var simRequest SimRequest
		if err := c.ShouldBindJSON(&simRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		players := make([]*player.Player, len(simRequest.Players))
		for i, playerRequest := range simRequest.Players {
			etroStats, err := api.GetFromEtro(playerRequest.EtroID)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			players[i] = player.NewPlayer(playerRequest.Name, 90, playerRequest.Job, etroStats.WeaponDamage, playerRequest.WeaponDelayMS, etroStats.MainStat, etroStats.Vitality, etroStats.CriticalHit, etroStats.DirectHit, etroStats.Determination, etroStats.SkillSpeed, etroStats.SpellSpeed, etroStats.Tenacity, etroStats.Piety, true)
		}

		// TODO: Implement downtimes
		downtimes := make([]encounter.Downtime, 0)
		// Get the current time as a Name String
		timeString := time.Now().Format("2006-01-02 15:04:05")
		// Make sure Duration and Iterations are within bounds
		if simRequest.Duration < 30 {
			simRequest.Duration = 30
		} else if simRequest.Duration > 1800 {
			simRequest.Duration = 1800
		}
		if simRequest.Iters < 2 {
			simRequest.Iters = 2
		} else if simRequest.Iters > 10000 {
			simRequest.Iters = 10000
		}
		// Create a new Encounter
		encounter := encounter.NewEncounter(timeString, players, simRequest.Duration, simRequest.Ping, simRequest.TickRate, simRequest.Iters, &downtimes)
		encounter.Run()
		c.JSON(200, encounter.Report())
	})

	g.Run(":8080")
}
