package main

import (
	"runtime"

	"github.com/ATroschke/xivsim/sim/encounter"
	"github.com/ATroschke/xivsim/sim/player"
)

func main() {
	// Set GOMAXPROCS to the number of CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Create a new List of Players
	players := make([]*player.Player, 1)
	player1 := player.NewPlayer("Test Player 1", 90, "WAR", 132, 3330, 3820, 2576, 940, 2182, 400, 400, 529, 400)
	players[0] = player1
	// Create a new Encounter
	encounter := encounter.NewEncounter("Test Encounter", players, 180, 0, 10, 1)
	encounter.Run()
	encounter.Report()
}
