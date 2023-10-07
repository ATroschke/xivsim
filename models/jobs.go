package models

// Job is an interface that represents a job in the game
type Job interface {
	GetName() string
	SelectNextSkill(player *Player, enemy *Enemy, encounterTime int64) *Skill
}
