package job

import (
	"github.com/ATroschke/xivsim/sim/buff"
	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Gunbreaker implements Job
type Gunbreaker struct {
	Skills      GunbreakerSkills
	Buffs       GunbreakerBuffs
	DOTs        GunbreakerDOTs
	PowderGauge int
}

type GunbreakerSkills struct {
	KeenEdge      skill.Skill
	BrutalShell   skill.Skill
	SolidBarrel   skill.Skill
	GnashingFang  skill.Skill
	SavageClaw    skill.Skill
	WickedTalon   skill.Skill
	JugularRip    skill.Skill
	AbdomenTear   skill.Skill
	EyeGouge      skill.Skill
	NoMercy       skill.Skill
	BowShock      skill.Skill
	SonicBreak    skill.Skill
	BlastingZone  skill.Skill
	Bloodfest     skill.Skill
	DoubleDown    skill.Skill
	BurstStrike   skill.Skill
	RoughDivide   skill.Skill
	Hypervelocity skill.Skill
	Tincture      skill.Skill
}

type GunbreakerBuffs struct {
	NoMercy     buff.Buff
	ReadyToRip  buff.Buff
	ReadyToTear buff.Buff
	ReadyToMaim buff.Buff
	ReadToBlast buff.Buff
	Tincture    buff.Buff
}

type GunbreakerDOTs struct {
	BowShock   buff.DOT
	SonicBreak buff.DOT
}

func NewGunbreaker() *Gunbreaker {
	return &Gunbreaker{
		Skills: GunbreakerSkills{
			KeenEdge:      KeenEdge,
			BrutalShell:   BrutalShell,
			SolidBarrel:   SolidBarrel,
			GnashingFang:  GnashingFang,
			SavageClaw:    SavageClaw,
			WickedTalon:   WickedTalon,
			JugularRip:    JugularRip,
			AbdomenTear:   AbdomenTear,
			EyeGouge:      EyeGouge,
			NoMercy:       NoMercy,
			BowShock:      BowShock,
			SonicBreak:    SonicBreak,
			BlastingZone:  BlastingZone,
			Bloodfest:     Bloodfest,
			DoubleDown:    DoubleDown,
			BurstStrike:   BurstStrike,
			RoughDivide:   RoughDivide,
			Hypervelocity: Hypervelocity,
			Tincture:      GNBTincture,
		},
		Buffs: GunbreakerBuffs{
			NoMercy:     NoMercyBuff,
			ReadyToRip:  ReadyToRip,
			ReadyToTear: ReadyToTear,
			ReadyToMaim: ReadyToMaim,
			ReadToBlast: ReadyToBlast,
			Tincture:    TinctureBuff,
		},
		DOTs: GunbreakerDOTs{
			BowShock:   BowShockDOT,
			SonicBreak: SonicBreakDOT,
		},
		PowderGauge: 0,
	}
}

func CopyGunbreaker(j *Job) *Gunbreaker {
	g := j.JobImpl.(*Gunbreaker)
	// Create a new instance of the Gunbreaker
	return &Gunbreaker{
		Skills: GunbreakerSkills{
			KeenEdge:      g.Skills.KeenEdge,
			BrutalShell:   g.Skills.BrutalShell,
			SolidBarrel:   g.Skills.SolidBarrel,
			GnashingFang:  g.Skills.GnashingFang,
			SavageClaw:    g.Skills.SavageClaw,
			WickedTalon:   g.Skills.WickedTalon,
			JugularRip:    g.Skills.JugularRip,
			AbdomenTear:   g.Skills.AbdomenTear,
			EyeGouge:      g.Skills.EyeGouge,
			NoMercy:       g.Skills.NoMercy,
			BowShock:      g.Skills.BowShock,
			SonicBreak:    g.Skills.SonicBreak,
			BlastingZone:  g.Skills.BlastingZone,
			Bloodfest:     g.Skills.Bloodfest,
			DoubleDown:    g.Skills.DoubleDown,
			BurstStrike:   g.Skills.BurstStrike,
			RoughDivide:   g.Skills.RoughDivide,
			Hypervelocity: g.Skills.Hypervelocity,
			Tincture:      g.Skills.Tincture,
		},
		Buffs: GunbreakerBuffs{
			NoMercy:     g.Buffs.NoMercy,
			ReadyToRip:  g.Buffs.ReadyToRip,
			ReadyToTear: g.Buffs.ReadyToTear,
			ReadyToMaim: g.Buffs.ReadyToMaim,
			ReadToBlast: g.Buffs.ReadToBlast,
			Tincture:    g.Buffs.Tincture,
		},
		DOTs: GunbreakerDOTs{
			BowShock:   g.DOTs.BowShock,
			SonicBreak: g.DOTs.SonicBreak,
		},
		PowderGauge: g.PowderGauge,
	}
}

func (g *Gunbreaker) GetSkills() []*skill.Skill {
	return []*skill.Skill{
		&g.Skills.KeenEdge,
		&g.Skills.BrutalShell,
		&g.Skills.SolidBarrel,
		&g.Skills.GnashingFang,
		&g.Skills.SavageClaw,
		&g.Skills.WickedTalon,
		&g.Skills.JugularRip,
		&g.Skills.AbdomenTear,
		&g.Skills.EyeGouge,
		&g.Skills.NoMercy,
		&g.Skills.BowShock,
		&g.Skills.SonicBreak,
		&g.Skills.BlastingZone,
		&g.Skills.Bloodfest,
		&g.Skills.DoubleDown,
		&g.Skills.BurstStrike,
		&g.Skills.RoughDivide,
		&g.Skills.Hypervelocity,
		&g.Skills.Tincture,
	}
}

func (g *Gunbreaker) GetBuffs() []*buff.Buff {
	return []*buff.Buff{
		&g.Buffs.NoMercy,
		&g.Buffs.ReadyToRip,
		&g.Buffs.ReadyToTear,
		&g.Buffs.ReadyToMaim,
		&g.Buffs.ReadToBlast,
		&g.Buffs.Tincture,
	}
}

func (g *Gunbreaker) GetDots() []*buff.DOT {
	return []*buff.DOT{
		&g.DOTs.BowShock,
		&g.DOTs.SonicBreak,
	}
}

func (g *Gunbreaker) GetJobMod() (int, int) {
	return 156, 100
}

func (g *Gunbreaker) NextSkill(j *Job, encounterTime int64) *skill.Skill {
	if encounterTime >= 0 {
		if encounterTime < j.AnimationLockUntil {
			return nil
		}
		if encounterTime >= j.GCDUntil {
			return g.SelectNextGCD(j, encounterTime)
		}
		// Check if the default animation lock would cut into the GCD
		if encounterTime+700 >= j.GCDUntil {
			return nil
		}
		return g.SelectNextOGCD(j, encounterTime)
	} else {
		// Gunbreaker has no Prepull
		return nil
	}
}

func (g *Gunbreaker) SelectNextGCD(j *Job, encounterTime int64) *skill.Skill {
	// If we have No Mercy active and Gnashing Fang is on Cooldown, use Sonic Break
	if g.Buffs.NoMercy.AppliedUntil > encounterTime && g.Skills.GnashingFang.Charges == 0 && g.Skills.SonicBreak.Charges > 0 {
		return &g.Skills.SonicBreak
	}
	// If we have No Mercy active and 2 or more Cartridges, use Double Down
	if g.Buffs.NoMercy.AppliedUntil > encounterTime && g.PowderGauge >= 2 && g.Skills.DoubleDown.Charges > 0 {
		return &g.Skills.DoubleDown
	}
	if j.NextCombo != nil {
		if j.NextCombo[0].Name == BrutalShell.Name {
			return &g.Skills.BrutalShell
		}
		if j.NextCombo[0].Name == SolidBarrel.Name {
			return &g.Skills.SolidBarrel
		}
		if j.NextCombo[0].Name == SavageClaw.Name {
			return &g.Skills.SavageClaw
		}
		if j.NextCombo[0].Name == WickedTalon.Name {
			return &g.Skills.WickedTalon
		}
	}
	// If we have a cartridge, use Gnashing Fang (Does this properly Loop?)
	if g.PowderGauge >= 1 && g.Skills.GnashingFang.Charges > 0 {
		return &g.Skills.GnashingFang
	}
	// If we have No Mercy active, use Burst Strike
	if g.Buffs.NoMercy.AppliedUntil > encounterTime && g.PowderGauge >= 1 {
		return &g.Skills.BurstStrike
	}
	return &g.Skills.KeenEdge
}

func (g *Gunbreaker) SelectNextOGCD(j *Job, encounterTime int64) *skill.Skill {
	// OPENER Relevant: If our next GCD is Solid Barrel and we can use Pot, do so
	if j.NextCombo != nil && j.NextCombo[0].Name == SolidBarrel.Name && g.Skills.Tincture.Charges > 0 && encounterTime < 30000 {
		return &g.Skills.Tincture
	}
	// If No Mercy is about to be ready, Pot
	if g.Skills.NoMercy.NextCharge+5000 > encounterTime && g.Skills.Tincture.Charges > 0 {
		return &g.Skills.Tincture
	}
	// If we have No Mercy active, use Bow Shock
	if g.Buffs.NoMercy.AppliedUntil > encounterTime && g.Skills.BowShock.Charges > 0 {
		return &g.Skills.BowShock
	}
	// If we have Ready to Rip, use Jugular Rip
	if g.Buffs.ReadyToRip.Stacks > 0 {
		return &g.Skills.JugularRip
	}
	// If we have Ready to Tear, use Abdomen Tear
	if g.Buffs.ReadyToTear.Stacks > 0 {
		return &g.Skills.AbdomenTear
	}
	// If we have Ready to Maim, use Eye Gouge
	if g.Buffs.ReadyToMaim.Stacks > 0 {
		return &g.Skills.EyeGouge
	}
	// If we have Ready to Blast, use Hypervelocity
	if g.Buffs.ReadToBlast.Stacks > 0 {
		return &g.Skills.Hypervelocity
	}
	// If we have No Mercy active or it is on cooldown for more than 25 seconds, use Blasting Zone
	if (g.Buffs.NoMercy.AppliedUntil > encounterTime || g.Skills.NoMercy.NextCharge-encounterTime > 25000) && g.Skills.BlastingZone.Charges > 0 {
		return &g.Skills.BlastingZone
	}
	// If we have No Mercy active, use Bloodfest
	if g.Buffs.NoMercy.AppliedUntil > encounterTime && g.Skills.Bloodfest.Charges > 0 {
		return &g.Skills.Bloodfest
	}
	// If we have No Mercy active, use Rough Divide
	if g.Buffs.NoMercy.AppliedUntil > encounterTime && g.Skills.RoughDivide.Charges > 0 {
		return &g.Skills.RoughDivide
	}
	// Try to use NoMercy as late in the GCD as possible
	if g.PowderGauge >= 1 && g.Skills.NoMercy.Charges > 0 && j.GCDUntil-encounterTime < 900 {
		return &g.Skills.NoMercy
	}
	return nil
}

func (g *Gunbreaker) ValidateAutoCDH(s *skill.Skill) bool {
	return false
}

func (g *Gunbreaker) GetBuffModifiers(encounterTime int64) float64 {
	if g.Buffs.NoMercy.AppliedUntil > encounterTime {
		return g.Buffs.NoMercy.DamageMod
	}
	return 1.0
}

func (g *Gunbreaker) GetDesiredPrepullDuration() int {
	return 0
}

func (g *Gunbreaker) CheckPotionStatus(time int64) bool {
	return g.Buffs.Tincture.AppliedUntil >= time
	//return false
}

func (g *Gunbreaker) AddCartridge(amount int) {
	g.PowderGauge += amount
	if g.PowderGauge > 3 {
		g.PowderGauge = 3
	}
}

func (g *Gunbreaker) SpendCartridge(amount int) {
	g.PowderGauge -= amount
	if g.PowderGauge < 0 {
		g.PowderGauge = 0
	}
}

// Gunbreaker Skills
var (
	// GCD
	KeenEdge = skill.Skill{
		Name:        "Keen Edge",
		Potency:     200,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&BrutalShell},
	}
	BrutalShell = skill.Skill{
		Name:        "Brutal Shell",
		Potency:     300,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&SolidBarrel},
	}
	SolidBarrel = skill.Skill{
		Name:        "Solid Barrel",
		Potency:     360,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Add a Cartridge
			g.AddCartridge(1)
		},
	}
	GnashingFang = skill.Skill{
		Name:        "Gnashing Fang",
		Potency:     380,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  1,
		Charges:     1,
		NextCombo:   []*skill.Skill{&SavageClaw},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Spend a Cartridge
			g.SpendCartridge(1)
			// Add Ready to Rip
			g.Buffs.ReadyToRip.AppliedUntil = time + ReadyToRip.DurationMS
			g.Buffs.ReadyToRip.Stacks++
		},
	}
	SavageClaw = skill.Skill{
		Name:        "Savage Claw",
		Potency:     460,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&WickedTalon},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Add Ready to Tear
			g.Buffs.ReadyToTear.AppliedUntil = time + ReadyToTear.DurationMS
			g.Buffs.ReadyToTear.Stacks++
		},
	}
	WickedTalon = skill.Skill{
		Name:        "Wicked Talon",
		Potency:     540,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Add Ready to Maim
			g.Buffs.ReadyToMaim.AppliedUntil = time + ReadyToMaim.DurationMS
			g.Buffs.ReadyToMaim.Stacks++
		},
	}
	SonicBreak = skill.Skill{
		Name:        "Sonic Break",
		Potency:     300,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			j := job.(*Job)
			g := j.JobImpl.(*Gunbreaker)
			// Check if No Mercy is active
			buffMod := 1.0
			if g.Buffs.NoMercy.AppliedUntil > time {
				buffMod = g.Buffs.NoMercy.DamageMod
			}
			// Apply the DOT
			target.ApplyDot(&g.DOTs.SonicBreak, &j.DamageDealt, time, buffMod)
		},
	}
	DoubleDown = skill.Skill{
		Name:        "Double Down",
		Potency:     1200,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			j := job.(*Job)
			g := j.JobImpl.(*Gunbreaker)
			// Spend 2 Cartridges
			g.SpendCartridge(2)
		},
	}
	BurstStrike = skill.Skill{
		Name:        "Burst Strike",
		Potency:     380,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Spend a Cartridge
			g.SpendCartridge(1)
			// Add Ready to Blast
			g.Buffs.ReadToBlast.AppliedUntil = time + ReadyToBlast.DurationMS
			g.Buffs.ReadToBlast.Stacks++
		},
	}
	// OGCD
	NoMercy = skill.Skill{
		Name:        "No Mercy",
		Potency:     0,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Apply the Buff
			g.Buffs.NoMercy.AppliedUntil = time + NoMercyBuff.DurationMS
		},
	}
	BowShock = skill.Skill{
		Name:        "Bow Shock",
		Potency:     150,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			j := job.(*Job)
			g := j.JobImpl.(*Gunbreaker)
			// Check if No Mercy is active
			buffMod := 1.0
			if g.Buffs.NoMercy.AppliedUntil > time {
				buffMod = g.Buffs.NoMercy.DamageMod
			}
			// Apply the DOT
			target.ApplyDot(&g.DOTs.BowShock, &j.DamageDealt, time, buffMod)
		},
	}
	JugularRip = skill.Skill{
		Name:        "Jugular Rip",
		Potency:     200,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Remove Ready to Rip
			g.Buffs.ReadyToRip.AppliedUntil = 0
			g.Buffs.ReadyToRip.Stacks = 0
		},
	}
	AbdomenTear = skill.Skill{
		Name:        "Abdomen Tear",
		Potency:     240,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Remove Ready to Tear
			g.Buffs.ReadyToTear.AppliedUntil = 0
			g.Buffs.ReadyToTear.Stacks = 0
		},
	}
	EyeGouge = skill.Skill{
		Name:        "Eye Gouge",
		Potency:     280,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Remove Ready to Maim
			g.Buffs.ReadyToMaim.AppliedUntil = 0
			g.Buffs.ReadyToMaim.Stacks = 0
		},
	}
	BlastingZone = skill.Skill{
		Name:        "Blasting Zone",
		Potency:     720,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  1,
		Charges:     1,
	}
	Bloodfest = skill.Skill{
		Name:        "Bloodfest",
		Potency:     0,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  120000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Add 3 Cartridges
			g.AddCartridge(3)
		},
	}
	RoughDivide = skill.Skill{
		Name:        "Rough Divide",
		Potency:     150,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      1000,
		CooldownMS:  30000,
		MaxCharges:  2,
		Charges:     2,
	}
	Hypervelocity = skill.Skill{
		Name:        "Hypervelocity",
		Potency:     180,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Remove Ready to Blast
			g.Buffs.ReadToBlast.AppliedUntil = 0
			g.Buffs.ReadToBlast.Stacks = 0
		},
	}
	GNBTincture = skill.Skill{
		Name:       "Grade 8 Tincture of Strength",
		GCD:        skill.OGCD,
		LockMS:     1300,
		CooldownMS: 270000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			g := job.(*Job).JobImpl.(*Gunbreaker)
			// Apply Pot Buff
			g.Buffs.Tincture.AppliedUntil = time + g.Buffs.Tincture.DurationMS
		},
	}
	// DOTs
	BowShockDOT = buff.DOT{
		Name:       "Bow Shock",
		Potency:    60,
		DurationMS: 15000,
	}
	SonicBreakDOT = buff.DOT{
		Name:       "Sonic Break",
		Potency:    60,
		DurationMS: 30000,
	}
	// Buffs
	NoMercyBuff = buff.Buff{
		Name:       "No Mercy",
		DamageMod:  1.25,
		DurationMS: 20000,
	}
	ReadyToRip = buff.Buff{
		Name:       "Ready to Rip",
		Stacks:     0,
		MaxStacks:  1,
		DurationMS: 10000,
	}
	ReadyToTear = buff.Buff{
		Name:       "Ready to Tear",
		Stacks:     0,
		MaxStacks:  1,
		DurationMS: 10000,
	}
	ReadyToMaim = buff.Buff{
		Name:       "Ready to Maim",
		Stacks:     0,
		MaxStacks:  1,
		DurationMS: 10000,
	}
	ReadyToBlast = buff.Buff{
		Name:       "Ready to Blast",
		Stacks:     0,
		MaxStacks:  1,
		DurationMS: 10000,
	}
)
