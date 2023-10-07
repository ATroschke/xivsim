package models

type AppliedBuff struct {
	SourceID     int
	AppliedUntil int64
	Buff         Buff
}

type ApplyType int

const (
	// ApplyTypeSelf is a buff that is applied to the player and overwrites any other buff with the same ID
	ApplyTypeSelf = iota
	// ApplyTypeSelfExtend is a buff that is applied to the player and extends any other buff with the same ID
	ApplyTypeSelfExtend
	// ApplyTypeEnemy is a debuff that is applied to the enemy and overwrites any other debuff with the same ID
	ApplyTypeEnemy
	// ApplyTypeGroup is a debuff that is applied to the group and overwrites any other debuff with the same ID
	ApplyTypeGroup
)

type Buff struct {
	ID        int
	Duration  int64
	ApplyType ApplyType
	DamageMod float64
}
