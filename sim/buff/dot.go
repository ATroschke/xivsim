package buff

import xivmath "github.com/ATroschke/xivsim/xiv-math"

type DOT struct {
	Name             string
	Source           *int // Points to the DamageDealt field of the source, so we can attribute the Damage to the source
	ID               int
	DurationMS       int64
	AppliedUntil     int64
	Potency          int
	Damage           int
	CritDamage       int
	DirectDamage     int
	CritDirectDamage int
	CritRate         float64
	DirectRate       float64
	BuffMod          float64
	Magical          bool
}

func (d *DOT) CalculateDamage(
	weaponDamage int,
	mainStat int,
	criticalHit int,
	directHit int,
	determination int,
	skillSpeed int,
	spellSpeed int,
	tenacity int,
	weaponDelay float64,
	m int,
	jobmod int,
) {
	if d.Magical {
		d.Damage = xivmath.MagicalDOTDamage(d.Potency, weaponDamage, m, mainStat, 390, jobmod, determination, 1900, tenacity, skillSpeed, criticalHit, 400, 100)
	} else {
		d.Damage = xivmath.PhysicalDOTDamage(d.Potency, weaponDamage, m, mainStat, 390, jobmod, determination, 1900, tenacity, skillSpeed, criticalHit, 400, 100)
	}
}
