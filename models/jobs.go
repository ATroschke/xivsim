package models

// Job is an interface that represents a job in the game
type Job interface {
	GetName() string
	UpdateCharges(encounterTime int64)
	SelectNextGCD(player *Player, enemy *Enemy, encounterTime int64) *Skill
	SelectNextOGCD(player *Player, enemy *Enemy, encounterTime int64, GCDUntil int64) *Skill
}

type Charge struct {
	SkillID       int
	SkillCooldown int
	Charges       int
	MaxCharges    int
	NextCharge    int64
}
