package ports

import "time"

// Clock abstracts time Source to make testing deterministic.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using system time.
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}
