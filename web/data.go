package web

import "time"

type DataHeader struct {
}

type DataFooter struct {
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
	Footer DataFooter
}

type DataCharacters_Entry struct {
	ID   int64
	Name string
}

type DataCharacters struct {
	Header  DataHeader
	Footer  DataFooter
	Entries []DataCharacters_Entry
}

type DataCharactersID struct {
	Header DataHeader
	Footer DataFooter
}

type DataEncounters struct {
	Header DataHeader
	Footer DataFooter
}

type DataEncountersID struct {
	Header DataHeader
	Footer DataFooter
	Name   string
}

type DataLogs_Entry struct {
	ID   int64
	Date time.Time
}

type DataLogs struct {
	Header  DataHeader
	Footer  DataFooter
	Entries []DataLogs_Entry
}

type DataLogsID_Entry struct {
	Name     *string
	Duration time.Duration
	Success  bool
	ID       int64
}

type DataLogsID struct {
	Header  DataHeader
	Footer  DataFooter
	ID      int64
	Entries []DataLogsID_Entry
}

type DataLogsIDEntry struct {
	Header DataHeader
	Footer DataFooter
}
