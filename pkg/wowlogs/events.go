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
