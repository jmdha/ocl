package db

import (
	"ocl/web"
	"time"
)

func WebDataHeader() (web.DataHeader, error) {
	return web.DataHeader{}, nil
}

func WebDataFooter() (web.DataFooter, error) {
	QueueActive, err := ImportStatusCount("processing")
	if err != nil {
		return web.DataFooter{}, err
	}

	QueueTotal, err := ImportStatusCount("pending")
	if err != nil {
		return web.DataFooter{}, err
	}

	storage, err := Size()
	if err != nil {
		return web.DataFooter{}, err
	}

	StorageGB := float64(storage) / 1024 / 1024 / 1024

	Requests24H, err := RequestCount24H()
	if err != nil {
		return web.DataFooter{}, err
	}

	RequestsTotal, err := RequestCount()
	if err != nil {
		return web.DataFooter{}, err
	}

	Visitors24H, err := RequestIPCount24H()
	if err != nil {
		return web.DataFooter{}, err
	}

	VisitorsTotal, err := RequestIPCount()
	if err != nil {
		return web.DataFooter{}, err
	}

	data := web.DataFooter{
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

	footer, err := WebDataFooter()
	if err != nil {
		return web.DataIndex{}, err
	}

	return web.DataIndex{
		Header: header,
		Footer: footer,
	}, nil
}

func WebDataCharacters_Entries() ([]web.DataCharacters_Entry, error) {
	var entries []web.DataCharacters_Entry

	rows, err := db.Query(`
		select id, name from character
		where guid like 'player-%'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		entry := web.DataCharacters_Entry{
			ID:   id,
			Name: name,
		}

		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return entries, nil
}

func WebDataCharacters() (web.DataCharacters, error) {
	entries, err := WebDataCharacters_Entries()
	if err != nil {
		return web.DataCharacters{}, err
	}
	return web.DataCharacters{
		Entries: entries,
	}, nil
}

func WebDataCharactersID() (web.DataCharactersID, error) {
	return web.DataCharactersID{}, nil
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

func WebDataLogs_Entries(limit int64, offset int64) ([]web.DataLogs_Entry, error) {
	var entries []web.DataLogs_Entry

	rows, err := db.Query(`
		select id, timestamp from import
		where status not like 'failed'
		order by timestamp desc
		limit ? offset ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var date int64

		err = rows.Scan(&id, &date)
		if err != nil {
			return nil, err
		}

		entry := web.DataLogs_Entry{
			ID:   id,
			Date: time.Unix(0, date),
		}

		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return entries, nil
}

func WebDataLogs(limit int64, offset int64) (web.DataLogs, error) {
	header, err := WebDataHeader()
	if err != nil {
		return web.DataLogs{}, err
	}

	footer, err := WebDataFooter()
	if err != nil {
		return web.DataLogs{}, err
	}

	entries, err := WebDataLogs_Entries(limit, offset)
	if err != nil {
		return web.DataLogs{}, err
	}

	return web.DataLogs{
		Header:  header,
		Footer:  footer,
		Entries: entries,
	}, nil
}

func WebDataLogsID_Entries(id int64) ([]web.DataLogsID_Entry, error) {
	var entries []web.DataLogsID_Entry

	rows, err := db.Query(`
		select id, encounter_name, (timestamp_end - timestamp_start), success
		from encounter
		where import_id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name *string
		var duration int64
		var success int

		err = rows.Scan(&id, &name, &duration, &success)
		if err != nil {
			return nil, err
		}

		entry := web.DataLogsID_Entry{
			Name:     name,
			ID:       id,
			Duration: time.Duration(duration),
			Success:  success == 1,
		}

		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return nil, err
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
