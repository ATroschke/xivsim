package player

import (
	"sync"

	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/job"
)

// Player is a struct that contains all the information about a player.

type Player struct {
	ID          int // Player ID to identify the Player in the Encounter
	Name        string
	Job         job.Job
	Level       int
	Stats       *Stats
	Speed       *job.Speed
	DamageDealt int
	// TODO: ? Gear (Struct that contains all the gear)
}

// NewPlayer creates a new player with the given name and level.
func NewPlayer(name string, level int, playerJob string, weapondamage int, weapondelay int, mainstat int, vitality int, criticalhit int, directhit int, determination int, skillspeed int, spellspeed int, tenacity int, piety int, partybonus bool) *Player {
	player := &Player{
		Name:  name,
		Level: level,
	}

	if partybonus {
		mainstat = int(float64(mainstat) * 1.05)
	}

	// Set and Calculate stats
	player.Stats = &Stats{
		WeaponDamage:  weapondamage,
		MainStat:      mainstat,
		Vitality:      vitality,
		CriticalHit:   criticalhit,
		DirectHit:     directhit,
		Determination: determination,
		SkillSpeed:    skillspeed,
		SpellSpeed:    spellspeed,
		Tenacity:      tenacity,
		Piety:         piety,
	}

	// TODO: Get LvSub and LvDiv from a table, based on the Level of the Player
	player.Speed = job.NewSpeed(player.Stats.SkillSpeed, weapondelay, 400, 1900)

	player.Job = *job.NewJob(playerJob, player.Speed)

	player.Job.CalculateSkills(
		player.Stats.WeaponDamage,
		player.Stats.MainStat,
		player.Stats.CriticalHit,
		player.Stats.DirectHit,
		player.Stats.Determination,
		player.Stats.SkillSpeed,
		player.Stats.SpellSpeed,
		player.Stats.Tenacity,
	)

	return player
}

// Create a new instance of a player with the same stats as the given player.
// We don't want to modify the original player (since we may run multiple simulations at once)
// as this would cause issues with Resources and Buffs. We also don't want to calculate all the static values again.
func CopyPlayer(player *Player) *Player {
	job := player.Job.CopyJob()
	p := &Player{
		Name:  player.Name,
		Level: player.Level,
		Stats: player.Stats,
		Speed: player.Speed,
		// We need to create a new instance of the Job, since otherwise instances would share a GCD and Animation Lock, and resources etc.
		Job: *job,
	}
	return p
}

// Runs the current Tick of the Player.
func (p *Player) Tick(tickWG *sync.WaitGroup, encounterTime int, enemy *enemy.Enemy, players []*Player) {
	defer tickWG.Done()
	// Run the current Tick of the Player's Job
	aaDamage, skillDamage := p.Job.Tick(enemy, int64(encounterTime))
	p.DamageDealt += (aaDamage + skillDamage)
}
