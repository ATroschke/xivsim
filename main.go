package main

import (
	"runtime"

	"github.com/ATroschke/xivsim/api"
	"github.com/ATroschke/xivsim/sim/encounter"
	"github.com/ATroschke/xivsim/sim/player"
)

func main() {
	// TODO: Formulae kontrollieren für NO MELDS set sollte 100 potency 2497 damage machen
	// Set GOMAXPROCS to the number of CPUs available
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Read a Gearset from Etro.gg
	war250, err := api.GetFromEtro("1103c082-1c80-4bf3-bb56-83734971d5ea")
	if err != nil {
		panic(err)
	}
	pld250, err := api.GetFromEtro("3a2d77ff-57e1-434d-b1da-72a7d4f44944")
	if err != nil {
		panic(err)
	}
	// Create a new List of Players
	players := make([]*player.Player, 2)
	player1 := player.NewPlayer("WAR - 2.50 - 1103c082-1c80-4bf3-bb56-83734971d5ea", 90, "WAR", war250.WeaponDamage, 3360, war250.MainStat, war250.Vitality, war250.CriticalHit, war250.DirectHit, war250.Determination, war250.SkillSpeed, war250.SpellSpeed, war250.Tenacity, war250.Piety, true)
	player2 := player.NewPlayer("PLD - 2.50 - 3a2d77ff-57e1-434d-b1da-72a7d4f44944", 90, "WAR", pld250.WeaponDamage, 2240, pld250.MainStat, pld250.Vitality, pld250.CriticalHit, pld250.DirectHit, pld250.Determination, pld250.SkillSpeed, pld250.SpellSpeed, pld250.Tenacity, pld250.Piety, true)
	players[0] = player1
	players[1] = player2
	// Specify downtimes
	downtimes := make([]encounter.Downtime, 0)
	/*downtimes[0] = encounter.Downtime{
		Start: -1,
		End:   -1,
	}*/
	// Create a new Encounter
	encounter := encounter.NewEncounter("WAR BiS Comparison", players, 60, 10, 50, 1, &downtimes)
	encounter.Run()
	encounter.Report()
}
