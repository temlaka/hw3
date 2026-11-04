package domain

import "time"

type Schedule struct {
	Enabled         bool
	IntervalSeconds int
	LastUpdate      time.Time
	NextUpdate      time.Time
}
