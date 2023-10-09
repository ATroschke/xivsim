package models

type Warrior struct {
	Name    string
	ID      int
	JobMod  int
	Skills  []*Skill
	Charges []Charge
	// Warrior specific resources
	BeastGauge int
}

func NewWarrior() Job {
	war := &Warrior{
		Name:   "Warrior",
		ID:     1,
		JobMod: 105,
		Skills: []*Skill{
			&HeavySwing,
			&Maim,
			&StormsEye,
			&StormsPath,
			&FellCleave,
			&Upheaval,
			&Onslaught,
			&InnerRelease,
			&Infuriate,
			&InnerChaos,
			&PrimalRend,
		},
		Charges: []Charge{},
	}
	// Initialize the charges (iterate over the skills and add any with 1 or more charges)
	for _, skill := range war.Skills {
		if skill.MaxCharges > 0 {
			war.Charges = append(war.Charges, Charge{
				SkillID:       skill.ID,
				SkillCooldown: skill.CooldownMS,
				Charges:       skill.MaxCharges,
				MaxCharges:    skill.MaxCharges,
				NextCharge:    0,
			})
		}
	}
	return war
}

func (w *Warrior) GetName() string {
	return w.Name
}

func (w *Warrior) UpdateCharges(encounterTime int64) {
	for i, charge := range w.Charges {
		if charge.NextCharge != 0 && encounterTime > charge.NextCharge && charge.Charges < charge.MaxCharges {
			w.Charges[i].Charges++
			w.Charges[i].NextCharge = encounterTime + int64(charge.SkillCooldown)
		}
	}
}

func (w *Warrior) SelectNextGCD(player *Player, enemy *Enemy, encounterTime int64) *Skill {
	var surgingTempestBuff *AppliedBuff
	// Check if the player has Surging Tempest
	for _, buff := range player.Buffs {
		if buff.Buff.ID == SurgingTempest.ID {
			surgingTempestBuff = &buff
			break
		}
	}
	// If the player has Nascent Chaos and Surging Tempest is up, use Inner Chaos
	for _, buff := range player.Buffs {
		if buff.Buff.ID == NascentChaos.ID && surgingTempestBuff != nil {
			return &InnerChaos
		}
	}
	// If the player has Primal Rend Ready and Surging Tempest is up, use Primal Rend
	for _, buff := range player.Buffs {
		if buff.Buff.ID == PrimalRendReady.ID && surgingTempestBuff != nil {
			return &PrimalRend
		}
	}
	// If the player has Inner Release and Surging Tempest is up, use Fell Cleave
	for _, buff := range player.Buffs {
		if buff.Buff.ID == InnerReleaseBuff.ID && surgingTempestBuff != nil {
			return &FellCleave
		}
	}
	// If the player has enough beast gauge for Fell Cleave and Surging Tempest is up, use it
	if w.BeastGauge >= 50 && surgingTempestBuff != nil {
		return &FellCleave
	}
	// If the player has a combo, try to continue it
	if player.NextCombo != nil {
		// If we can use either Storm's Eye or Storm's Path, check if Surging Tempest is > 7 seconds duration
		if player.NextCombo[0].ID == StormsEye.ID || player.NextCombo[0].ID == StormsPath.ID {
			if surgingTempestBuff != nil && surgingTempestBuff.AppliedUntil-encounterTime > 7000 {
				return &StormsPath
			} else {
				return &StormsEye
			}
		}
		return player.NextCombo[0]
	}
	return &HeavySwing
}

func (w *Warrior) SelectNextOGCD(player *Player, enemy *Enemy, encounterTime int64, GCDUntil int64) *Skill {
	var surgingTempestBuff *AppliedBuff
	var nascentChaosBuff *AppliedBuff
	// Check if the player has Surging Tempest or nascent chaos
	for _, buff := range player.Buffs {
		if buff.Buff.ID == SurgingTempest.ID {
			surgingTemp := buff
			surgingTempestBuff = &surgingTemp
		}
		if buff.Buff.ID == NascentChaos.ID {
			nascentTemp := buff
			nascentChaosBuff = &nascentTemp
		}
	}
	// If Infuriate has a charge and we don't have Nascent Chaos and would not overcap beast gauge, use it
	for i, charge := range w.Charges {
		if charge.SkillID == Infuriate.ID && charge.Charges > 0 && nascentChaosBuff == nil && w.BeastGauge <= 50 {
			// Make sure we don't use the OGCD if it would cause us to clip the GCD
			if encounterTime+int64(Infuriate.LockMS) < GCDUntil {
				w.Charges[i].Charges--
				w.Charges[i].NextCharge = encounterTime + int64(w.Charges[i].SkillCooldown)
				return &Infuriate
			}
		}
	}
	// If Inner Release has a charge and we have Surging Tempest up (with more than 7s), use it
	for i, charge := range w.Charges {
		if charge.SkillID == InnerRelease.ID && charge.Charges > 0 && surgingTempestBuff != nil && surgingTempestBuff.AppliedUntil-encounterTime > 7000 {
			// Make sure we don't use the OGCD if it would cause us to clip the GCD
			if encounterTime+int64(InnerRelease.LockMS) < GCDUntil {
				w.Charges[i].Charges--
				w.Charges[i].NextCharge = encounterTime + int64(w.Charges[i].SkillCooldown)
				return &InnerRelease
			}
		}
	}
	// If Upheaval has a charge, use it
	for i, charge := range w.Charges {
		if charge.SkillID == Upheaval.ID && charge.Charges > 0 && surgingTempestBuff != nil && surgingTempestBuff.AppliedUntil-encounterTime > 3000 {
			// Make sure we don't use the OGCD if it would cause us to clip the GCD
			if encounterTime+int64(Upheaval.LockMS) < GCDUntil {
				w.Charges[i].Charges--
				w.Charges[i].NextCharge = encounterTime + int64(w.Charges[i].SkillCooldown)
				return &Upheaval
			}
		}
	}
	// If Onslaught has a charge, use it
	for i, charge := range w.Charges {
		if charge.SkillID == Onslaught.ID && charge.Charges > 0 && surgingTempestBuff != nil && surgingTempestBuff.AppliedUntil-encounterTime > 3000 {
			// Make sure we don't use the OGCD if it would cause us to clip the GCD
			if encounterTime+int64(Upheaval.LockMS) < GCDUntil {
				w.Charges[i].Charges--
				w.Charges[i].NextCharge = encounterTime + int64(w.Charges[i].SkillCooldown)
				return &Onslaught
			}
		}
	}
	// Return nil if no OGCDs are available
	return nil
}

// Skills
var HeavySwing = Skill{
	Name:        "Heavy Swing",
	ID:          31,
	Potency:     200,
	GCD:         GCD2_5,
	BreaksCombo: true,
	LockMS:      700,
	NextCombo:   []*Skill{&Maim},
	CustomLogic: nil,
	MaxCharges:  0,
}

var Maim = Skill{
	Name:        "Maim",
	ID:          37,
	Potency:     300,
	GCD:         GCD2_5,
	BreaksCombo: true,
	LockMS:      700,
	NextCombo:   []*Skill{&StormsEye, &StormsPath},
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		j := job.(*Warrior)
		j.IncreaseBeastGauge(10)
		return false, false
	},
	MaxCharges: 0,
}

var StormsEye = Skill{
	Name:        "Storm's Eye",
	ID:          45,
	Potency:     440,
	GCD:         GCD2_5,
	NextCombo:   nil,
	BreaksCombo: true,
	LockMS:      700,
	AppliesBuffs: []Buff{
		SurgingTempest,
	},
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		j := job.(*Warrior)
		j.IncreaseBeastGauge(10)
		return false, false
	},
	MaxCharges: 0,
}

var StormsPath = Skill{
	Name:        "Storm's Path",
	ID:          42,
	Potency:     440,
	GCD:         GCD2_5,
	BreaksCombo: true,
	NextCombo:   nil,
	LockMS:      700,
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		j := job.(*Warrior)
		j.IncreaseBeastGauge(20)
		return false, false
	},
	MaxCharges: 0,
}

var FellCleave = Skill{
	Name:        "Fell Cleave",
	ID:          3549,
	Potency:     520,
	GCD:         GCD2_5,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      700,
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		// Check if the player is under the effect of Inner Release
		hasIR := false
		for i, buff := range player.Buffs {
			if buff.Buff.ID == InnerReleaseBuff.ID {
				hasIR = true
				// If so, remove one stack of Inner Release
				player.Buffs[i].Stacks--
				// If the buff has no stacks left, remove it
				if player.Buffs[i].Stacks <= 0 {
					player.Buffs = append(player.Buffs[:i], player.Buffs[i+1:]...)
				}
				// Break out of the loop
				break
			}
		}
		// If the player does not have Inner Release, decrease the beast gauge by 50
		if !hasIR {
			j := job.(*Warrior)
			j.DecreaseBeastGauge(50)
		}
		return hasIR, hasIR
	},
	MaxCharges: 0,
}

var InnerChaos = Skill{
	Name:        "Inner Chaos",
	ID:          16465,
	Potency:     660,
	GCD:         GCD2_5,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      700,
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		// Remove Nascent Chaos
		for i, buff := range player.Buffs {
			if buff.Buff.ID == NascentChaos.ID {
				player.Buffs = append(player.Buffs[:i], player.Buffs[i+1:]...)
				break
			}
		}
		return true, true
	},
	MaxCharges: 0,
}

var PrimalRend = Skill{
	Name:        "Primal Rend",
	ID:          25753,
	Potency:     700,
	GCD:         GCD2_5,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      1000,
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		// Remove the Primal Rend Ready buff
		for i, buff := range player.Buffs {
			if buff.Buff.ID == PrimalRendReady.ID {
				player.Buffs = append(player.Buffs[:i], player.Buffs[i+1:]...)
				break
			}
		}
		return true, true
	},
	MaxCharges: 0,
}

// oGCDs
var Upheaval = Skill{
	Name:        "Upheaval",
	ID:          7387,
	Potency:     400,
	GCD:         OGCD,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      700,
	CooldownMS:  30000,
	CustomLogic: nil,
	MaxCharges:  1,
}

var Onslaught = Skill{
	Name:        "Onslaught",
	ID:          7386,
	Potency:     150,
	GCD:         OGCD,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      700,
	CooldownMS:  30000,
	CustomLogic: nil,
	MaxCharges:  3,
}

var Infuriate = Skill{
	Name:        "Infuriate",
	ID:          52,
	Potency:     0,
	GCD:         OGCD,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      700,
	CooldownMS:  60000,
	AppliesBuffs: []Buff{
		NascentChaos,
	},
	CustomLogic: func(player *Player, job Job) (bool, bool) {
		j := job.(*Warrior)
		j.IncreaseBeastGauge(50)
		return false, false
	},
	MaxCharges: 2,
}

var InnerRelease = Skill{
	Name:        "Inner Release",
	ID:          7389,
	Potency:     0,
	GCD:         OGCD,
	BreaksCombo: false,
	NextCombo:   nil,
	LockMS:      700,
	CooldownMS:  60000,
	CustomLogic: nil,
	MaxCharges:  1,
	AppliesBuffs: []Buff{
		InnerReleaseBuff,
		PrimalRendReady,
	},
}

// Buffs
var SurgingTempest = Buff{
	ID:        2677,
	Duration:  30000,
	ApplyType: ApplyTypeSelfExtend,
	DamageMod: 1.1,
	Stacks:    0,
}

var InnerReleaseBuff = Buff{
	ID:        1177,
	Duration:  0, // Doesn't expire
	ApplyType: ApplyTypeSelf,
	DamageMod: 1,
	Stacks:    3,
}

var NascentChaos = Buff{
	ID:        1897,
	Duration:  30000,
	ApplyType: ApplyTypeSelf,
	DamageMod: 1,
	Stacks:    1,
}

var PrimalRendReady = Buff{
	ID:        2624,
	Duration:  0, // Doesn't expire
	ApplyType: ApplyTypeSelf,
	DamageMod: 1,
	Stacks:    1,
}

func (w *Warrior) IncreaseBeastGauge(amount int) {
	w.BeastGauge += amount
	// If the beast gauge would go over 100, set it to 100
	if w.BeastGauge > 100 {
		w.BeastGauge = 100
	}
}

func (w *Warrior) DecreaseBeastGauge(amount int) {
	w.BeastGauge -= amount
	// If the beast gauge would go below 0, panic
	if w.BeastGauge < 0 {
		panic("Beast gauge went below 0")
	}
}
