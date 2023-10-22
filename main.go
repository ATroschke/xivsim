package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ATroschke/xivsim/api"
	"github.com/ATroschke/xivsim/sim/encounter"
	"github.com/ATroschke/xivsim/sim/player"
)

/*
	main.go is used as a helper entry-point during development and won't be used in the actual API.
	The Idea is to easily run a simulation locally without setting up a server.
	Simply run `go mod tidy` (once) and `go run main.go` to start a simulation.
*/

func dev() {
	// Set GOMAXPROCS to the number of CPUs available, we want to punch in as much performance as we can
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
	gnb244, err := api.GetFromEtro("b3847bc9-2fda-4b7b-a8a1-e1720b51a46e")
	if err != nil {
		panic(err)
	}
	drk250, err := api.GetFromEtro("dcd2eb34-7c43-4840-a17b-2eb790f19cf4")
	if err != nil {
		panic(err)
	}
	// Create a new List of Players
	players := make([]*player.Player, 4)
	war := player.NewPlayer("WAR - 2.50 - 1103c082-1c80-4bf3-bb56-83734971d5ea", 90, "WAR", war250.WeaponDamage, 3360, war250.MainStat, war250.Vitality, war250.CriticalHit, war250.DirectHit, war250.Determination, war250.SkillSpeed, war250.SpellSpeed, war250.Tenacity, war250.Piety, true)
	pld := player.NewPlayer("PLD - 2.50 - 3a2d77ff-57e1-434d-b1da-72a7d4f44944", 90, "PLD", pld250.WeaponDamage, 2240, pld250.MainStat, pld250.Vitality, pld250.CriticalHit, pld250.DirectHit, pld250.Determination, pld250.SkillSpeed, pld250.SpellSpeed, pld250.Tenacity, pld250.Piety, true)
	gnb := player.NewPlayer("GNB - 2.44 - b3847bc9-2fda-4b7b-a8a1-e1720b51a46e", 90, "GNB", gnb244.WeaponDamage, 2800, gnb244.MainStat, gnb244.Vitality, gnb244.CriticalHit, gnb244.DirectHit, gnb244.Determination, gnb244.SkillSpeed, gnb244.SpellSpeed, gnb244.Tenacity, gnb244.Piety, true)
	drk := player.NewPlayer("DRK - 2.50 - dcd2eb34-7c43-4840-a17b-2eb790f19cf4", 90, "DRK", drk250.WeaponDamage, 2960, drk250.MainStat, drk250.Vitality, drk250.CriticalHit, drk250.DirectHit, drk250.Determination, drk250.SkillSpeed, drk250.SpellSpeed, drk250.Tenacity, drk250.Piety, true)
	players[0] = war
	players[1] = pld
	players[2] = gnb
	players[3] = drk
	// WIP: Specify downtimes
	downtimes := make([]encounter.Downtime, 0)
	/*downtimes[0] = encounter.Downtime{
		Start: -1,
		End:   -1,
	}*/
	// Create a new Encounter
	encounter := encounter.NewEncounter("Tank Sim", players, 620, 10, 50, 1000, &downtimes)
	encounter.Run()
	report := encounter.Report()
	printReport(&report)
}

// Helper function to Print the Report and optionally a specific rotation for debug/development purposes
func printReport(report *encounter.EncounterResults) {
	fmt.Printf("%s\n", report.Name)
	fmt.Printf("Duration: %d\n", report.Duration)
	// Format Iterations with dots
	formattedIterations := fmt.Sprintf("%d", report.NrIterations)
	for i := len(formattedIterations) - 3; i > 0; i -= 3 {
		formattedIterations = formattedIterations[:i] + "." + formattedIterations[i:]
	}
	fmt.Printf("Iterations: %s\n", formattedIterations)
	simulationTime := time.Duration(report.SimulationTime) * time.Millisecond
	fmt.Printf("Simulated in: %s\n", simulationTime.String())
	fmt.Printf("Margin of Error: %.2f -  %.4f%%\n", report.MarginOfError, report.ConfidencePct)
	fmt.Printf("Confidence Interval: %.2f - %.2f\n", report.ConfidenceLow, report.ConfidenceHigh)
	for _, player := range report.PlayerResults {
		fmt.Printf("%s: DPS: %d (min:%d - max:%d)\n", player.Player, player.AverageDPS, player.LowestDPS, player.HighestDPS)
		if player.Player == "DRK - 2.50 - dcd2eb34-7c43-4840-a17b-2eb790f19cf4" && false {
			// Print DRK Rotation
			fmt.Printf("DRK Rotation:\n")
			for _, skill := range player.HighestDamageLog {
				// Want to filter out Auto Attacks
				if skill.Name != "Attack" {
					// Format the skill Damage value
					skillDamage := fmt.Sprintf("%d", skill.Damage)
					if skill.Crit && skill.DirectHit {
						skillDamage += "!!"
					} else if skill.Crit {
						skillDamage += "!"
					} else if skill.DirectHit {
						skillDamage += "*"
					}
					fmt.Printf("%d: %s - %s\n", skill.Time, skill.Name, skillDamage)
				}
			}
		}
	}
}
