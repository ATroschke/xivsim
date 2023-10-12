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
	players := make([]*player.Player, 4)
	player0 := player.NewPlayer("WAR - 2.48 - NO MELDS", 90, "WAR", 132, 3360, 3337, 3687, 2344, 400, 1901, 502, 400, 529, 400, true)
	player1 := player.NewPlayer("WAR - 2.50 - BAKED EGGPLANT", 90, "WAR", 132, 3360, 3330, 3820, 2576, 940, 2182, 400, 400, 529, 400, true)
	player2 := player.NewPlayer("WAR - 2.45 - BABA GHANOUSH", 90, "WAR", 132, 3360, 3340, 3833, 2627, 904, 1901, 671, 400, 529, 400, true)
	player3 := player.NewPlayer("WAR - 2.40 - BABA GHANOUSH", 90, "WAR", 132, 3360, 3340, 3833, 2627, 616, 1901, 960, 400, 529, 400, true)
	players[0] = player0
	players[1] = player1
	players[2] = player2
	players[3] = player3
	// Create a new Encounter
	encounter := encounter.NewEncounter("WAR BiS Comparison", players, 600, 10, 50, 1000)
	encounter.Run()
	encounter.Report()
}
