package job

import (
	"github.com/ATroschke/xivsim/sim/buff"
	"github.com/ATroschke/xivsim/sim/enemy"
	"github.com/ATroschke/xivsim/sim/skill"
)

// Paladin implements Job
type Paladin struct {
	Skills PaladinSkills
	Buffs  PaladinBuffs
	DOTs   PaladinDOTs
	MP     int
}

type PaladinSkills struct {
	// GCDs
	FastBlade             skill.Skill
	RiotBlade             skill.Skill
	RoyalAuthority        skill.Skill
	Atonement             skill.Skill
	GoringBlade           skill.Skill
	Confiteor             skill.Skill
	BladeOfFaith          skill.Skill
	BladeOfTruth          skill.Skill
	BladeOfValor          skill.Skill
	HolySpirit            skill.Skill
	HolySpiritReq         skill.Skill
	HolySpiritDivineMight skill.Skill
	// OGCDs
	FightOrFlight skill.Skill
	Requiescat    skill.Skill
	CircleOfScorn skill.Skill
	Expiacion     skill.Skill
	Intervene     skill.Skill
}

type PaladinBuffs struct {
	SwordOath     buff.Buff
	FightOrFlight buff.Buff
	Requiescat    buff.Buff
	Confiteor     buff.Buff
	DivineMight   buff.Buff
	Tincture      buff.Buff
}

type PaladinDOTs struct {
	CircleOfScorn buff.DOT
}

func NewPaladin() *Paladin {
	return &Paladin{
		Skills: PaladinSkills{
			FastBlade:             FastBlade,
			RiotBlade:             RiotBlade,
			RoyalAuthority:        RoyalAuthority,
			Atonement:             Atonement,
			FightOrFlight:         FightOrFlight,
			GoringBlade:           GoringBlade,
			Requiescat:            Requiescat,
			CircleOfScorn:         CircleOfScorn,
			Expiacion:             Expiacion,
			Intervene:             Intervene,
			Confiteor:             Confiteor,
			BladeOfFaith:          BladeOfFaith,
			BladeOfTruth:          BladeOfTruth,
			BladeOfValor:          BladeOfValor,
			HolySpirit:            HolySpirit,
			HolySpiritReq:         HolySpiritRequiescat,
			HolySpiritDivineMight: HolySpiritDivineMight,
		},
		Buffs: PaladinBuffs{
			SwordOath:     SwordOath,
			FightOrFlight: FightOrFlightBuff,
			Requiescat:    RequiescatBuff,
			Confiteor:     ConfiteorReady,
			DivineMight:   DivineMight,
			Tincture:      TinctureBuff,
		},
		DOTs: PaladinDOTs{
			CircleOfScorn: CircleOfScornDOT,
		},
		MP: 10000,
	}
}

func CopyPaladin(j *Job) *Paladin {
	p := j.JobImpl.(*Paladin)
	// Create a new instance of the Paladin
	return &Paladin{
		Skills: PaladinSkills{
			FastBlade:             p.Skills.FastBlade,
			RiotBlade:             p.Skills.RiotBlade,
			RoyalAuthority:        p.Skills.RoyalAuthority,
			Atonement:             p.Skills.Atonement,
			FightOrFlight:         p.Skills.FightOrFlight,
			GoringBlade:           p.Skills.GoringBlade,
			Requiescat:            p.Skills.Requiescat,
			CircleOfScorn:         p.Skills.CircleOfScorn,
			Expiacion:             p.Skills.Expiacion,
			Intervene:             p.Skills.Intervene,
			Confiteor:             p.Skills.Confiteor,
			BladeOfFaith:          p.Skills.BladeOfFaith,
			BladeOfTruth:          p.Skills.BladeOfTruth,
			BladeOfValor:          p.Skills.BladeOfValor,
			HolySpirit:            p.Skills.HolySpirit,
			HolySpiritReq:         p.Skills.HolySpiritReq,
			HolySpiritDivineMight: p.Skills.HolySpiritDivineMight,
		},
		Buffs: PaladinBuffs{
			SwordOath:     p.Buffs.SwordOath,
			FightOrFlight: p.Buffs.FightOrFlight,
			Requiescat:    p.Buffs.Requiescat,
			Confiteor:     p.Buffs.Confiteor,
			DivineMight:   p.Buffs.DivineMight,
			Tincture:      p.Buffs.Tincture,
		},
		DOTs: PaladinDOTs{
			CircleOfScorn: p.DOTs.CircleOfScorn,
		},
		MP: p.MP,
	}
}

func (p *Paladin) GetSkills() []*skill.Skill {
	return []*skill.Skill{
		&p.Skills.FastBlade,
		&p.Skills.RiotBlade,
		&p.Skills.RoyalAuthority,
		&p.Skills.Atonement,
		&p.Skills.FightOrFlight,
		&p.Skills.GoringBlade,
		&p.Skills.Requiescat,
		&p.Skills.CircleOfScorn,
		&p.Skills.Expiacion,
		&p.Skills.Intervene,
		&p.Skills.Confiteor,
		&p.Skills.BladeOfFaith,
		&p.Skills.BladeOfTruth,
		&p.Skills.BladeOfValor,
		&p.Skills.HolySpirit,
		&p.Skills.HolySpiritReq,
		&p.Skills.HolySpiritDivineMight,
	}
}

func (p *Paladin) GetBuffs() []*buff.Buff {
	return []*buff.Buff{
		&p.Buffs.SwordOath,
		&p.Buffs.FightOrFlight,
		&p.Buffs.Requiescat,
		&p.Buffs.Confiteor,
		&p.Buffs.DivineMight,
		&p.Buffs.Tincture,
	}
}

func (p *Paladin) GetDots() []*buff.DOT {
	return []*buff.DOT{
		&p.DOTs.CircleOfScorn,
	}
}

func (p *Paladin) GetJobMod() (int, int) {
	return 156, 100
}

func (p *Paladin) NextSkill(j *Job, encounterTime int64) *skill.Skill {
	if encounterTime >= 0 {
		if encounterTime < j.AnimationLockUntil {
			return nil
		}
		if encounterTime >= j.GCDUntil {
			return p.SelectNextGCD(j, encounterTime)
		}
		// Check if the default animation lock would cut into the GCD
		if encounterTime+700 >= j.GCDUntil {
			return nil
		}
		return p.SelectNextOGCD(encounterTime)
	} else {
		// Prepull
		// TODO: Casting and Casting Animation Lock
		// Cast Holy Spirit so it Hits around when the fight starts (0)
		if p.MP == 10000 && (encounterTime+2500 < 60 && encounterTime+2500 > -60) {
			return &p.Skills.HolySpirit
		}
		return nil
	}
}

func (p *Paladin) SelectNextGCD(job *Job, encounterTime int64) *skill.Skill {
	// Check if we are in a combo
	if job.NextCombo != nil {
		// We are in a combo, so we will select the next skill in the combo
		if job.NextCombo[0].Name == RiotBlade.Name {
			return &p.Skills.RiotBlade
		}
		if job.NextCombo[0].Name == RoyalAuthority.Name {
			return &p.Skills.RoyalAuthority
		}
		if job.NextCombo[0].Name == BladeOfFaith.Name {
			return &p.Skills.BladeOfFaith
		}
		if job.NextCombo[0].Name == BladeOfTruth.Name {
			return &p.Skills.BladeOfTruth
		}
		if job.NextCombo[0].Name == BladeOfValor.Name {
			return &p.Skills.BladeOfValor
		}
	}
	// If Goring Blade is ready and we have Fight or Flight, we will use Goring Blade
	if p.Skills.GoringBlade.Charges > 0 && p.Buffs.FightOrFlight.AppliedUntil > encounterTime {
		return &p.Skills.GoringBlade
	}
	// If we have Requiescat and Confiteor Ready, we will use Confiteor
	if p.Buffs.Confiteor.Stacks > 0 && p.Buffs.Requiescat.Stacks > 0 {
		return &p.Skills.Confiteor
	}
	// If we have Sword Oath, we will use Atonement
	if p.Buffs.SwordOath.Stacks > 0 {
		return &p.Skills.Atonement
	}
	// If we have Divine Might, we will use Holy Spirit
	if p.Buffs.DivineMight.Stacks > 0 {
		return &p.Skills.HolySpiritDivineMight
	}
	// If we have a Stack of Requiescat, we will use Holy Spirit
	if p.Buffs.Requiescat.Stacks > 0 {
		return &p.Skills.HolySpiritReq
	}
	// The Player is not in a combo, so we will select our first filler GCD
	return &p.Skills.FastBlade
}

func (p *Paladin) SelectNextOGCD(encounterTime int64) *skill.Skill {
	// If FoF is ready, we have 3 Stacks of Sword Oath,
	// Goring Blade is ready in the next GCD and Requiestcat is almost ready,
	// we will use FoF
	if p.Skills.FightOrFlight.Charges > 0 && p.Buffs.SwordOath.Stacks == 3 && (p.Skills.Requiescat.NextCharge <= encounterTime+2500 || p.Skills.Requiescat.Charges > 0) && (p.Skills.GoringBlade.NextCharge <= encounterTime+2500 || p.Skills.GoringBlade.Charges > 0) {
		return &p.Skills.FightOrFlight
	}
	// If Requiescat is ready and we have FoF, we will use Requiescat
	if p.Skills.Requiescat.Charges > 0 && p.Buffs.FightOrFlight.AppliedUntil > encounterTime {
		return &p.Skills.Requiescat
	}
	// If we have FoF or FoF is more than 25s away, we will use Circle of Scorn
	if p.Skills.CircleOfScorn.Charges > 0 && (p.Buffs.FightOrFlight.AppliedUntil > encounterTime || p.Skills.FightOrFlight.NextCharge > encounterTime+25000) {
		return &p.Skills.CircleOfScorn
	}
	// If we have FoF or FoF is more than 25s away, we will use Expiacion
	if p.Skills.Expiacion.Charges > 0 && (p.Buffs.FightOrFlight.AppliedUntil > encounterTime || p.Skills.FightOrFlight.NextCharge > encounterTime+25000) {
		return &p.Skills.Expiacion
	}
	// If we have FoF and Intervene is ready, we will use Intervene
	if p.Skills.Intervene.Charges > 0 && p.Buffs.FightOrFlight.AppliedUntil > encounterTime {
		return &p.Skills.Intervene
	}
	return nil
}

func (p *Paladin) ValidateAutoCDH(s *skill.Skill) bool {
	return false
}

func (p *Paladin) GetBuffModifiers(encounterTime int64) float64 {
	if p.Buffs.FightOrFlight.AppliedUntil > encounterTime {
		return p.Buffs.FightOrFlight.DamageMod
	}
	return 1
}

func (p *Paladin) GetDesiredPrepullDuration() int {
	return 2000
}

func (p *Paladin) CheckPotionStatus(time int64) bool {
	//return p.Buffs.Tincture.AppliedUntil >= time
	return false
}

func (p *Paladin) IncreaseMP(amount int) {
	p.MP += amount
	if p.MP > 10000 {
		p.MP = 10000
	}
}

func (p *Paladin) DecreaseMP(amount int) {
	p.MP -= amount
	if p.MP < 0 {
		p.MP = 0
	}
}

// Paladin Sills
var (
	// GCD
	FastBlade = skill.Skill{
		Name:        "Fast Blade",
		ID:          9,
		Potency:     200,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&RiotBlade},
	}
	RiotBlade = skill.Skill{
		Name:        "Riot Blade",
		ID:          15,
		Potency:     300,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&RoyalAuthority},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.IncreaseMP(1000)
		},
	}
	RoyalAuthority = skill.Skill{
		Name:        "Royal Authority",
		ID:          3539,
		Potency:     400,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.SwordOath.Stacks = p.Buffs.SwordOath.MaxStacks
			p.Buffs.DivineMight.Stacks = p.Buffs.DivineMight.MaxStacks
		},
	}
	Atonement = skill.Skill{
		Name:        "Atonement",
		ID:          16460,
		Potency:     400,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		NextCombo:   nil,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.SwordOath.Stacks--
			p.IncreaseMP(400)
		},
	}
	GoringBlade = skill.Skill{
		Name:        "Goring Blade",
		ID:          3538,
		Potency:     700,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		NextCombo:   nil,
		CooldownMS:  60000,
		Charges:     1,
		MaxCharges:  1,
	}
	Confiteor = skill.Skill{
		Name:        "Confiteor",
		ID:          16459,
		Potency:     920,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&BladeOfFaith},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.DecreaseMP(1000)
			p.Buffs.Requiescat.Stacks--
			p.Buffs.Confiteor.Stacks = 0
		},
	}
	BladeOfFaith = skill.Skill{
		Name:        "Blade of Faith",
		ID:          25748,
		Potency:     720,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&BladeOfTruth},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.Requiescat.Stacks--
			p.DecreaseMP(1000)
		},
	}
	BladeOfTruth = skill.Skill{
		Name:        "Blade of Truth",
		ID:          25749,
		Potency:     820,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		NextCombo:   []*skill.Skill{&BladeOfValor},
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.Requiescat.Stacks--
			p.DecreaseMP(1000)
		},
	}
	BladeOfValor = skill.Skill{
		Name:        "Blade of Valor",
		ID:          25750,
		Potency:     920,
		GCD:         skill.GCD2_5,
		BreaksCombo: true,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.Requiescat.Stacks--
			p.DecreaseMP(1000)
		},
	}
	HolySpirit = skill.Skill{
		Name:        "Holy Spirit",
		ID:          7384,
		Potency:     350,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.DecreaseMP(1000)
		},
	}
	HolySpiritDivineMight = skill.Skill{
		Name:        "Holy Spirit",
		ID:          7384,
		Potency:     450,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.DivineMight.Stacks--
		},
	}
	HolySpiritRequiescat = skill.Skill{
		Name:        "Holy Spirit",
		ID:          7384,
		Potency:     650,
		GCD:         skill.GCD2_5,
		BreaksCombo: false,
		LockMS:      700,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.Requiescat.Stacks--
		},
	}
	// OGCDs
	FightOrFlight = skill.Skill{
		Name:        "Fight or Flight",
		ID:          20,
		Potency:     0,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.Buffs.FightOrFlight.AppliedUntil = time + p.Buffs.FightOrFlight.DurationMS
		},
	}
	Requiescat = skill.Skill{
		Name:        "Requiescat",
		ID:          7383,
		Potency:     320,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  60000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			// Add 4 Requiescat stacks
			p.Buffs.Requiescat.Stacks = p.Buffs.Requiescat.MaxStacks
			// Add Confiteor Ready
			p.Buffs.Confiteor.Stacks = p.Buffs.Confiteor.MaxStacks
		},
	}
	CircleOfScorn = skill.Skill{
		Name:        "Circle of Scorn",
		ID:          23,
		Potency:     140,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			j := job.(*Job)
			p := j.JobImpl.(*Paladin)
			// Check if Fight or Flight is active
			buffMod := 1.0
			if p.Buffs.FightOrFlight.AppliedUntil > time {
				buffMod = p.Buffs.FightOrFlight.DamageMod
			}
			// Apply the DoT
			target.ApplyDot(&p.DOTs.CircleOfScorn, &j.DamageDealt, time, buffMod)
		},
	}
	Expiacion = skill.Skill{
		Name:        "Expiacion",
		ID:          25747,
		Potency:     450,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  1,
		Charges:     1,
		CustomLogic: func(job any, target *enemy.Enemy, time int64) {
			p := job.(*Job).JobImpl.(*Paladin)
			p.IncreaseMP(500)
		},
	}
	Intervene = skill.Skill{
		Name:        "Intervene",
		ID:          16461,
		Potency:     150,
		GCD:         skill.OGCD,
		BreaksCombo: false,
		LockMS:      700,
		CooldownMS:  30000,
		MaxCharges:  2,
		Charges:     2,
	}
	// DOTs
	CircleOfScornDOT = buff.DOT{
		Name:       "Circle of Scorn",
		ID:         248,
		DurationMS: 15000,
		Potency:    30,
	}
	// Buffs
	SwordOath = buff.Buff{
		Name:      "Sword Oath",
		ID:        1901,
		Stacks:    0,
		MaxStacks: 3,
	}
	FightOrFlightBuff = buff.Buff{
		Name:         "Fight or Flight",
		ID:           76,
		DamageMod:    1.25,
		DurationMS:   20000,
		AppliedUntil: 0,
	}
	RequiescatBuff = buff.Buff{
		Name:      "Requiescat",
		ID:        1368,
		Stacks:    0,
		MaxStacks: 4,
	}
	ConfiteorReady = buff.Buff{
		Name:      "Confiteor Ready",
		ID:        3019,
		Stacks:    0,
		MaxStacks: 1,
	}
	DivineMight = buff.Buff{
		Name:      "Divine Might",
		ID:        3018,
		Stacks:    0,
		MaxStacks: 1,
	}
)
