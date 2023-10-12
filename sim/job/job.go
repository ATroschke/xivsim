package job

import (
	"github.com/ATroschke/xivsim/sim/enemy"
	xivmath "github.com/ATroschke/xivsim/xiv-math"
)

type Job interface {
	CalculateSkills(
		weaponDamage int,
		mainStat int,
		criticalHit int,
		directHit int,
		determination int,
		skillSpeed int,
		spellSpeed int,
		tenacity int,
	)
	//AddBuff()
	Tick(enemy *enemy.Enemy, encounterTime int64) (int, int)
	Report()
}

type DamageModifiers struct {
	CritRate   float64
	CritDamage float64
	DirectRate float64
}

func CalculateDamageModifiers(criticalHit int, directHit int, LvSub int, LvDiv int) DamageModifiers {
	return DamageModifiers{
		CritRate:   xivmath.P_CHR(criticalHit, LvSub, LvDiv),
		CritDamage: xivmath.F_CRIT(criticalHit, LvSub, LvDiv),
		DirectRate: xivmath.P_DHR(directHit, LvSub, LvDiv),
	}
}
