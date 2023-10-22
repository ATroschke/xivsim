package buff

type Debuff struct {
	Name         string
	Source       *int // Points to the DamageDealt field of the source, so we can attribute the Damage to the source
	ID           int
	AppliedUntil int64
	DamageMod    float64
}
