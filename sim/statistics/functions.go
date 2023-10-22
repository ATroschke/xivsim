package statistics

import (
	"math"
)

type XivSimStatistics struct {
	Damage []int
}

func (s *XivSimStatistics) AddDamage(damage int) {
	s.Damage = append(s.Damage, damage)
}

func (s *XivSimStatistics) Mean() float64 {
	sum := 0
	for _, d := range s.Damage {
		sum += d
	}
	return float64(sum) / float64(len(s.Damage))
}

func (s *XivSimStatistics) StandardDeviation() float64 {
	mean := s.Mean()
	sum := 0.0
	for _, d := range s.Damage {
		sum += math.Pow(float64(d)-mean, 2)
	}
	return math.Sqrt(sum / float64(len(s.Damage)-1))
}

func (s *XivSimStatistics) StanardError() float64 {
	sdev := s.StandardDeviation()
	return sdev / math.Sqrt(float64(len(s.Damage)))
}

// CriticalValue is the value from the t-distribution table
// for the given confidence level and degrees of freedom
// 95% confidence level -> 1.96
// 99% confidence level -> 2.58
func (s *XivSimStatistics) MarginOfError(CriticalValue float64) float64 {
	return s.StanardError() * CriticalValue
}

func (s *XivSimStatistics) ConfidenceInterval(CriticalValue float64) (float64, float64) {
	mean := s.Mean()
	marginOfError := s.MarginOfError(CriticalValue)
	return mean - marginOfError, mean + marginOfError
}
