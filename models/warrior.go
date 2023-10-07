package models

type Warrior struct {
	Name   string
	ID     int
	JobMod int
	Skills []*Skill
}

func NewWarrior() Job {
	return &Warrior{
		Name:   "Warrior",
		ID:     1,
		JobMod: 105,
		Skills: []*Skill{
			&HeavySwing,
			&Maim,
			&StormsEye,
			&StormsPath,
		},
	}
}

func (w *Warrior) GetName() string {
	return w.Name
}

func (w *Warrior) SelectNextSkill(player *Player, enemy *Enemy, encounterTime int64) *Skill {
	// If the player has a combo, try to continue it
	if player.NextCombo != nil {
		// If we can use either Storm's Eye or Storm's Path, check if Surging Tempest is > 3 seconds duration
		if player.NextCombo[0].ID == StormsEye.ID || player.NextCombo[0].ID == StormsPath.ID {
			for _, buff := range player.Buffs {
				if buff.Buff.ID == SurgingTempest.ID && buff.AppliedUntil-encounterTime > 3000 {
					return &StormsPath
				}
			}
		}
		return player.NextCombo[0]
	}
	return &HeavySwing
}

// Skills
var HeavySwing = Skill{
	Name:      "Heavy Swing",
	ID:        31,
	Potency:   200,
	GCD:       GCD2_5,
	NextCombo: []*Skill{&Maim},
}

var Maim = Skill{
	Name:      "Maim",
	ID:        37,
	Potency:   300,
	GCD:       GCD2_5,
	NextCombo: []*Skill{&StormsEye, &StormsPath},
}

var StormsEye = Skill{
	Name:      "Storm's Eye",
	ID:        45,
	Potency:   440,
	GCD:       GCD2_5,
	NextCombo: nil,
	AppliesBuffs: []Buff{
		SurgingTempest,
	},
}

var StormsPath = Skill{
	Name:      "Storm's Path",
	ID:        42,
	Potency:   440,
	GCD:       GCD2_5,
	NextCombo: nil,
}

// Buffs
var SurgingTempest = Buff{
	ID:        2677,
	Duration:  30000,
	ApplyType: ApplyTypeSelfExtend,
	DamageMod: 1.1,
}
