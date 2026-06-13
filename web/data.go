package web

import "time"

type DataHeader struct {
	QueueActive int64
	QueueTotal  int64
	StorageGB   float64
	Requests    int64
	Requests24H int64
	Visitors    int64
	Visitors24H int64
}

type DataIndex struct {
	Header DataHeader
}

type DataEncounters struct {
	Header DataHeader
}

type DataEncountersID struct {
	Header DataHeader
	Name   string
}

type DataLogs struct {
	Header DataHeader
}

type DataLogsID_Entry struct {
	Name     *string
	Duration time.Duration
	Success  bool
	ID       int64
}

type DataLogsID struct {
	Header  DataHeader
	ID      int64
	Entries []DataLogsID_Entry
}

type DataLogsIDEntry struct {
	Header DataHeader
}
