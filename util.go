package echo

import "time"

type Sleeper struct {
	Duration time.Duration
	Step     time.Duration
	Max      time.Duration
}

func NewSleeper(step time.Duration, max time.Duration) *Sleeper {
	return &Sleeper{
		Duration: step,
		Step:     step,
		Max:      max,
	}
}

func (sleeper *Sleeper) Sleep() {
	time.Sleep(sleeper.Duration)
	if sleeper.Duration < sleeper.Max {
		sleeper.Duration += sleeper.Step
	}
}

func (sleeper *Sleeper) Reset() {
	sleeper.Duration = sleeper.Step
}
