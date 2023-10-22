package job

import (
	"github.com/ATroschke/xivsim/sim/buff"
	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

type Darkknight struct {
	Skills     DarkknightSkills
	Buffs      DarkknightBuffs
	DOTs       DarkknightDOTs
	MP         int
	BloodGauge int
	LSH        LivingShadowHandler
}

type DarkknightSkills struct {
	HardSlash       skill.Skill
	SyphonStrike    skill.Skill
	SoulEater       skill.Skill
	Bloodspiller    skill.Skill
	EdgeOfShadow    skill.Skill
	Delirium        skill.Skill
	BloodWeapon     skill.Skill
	LivingShadow    skill.Skill
	SaltedEarth     skill.Skill
	Shadowbringer   skill.Skill
	CarveAndSpit    skill.Skill
	Plunge          skill.Skill
	SaltAndDarkness skill.Skill
	// Living Shadow Skills
	LSAbyssalDrain  skill.Skill
	LSPlunge        skill.Skill
	LSFloodOfShadow skill.Skill
	LSEdgeOfShadow  skill.Skill
	LSBloodspiller  skill.Skill
	LSShadowbringer skill.Skill
	// Tincture
	Tincture skill.Skill
}

type DarkknightBuffs struct {
	Darkside      buff.Buff
	Delirium      buff.Buff
	BloodWeapon   buff.Buff
	SaltedEarthUP buff.Buff
	Tincture      buff.Buff
}

type DarkknightDOTs struct {
	SaltedEarth buff.DOT
}

func NewDarkknight() *Darkknight {
	return &Darkknight{
		Skills: DarkknightSkills{
			HardSlash:       HardSlash,
			SyphonStrike:    SyphonStrike,
			SoulEater:       SoulEater,
			Bloodspiller:    Bloodspiller,
			EdgeOfShadow:    EdgeOfShadow,
			Delirium:        Delirium,
			BloodWeapon:     BloodWeapon,
			LivingShadow:    LivingShadow,
			SaltedEarth:     SaltedEarth,
			Shadowbringer:   Shadowbringer,
			CarveAndSpit:    CarveAndSpit,
			Plunge:          Plunge,
			SaltAndDarkness: SaltAndDarkness,
			LSAbyssalDrain:  LSAbyssalDrain,
			LSPlunge:        LSPlunge,
			LSFloodOfShadow: LSFloodOfShadow,
			LSEdgeOfShadow:  LSEdgeOfShadow,
			LSBloodspiller:  LSBloodspiller,
			LSShadowbringer: LSShadowbringer,
			Tincture:        DRKTincture,
		},
		Buffs: DarkknightBuffs{
			Darkside:      Darkside,
			Delirium:      DeliriumBuff,
			BloodWeapon:   BloodWeaponBuff,
			SaltedEarthUP: SaltedEarthUP,
			Tincture:      TinctureBuff,
		},
		DOTs: DarkknightDOTs{
			SaltedEarth: SaltedEarthDOT,
		},
		MP:         10000,
		BloodGauge: 0,
		LSH: LivingShadowHandler{
			NextAction: 0,
			Step:       0,
		},
	}
}

func CopyDarkknight(j *Job) *Darkknight {
	d := j.JobImpl.(*Darkknight)
	return &Darkknight{
		Skills: DarkknightSkills{
			HardSlash:       d.Skills.HardSlash,
			SyphonStrike:    d.Skills.SyphonStrike,
			SoulEater:       d.Skills.SoulEater,
			Bloodspiller:    d.Skills.Bloodspiller,
			EdgeOfShadow:    d.Skills.EdgeOfShadow,
			Delirium:        d.Skills.Delirium,
			BloodWeapon:     d.Skills.BloodWeapon,
			LivingShadow:    d.Skills.LivingShadow,
			SaltedEarth:     d.Skills.SaltedEarth,
			Shadowbringer:   d.Skills.Shadowbringer,
			CarveAndSpit:    d.Skills.CarveAndSpit,
			Plunge:          d.Skills.Plunge,
			SaltAndDarkness: d.Skills.SaltAndDarkness,
			LSAbyssalDrain:  d.Skills.LSAbyssalDrain,
			LSPlunge:        d.Skills.LSPlunge,
			LSFloodOfShadow: d.Skills.LSFloodOfShadow,
			LSEdgeOfShadow:  d.Skills.LSEdgeOfShadow,
			LSBloodspiller:  d.Skills.LSBloodspiller,
			LSShadowbringer: d.Skills.LSShadowbringer,
			Tincture:        d.Skills.Tincture,
		},
		Buffs: DarkknightBuffs{
			Darkside:      d.Buffs.Darkside,
			Delirium:      d.Buffs.Delirium,
			BloodWeapon:   d.Buffs.BloodWeapon,
			SaltedEarthUP: d.Buffs.SaltedEarthUP,
			Tincture:      d.Buffs.Tincture,
		},
		DOTs: DarkknightDOTs{
			SaltedEarth: d.DOTs.SaltedEarth,
		},
		MP:         d.MP,
		BloodGauge: d.BloodGauge,
		LSH:        d.LSH,
	}
}

func (d *Darkknight) GetSkills() []*skill.Skill {
	return []*skill.Skill{
		&d.Skills.HardSlash,
		&d.Skills.SyphonStrike,
		&d.Skills.SoulEater,
		&d.Skills.Bloodspiller,
		&d.Skills.EdgeOfShadow,
		&d.Skills.Delirium,
		&d.Skills.BloodWeapon,
		&d.Skills.LivingShadow,
		&d.Skills.SaltedEarth,
		&d.Skills.Shadowbringer,
		&d.Skills.CarveAndSpit,
		&d.Skills.Plunge,
		&d.Skills.SaltAndDarkness,
		&d.Skills.LSAbyssalDrain,
		&d.Skills.LSPlunge,
		&d.Skills.LSFloodOfShadow,
		&d.Skills.LSEdgeOfShadow,
		&d.Skills.LSBloodspiller,
		&d.Skills.LSShadowbringer,
		&d.Skills.Tincture,
	}
}

func (d *Darkknight) GetBuffs() []*buff.Buff {
	return []*buff.Buff{
		&d.Buffs.Darkside,
		&d.Buffs.Delirium,
		&d.Buffs.BloodWeapon,
		&d.Buffs.SaltedEarthUP,
		&d.Buffs.Tincture,
	}
}

func (d *Darkknight) GetDots() []*buff.DOT {
	return []*buff.DOT{
		&d.DOTs.SaltedEarth,
	}
}

func (d *Darkknight) GetJobMod() (int, int) {
	return 156, 105
}

func (d *Darkknight) NextSkill(job *Job, encounterTime int64) *skill.Skill {
	if encounterTime >= 0 {
		// Check if LS can use a Skill, and if it can, return that
		lsSkill := d.SelectNextLSSkill(encounterTime)
		if lsSkill != nil {
			return lsSkill
		}
		// Check if the player is animation locked
		if encounterTime < job.AnimationLockUntil {
			return nil
		}
		// Check if the GCD is ready
		if encounterTime >= job.GCDUntil {
			return d.SelectNextGCD(job, encounterTime)
		}
		// Check if the default animation lock (700ms) would cut into the GCD
		if encounterTime+700 >= job.GCDUntil {
			return nil
		}
		return d.SelectNextOGCD(encounterTime)
	} else {
		// TODO: DRK Prepull
		if encounterTime >= -4000 && d.Skills.BloodWeapon.Charges > 0 {
			return &d.Skills.BloodWeapon
		}
		if encounterTime >= -2000 && d.Skills.Tincture.Charges > 0 {
			return &d.Skills.Tincture
		}
		return nil
	}
}

func (d *Darkknight) SelectNextGCD(job *Job, encounterTime int64) *skill.Skill {
	// Check if we have Delirium Stacks, if yes, use Bloodspiller
	// BUT prioritize getting Living Shadow out first
	if d.Buffs.Delirium.Stacks > 0 && d.Skills.LivingShadow.Charges == 0 {
		return &d.Skills.Bloodspiller
	}
	// Use a Bloodspiller if we can and have Pot up
	if d.BloodGauge >= 50 && d.Buffs.Tincture.AppliedUntil >= encounterTime {
		return &d.Skills.Bloodspiller
	}
	if job.NextCombo != nil {
		if job.NextCombo[0].Name == SyphonStrike.Name {
			return &d.Skills.SyphonStrike
		}
		if job.NextCombo[0].Name == SoulEater.Name {
			// Would we overcap on Bloodgauge with this Soul Eater?
			if d.BloodGauge >= 80 {
				// Yeah, so we should use a Bloodspiller
				return &d.Skills.Bloodspiller
			}
			return &d.Skills.SoulEater
		}
	}

	return &d.Skills.HardSlash
}

func (d *Darkknight) SelectNextOGCD(encounterTime int64) *skill.Skill {
	// Check if Darkside is up/about to expire > Edge of Shadow
	if (d.Buffs.Darkside.AppliedUntil <= encounterTime || d.Buffs.Darkside.AppliedUntil <= encounterTime+3000) && d.MP >= 3000 && d.Skills.EdgeOfShadow.Charges > 0 {
		return &d.Skills.EdgeOfShadow
	}
	// Pot if possible and Blood Weapon is about to be ready
	if d.Skills.BloodWeapon.NextCharge+5000 > encounterTime && d.Skills.Tincture.Charges > 0 {
		return &d.Skills.Tincture
	}
	// Blood Weapon is ready
	if d.Skills.BloodWeapon.Charges > 0 {
		return &d.Skills.BloodWeapon
	}
	// Delirium is ready
	if d.Skills.Delirium.Charges > 0 {
		return &d.Skills.Delirium
	}
	// Living Shadow is ready
	if d.Skills.LivingShadow.Charges > 0 && d.BloodGauge >= 50 {
		return &d.Skills.LivingShadow
	}
	// Salted Earth is ready and Darkside is active
	if d.Skills.SaltedEarth.Charges > 0 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.SaltedEarth
	}
	// Shadowbringer is ready and has 2 Stacks
	if d.Skills.Shadowbringer.Charges == 2 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.Shadowbringer
	}
	// Edge of Shadow if MP would overcap soon
	if d.MP >= 8000 && d.Skills.EdgeOfShadow.Charges > 0 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.EdgeOfShadow
	}
	// Carve and Spit is ready
	if d.Skills.CarveAndSpit.Charges > 0 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.CarveAndSpit
	}
	// Plunge is ready and has 2 Stacks
	if d.Skills.Plunge.Charges == 2 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.Plunge
	}
	// Shadowbringer is ready
	if d.Skills.Shadowbringer.Charges > 0 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.Shadowbringer
	}
	// Edge of Shadow can be used
	if d.MP >= 3000 && d.Skills.EdgeOfShadow.Charges > 0 && d.Buffs.Darkside.AppliedUntil > encounterTime {
		return &d.Skills.EdgeOfShadow
	}
	// Salt and Darkness can be used
	if d.Skills.SaltAndDarkness.Charges > 0 && d.Buffs.SaltedEarthUP.Stacks > 0 && d.Buffs.SaltedEarthUP.AppliedUntil > encounterTime {
		return &d.Skills.SaltAndDarkness
	}
	// Plunge is ready
	if d.Skills.Plunge.Charges > 0 {
		return &d.Skills.Plunge
	}
	return nil
}

func (d *Darkknight) SelectNextLSSkill(encounterTime int64) *skill.Skill {
	// Manage Living Shadow Skills....
	if d.LSH.Step != 0 && encounterTime > d.LSH.NextAction {
		if d.LSH.Step == 1 {
			return &d.Skills.LSAbyssalDrain
		}
		if d.LSH.Step == 2 {
			return &d.Skills.LSPlunge
		}
		if d.LSH.Step == 3 {
			return &d.Skills.LSFloodOfShadow
		}
		if d.LSH.Step == 4 {
			return &d.Skills.LSEdgeOfShadow
		}
		if d.LSH.Step == 5 {
			return &d.Skills.LSBloodspiller
		}
		if d.LSH.Step == 6 {
			return &d.Skills.LSShadowbringer
		}
	}
	return nil
}

func (d *Darkknight) ValidateAutoCDH(s *skill.Skill) bool {
	return s.AutoCDH
}

func (d *Darkknight) GetBuffModifiers(encounterTime int64) float64 {
	if d.Buffs.Darkside.AppliedUntil > encounterTime {
		return d.Buffs.Darkside.DamageMod
	}
	return 1
}

func (d *Darkknight) GetDesiredPrepullDuration() int {
	return 4000
}

func (d *Darkknight) CheckPotionStatus(time int64) bool {
	return d.Buffs.Tincture.AppliedUntil >= time
}

func (d *Darkknight) IncreaseMP(amount int) {
	d.MP += amount
	if d.MP > 10000 {
		d.MP = 10000
	}
}

func (d *Darkknight) DecreaseMP(amount int) {
	d.MP -= amount
	if d.MP < 0 {
		d.MP = 0
	}
}

func (d *Darkknight) IncreaseBloodGauge(amount int) {
	d.BloodGauge += amount
	if d.BloodGauge > 100 {
		d.BloodGauge = 100
	}
}

func (d *Darkknight) DecreaseBloodGauge(amount int) {
	d.BloodGauge -= amount
	if d.BloodGauge < 0 {
		d.BloodGauge = 0
	}
}

func (d *Darkknight) BloodWeapon() {
	// This function will check if Bloodweapon is up and handle it's effects
	// It should be plugged in every GCD's CustomLogic
	if d.Buffs.BloodWeapon.Stacks > 0 {
		d.Buffs.BloodWeapon.Stacks--
		d.IncreaseBloodGauge(10)
		d.IncreaseMP(600)
	}
}

// Darkknight Skills
var (
	// GCDs
	HardSlash = skill.Skill{
		Name:        "Hard Slash",
		Potency:     170,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&SyphonStrike},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.BloodWeapon()
		},
	}
	SyphonStrike = skill.Skill{
		Name:        "Syphon Strike",
		Potency:     260,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&SoulEater},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.IncreaseMP(600)
			d.BloodWeapon()
		},
	}
	SoulEater = skill.Skill{
		Name:        "Soul Eater",
		Potency:     340,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.IncreaseBloodGauge(20)
			d.BloodWeapon()
		},
	}
	Bloodspiller = skill.Skill{
		Name:        "Bloodspiller",
		Potency:     500,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			// Check if we have Delirium Stacks
			if d.Buffs.Delirium.Stacks > 0 {
				// Remove a stack and restore some MP
				d.Buffs.Delirium.Stacks--
				d.IncreaseMP(200)
			} else {
				d.DecreaseBloodGauge(50)
			}
			d.BloodWeapon()
		},
	}
	// OGCDs
	EdgeOfShadow = skill.Skill{
		Name:        "Edge of Shadow",
		Potency:     460,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		NextCombo:   nil,
		CooldownMS:  1000,
		Charges:     1,
		MaxCharges:  1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.Buffs.Darkside.AppliedUntil = time + d.Buffs.Darkside.DurationMS
			d.DecreaseMP(3000)
		},
	}
	Delirium = skill.Skill{
		Name:       "Delirium",
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 60000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.Buffs.Delirium.Stacks = d.Buffs.Delirium.MaxStacks
		},
	}
	BloodWeapon = skill.Skill{
		Name:       "Blood Weapon",
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 60000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.Buffs.BloodWeapon.Stacks = d.Buffs.BloodWeapon.MaxStacks
		},
	}
	LivingShadow = skill.Skill{
		Name:       "Living Shadow",
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 120000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.DecreaseBloodGauge(50)
			d.LSH.NextAction = time + 6800
			d.LSH.Step = 1
		},
	}
	SaltedEarth = skill.Skill{
		Name:       "Salted Earth",
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 90000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			j := job.(*Job)
			d := j.JobImpl.(*Darkknight)
			// Set the internal Buff to Signal that Salted Earth is now UP
			d.Buffs.SaltedEarthUP.Stacks = d.Buffs.SaltedEarthUP.MaxStacks
			d.Buffs.SaltedEarthUP.AppliedUntil = time + d.Buffs.SaltedEarthUP.DurationMS
			// Check if Darkside if active
			buffMod := 1.0
			if d.Buffs.Darkside.AppliedUntil > time {
				buffMod = d.Buffs.Darkside.DamageMod
			}
			// Apply the DoT
			target.ApplyDot(&d.DOTs.SaltedEarth, &j.DamageDealt, time, buffMod)
		},
	}
	Shadowbringer = skill.Skill{
		Name:       "Shadowbringer",
		Potency:    600,
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 60000,
		Charges:    2,
		MaxCharges: 2,
	}
	CarveAndSpit = skill.Skill{
		Name:       "Carve and Spit",
		Potency:    510,
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 60000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.IncreaseMP(600)
		},
	}
	Plunge = skill.Skill{
		Name:       "Plunge",
		Potency:    150,
		GCD:        skill.OGCD,
		LockMS:     1000,
		CooldownMS: 30000,
		Charges:    2,
		MaxCharges: 2,
	}
	SaltAndDarkness = skill.Skill{
		Name:       "Salt and Darkness",
		Potency:    500,
		GCD:        skill.OGCD,
		LockMS:     700,
		CooldownMS: 20000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.Buffs.SaltedEarthUP.Stacks--
		},
	}
	DRKTincture = skill.Skill{
		Name:       "Grade 8 Tincture of Strength",
		GCD:        skill.OGCD,
		LockMS:     1300,
		CooldownMS: 270000,
		Charges:    1,
		MaxCharges: 1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			// Apply Pot Buff
			d.Buffs.Tincture.AppliedUntil = time + d.Buffs.Tincture.DurationMS
		},
	}
	// DOTs
	SaltedEarthDOT = buff.DOT{
		Name:       "Salted Earth",
		DurationMS: 15000,
		Potency:    50,
	}
	// Buffs
	Darkside = buff.Buff{
		Name:         "Darkside",
		DamageMod:    1.1,
		DurationMS:   30000,
		AppliedUntil: 0,
	}
	DeliriumBuff = buff.Buff{
		Name:       "Delirium",
		DurationMS: 15000,
		Stacks:     0,
		MaxStacks:  3,
	}
	BloodWeaponBuff = buff.Buff{
		Name:       "Blood Weapon",
		DurationMS: 15000,
		Stacks:     0,
		MaxStacks:  5,
	}
	SaltedEarthUP = buff.Buff{
		// This buff is internal and handles the prerequisite for Salt and Darkness
		Name:       "INTERNAL: Salted Earth Up",
		DurationMS: 15000,
		Stacks:     0,
		MaxStacks:  1,
	}
	// TODO: Calculate those properly... (what the hell did SE think when developin this)
	// Living Shadow Actions
	LSAbyssalDrain = skill.Skill{
		Name:    "Abyssal Dran (Living Shadow)",
		Potency: 350,
		GCD:     skill.OGCD, // Making sure this doesn't trigger the GCD
		LockMS:  0,          // Making sure this doesn't trigger the Animation Lock
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.LSH.Step++
			d.LSH.NextAction = time + 2180
		},
		MOverride: 195,
	}
	LSPlunge = skill.Skill{
		Name:    "Plunge (Living Shadow)",
		Potency: 350,
		GCD:     skill.OGCD,
		LockMS:  0,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.LSH.Step++
			d.LSH.NextAction = time + 2180
		},
		MOverride: 195,
	}
	LSFloodOfShadow = skill.Skill{
		Name:    "Flood of Shadow (Living Shadow)",
		Potency: 350,
		GCD:     skill.OGCD,
		LockMS:  0,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.LSH.Step++
			d.LSH.NextAction = time + 2180
		},
		MOverride: 195,
	}
	LSEdgeOfShadow = skill.Skill{
		Name:    "Edge of Shadow (Living Shadow)",
		Potency: 350,
		GCD:     skill.OGCD,
		LockMS:  0,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.LSH.Step++
			d.LSH.NextAction = time + 2180
		},
		MOverride: 195,
	}
	LSBloodspiller = skill.Skill{
		Name:    "Bloodspiller (Edge of Shadow)",
		Potency: 350,
		GCD:     skill.OGCD,
		LockMS:  0,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.LSH.Step++
			d.LSH.NextAction = time + 2180
		},
		MOverride: 195,
	}
	LSShadowbringer = skill.Skill{
		Name:    "Shadowbringer (Living Shadow)",
		Potency: 500,
		GCD:     skill.OGCD,
		LockMS:  0,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			d := job.(*Job).JobImpl.(*Darkknight)
			d.LSH.Step = 0
			d.LSH.NextAction = 0
		},
		MOverride: 195,
	}
)

// Helper class to handle Living Shadows "Rotation"
type LivingShadowHandler struct {
	NextAction int64
	Step       int
}
