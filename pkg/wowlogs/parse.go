package wowlogs

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"
	"strings"
	"time"
)

func ParseFile(path string) ([]Event, error) {
	var events []Event

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event, err := Parse(scanner.Text())
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func Parse(str string) (Event, error) {
	elems := strings.Split(str, "  ")
	if len(elems) != 2 {
		return nil, errors.New("malformed event")
	}

	ts, err := time.Parse("1/2/2006 15:4:5.0", elems[0])
	if err != nil {
		return nil, err
	}

	return parseEvent(ts, elems[1])
}

func parseEvent(ts time.Time, str string) (Event, error) {
	r := csv.NewReader(strings.NewReader(str))
	fields, err := r.Read()
	if err != nil {
		return nil, err
	}

	if len(fields) < 1 {
		return nil, errors.New("malformed event fields")
	}

	p := fieldParser{fields: fields, pos: 1}
	switch fields[0] {
	case "COMBAT_LOG_VERSION":
		e := EventVersion{
			Time:     ts,
			Log:      p.Int(),
			Advanced: p.BoolSkip(),
			Version:  p.StringSkip(),
			Project:  p.IntSkip(),
		}
		return e, p.err
	case "ZONE_CHANGE":
		return EventZoneChange{
			Time:       ts,
			Instance:   p.Int(),
			Zone:       p.String(),
			Difficulty: p.Int(),
		}, nil
	case "MAP_CHANGE":
		return EventMapChange{
			Time: ts,
			ID:   p.Int(),
			Name: p.String(),
			X0:   p.Float(),
			Y0:   p.Float(),
			X1:   p.Float(),
			Y1:   p.Float(),
		}, nil
	default:
		return nil, errors.New("unknown event type")
	}
}
