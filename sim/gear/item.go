package gear

const (
	SLOT_MAINHAND = iota
	SLOT_OFFHAND
	SLOT_HEAD
	SLOT_BODY
	SLOT_HANDS
	SLOT_WAIST // Lol...
	SLOT_LEGS
	SLOT_FEET
	SLOT_EARRINGS
	SLOT_NECKLACE
	SLOT_BRACELET
	SLOT_RING1
	SLOT_RING2
	SLOT_FOOD
	SLOT_POTION
)

type Item struct {
	Name       string
	ItemLevel  int
	Slot       int
	Crafted    bool
	Materia    []Materia
	Stats      ItemStats
	MaxSubstat int // Max Substat Value for the Item Level, currently just taken from the higher value of the two possible substats, used to limit materia melds
}

type ItemStats struct {
	WeaponDamage  int
	Strength      int
	Dexterity     int
	Intelligence  int
	Mind          int
	CriticalHit   int
	Determination int
	DirectHit     int
	SkillSpeed    int
	SpellSpeed    int
	Tenacity      int
	Piety         int
}
