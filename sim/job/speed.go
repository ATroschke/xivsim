package job

import (
	"math"

	"github.com/ATroschke/xivsim/sim/skill"
)

// Speed holds GCDs and Cast Times in ms
type Speed struct {
	AA_DOT_MOD float64
	GCD1_5     int
	GCD2       int
	GCD2_5     int
	GCD2_8     int
	GCD3       int
	GCD3_5     int
	GCD4       int
}

// TODO: Do we need both SkillSpeed and SpellSpeed for casters? Does SkillSpeed scale Auto Attacks?
func NewSpeed(speedAmount, LvSub, LvDiv int) *Speed {
	speed := Speed{}

	speed.calculateSpeed(speedAmount, LvSub, LvDiv)

	return &speed
}

func (s *Speed) calculateSpeed(speed, LvSub, LvDiv int) {
	// AA and DOT Modifier is calculated as follows: f(SPD) = ( 1000 + ⌊ 130 × ( Speed - Level Lv, SUB )/ Level Lv, DIV ⌋ ) / 1000
	s.AA_DOT_MOD = (1000 + math.Floor(130*(float64(speed)-float64(LvSub))/float64(LvDiv))) / 1000.0
	// GCD is calculated as follows: =(INT(GCD*(1000+CEILING(130*(400-Speed)/1900))/10000)/100)
	s.GCD1_5 = int(1500*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
	s.GCD2 = int(2000*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
	s.GCD2_5 = int(2500*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
	s.GCD2_8 = int(2800*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
	s.GCD3 = int(3000*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
	s.GCD3_5 = int(3500*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
	s.GCD4 = int(4000*(1000+math.Ceil(130*(float64(LvSub-speed))/float64(LvDiv)))/10000) * 10
}

func (s *Speed) GetGCD(gcd skill.GCD) int {
	switch gcd {
	case skill.GCD1_5:
		return s.GCD1_5
	case skill.GCD2:
		return s.GCD2
	case skill.GCD2_5:
		return s.GCD2_5
	case skill.GCD2_8:
		return s.GCD2_8
	case skill.GCD3:
		return s.GCD3
	case skill.GCD3_5:
		return s.GCD3_5
	case skill.GCD4:
		return s.GCD4
	default:
		return 0
	}
}
