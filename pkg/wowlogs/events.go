package wowlogs

type Event any // why is there not a union type?

type EventVersion struct {
	Log      uint
	Major    uint
	Minor    uint
	Patch    uint
	Project  uint
	Advanced bool
}

type EventZoneChange struct {
	Instance   uint
	Zone       string
	Difficulty uint
}

type EventMapChange struct {
	ID             uint
	Name           string
	X0, Y0, X1, Y1 float64
}
