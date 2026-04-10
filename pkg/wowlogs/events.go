package wowlogs

import "time"

type Event any // why is there not a union type?

type EventVersion struct {
	Time     time.Time
	Log      int
	Version  string
	Project  int
	Advanced bool
}

type EventZoneChange struct {
	Time       time.Time
	Instance   int
	Zone       string
	Difficulty int
}

type EventMapChange struct {
	Time time.Time
	ID   int
	Name string
	X0   float64
	Y0   float64
	X1   float64
	Y1   float64
}
