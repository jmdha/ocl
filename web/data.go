package web

import "time"

type DataIndex struct {
}

type DataCharacters_Entry struct {
	ID   int64
	Name string
}

type DataCharacters struct {
	Entries []DataCharacters_Entry
}

type DataCharactersID struct {
}

type DataEncounters struct {
}

type DataEncountersID struct {
	Name string
}

type DataLogs_Entry struct {
	ID   int64
	Date time.Time
}

type DataLogs struct {
	Entries []DataLogs_Entry
}

type DataLogsID_Entry struct {
	Name     *string
	Duration time.Duration
	Success  bool
	ID       int64
}

type DataLogsID struct {
	ID      int64
	Entries []DataLogsID_Entry
}

type DataLogsIDEntry_Entry struct {
}

type DataLogsIDEntry struct {
	ID      int64
	IDEntry int64
	Entries []DataLogsIDEntry_Entry
}
