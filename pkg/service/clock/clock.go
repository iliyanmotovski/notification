package clock

import "time"

type IClock interface {
	Now() time.Time
}

func NewClockService() IClock {
	return &SysClock{}
}

type SysClock struct{}

func (c *SysClock) Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
