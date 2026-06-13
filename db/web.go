package db

import (
	"ocl/web"
)

func WebDataHeader() (web.DataHeader, error) {
	QueueActive, err := ImportStatusCount("processing")
	if err != nil {
		return web.DataHeader{}, err
	}

	QueueTotal, err := ImportStatusCount("pending")
	if err != nil {
		return web.DataHeader{}, err
	}

	storage, err := Size()
	if err != nil {
		return web.DataHeader{}, err
	}

	StorageGB := float64(storage) / 1024 / 1024 / 1024

	Requests24H, err := RequestCount24H()
	if err != nil {
		return web.DataHeader{}, err
	}

	RequestsTotal, err := RequestCount()
	if err != nil {
		return web.DataHeader{}, err
	}

	Visitors24H, err := RequestIPCount24H()
	if err != nil {
		return web.DataHeader{}, err
	}

	VisitorsTotal, err := RequestIPCount()
	if err != nil {
		return web.DataHeader{}, err
	}

	data := web.DataHeader{
		QueueActive: QueueActive,
		QueueTotal:  QueueTotal,
		StorageGB:   StorageGB,
		Requests24H: Requests24H,
		Requests:    RequestsTotal,
		Visitors24H: Visitors24H,
		Visitors:    VisitorsTotal,
	}

	return data, nil
}

func WebDataIndex() (web.DataIndex, error) {
	header, err := WebDataHeader()
	if err != nil {
		return web.DataIndex{}, err
	}

	return web.DataIndex{
		Header: header,
	}, nil
}

func WebDataEncounters() (web.DataEncounters, error) {
	return web.DataEncounters{}, nil
}

func WebDataEncountersID(id int64) (web.DataEncountersID, error) {
	header, err := WebDataHeader()
	if err != nil {
		return web.DataEncountersID{}, err
	}

	return web.DataEncountersID{
		Header: header,
	}, nil
}

func WebDataLogs() (web.DataLogs, error) {
	return web.DataLogs{}, nil
}

func WebDataLogsID_Entries(id int64) ([]web.DataLogsID_Entry, error) {
	var entries []web.DataLogsID_Entry

	rows, err := db.Query(`
		select event_encounter_start.id, encounters.name
		from event_encounter_start
		left join encounters on encounters.id = event_encounter_start.encounter_id
		where event_encounter_start.import_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var name *string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		entry := web.DataLogsID_Entry{
			Name: name,
			ID:   id,
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func WebDataLogsID(id int64) (web.DataLogsID, error) {
	header, err := WebDataHeader()
	if err != nil {
		return web.DataLogsID{}, err
	}

	entries, err := WebDataLogsID_Entries(id)
	if err != nil {
		return web.DataLogsID{}, err
	}

	return web.DataLogsID{
		Header:  header,
		ID:      id,
		Entries: entries,
	}, nil
}

func WebDataLogsIDEntry() (web.DataLogsIDEntry, error) {
	return web.DataLogsIDEntry{}, nil
}
