package buff

type Buff struct {
	Name          string
	ID            int
	DurationMS    int64   // Duration in milliseconds
	MaxDurationMS int64   // Max duration in milliseconds
	AppliedUntil  int64   // the time this buff is applied until
	MaxStacks     int     // the maximum number of stacks this buff can have
	Stacks        int     // the current number of stacks this buff has
	DamageMod     float64 // the damage modifier of this buff (1.0 = 100%, default)
}
