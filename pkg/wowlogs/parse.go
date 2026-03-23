package wowlogs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var regexpFormat = regexp.MustCompile(`(.+ .+)  ([a-zA-Z0-9_]+),(.+)`)
var regexpTime = regexp.MustCompile(`(\d+)\/(\d+)\/(\d+) (\d+):(\d+):(\d+).(\d+)`)
var regexpVersion = regexp.MustCompile(`(\d+),ADVANCED_LOG_ENABLED,(\d),BUILD_VERSION,(\d+).(\d+).(\d+),PROJECT_ID,(\d+)`)
var regexpZoneChange = regexp.MustCompile(`(\d+),\"(.*)\",(\d+)`)

var eventMap = map[string]func(string) (Event, error){
	"COMBAT_LOG_VERSION": parseVersion,
	"ZONE_CHANGE":        parseZoneChange,
}

func ParseFile(path string) ([]time.Time, []Event, error) {
	var times []time.Time
	var events []Event

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		time, event, err := Parse(scanner.Text())
		if err != nil {
			return nil, nil, err
		}
		times = append(times, time)
		events = append(events, event)
	}

	return times, events, nil
}

func Parse(str string) (time.Time, Event, error) {
	timeStr, name, fields, err := parseFormat(str)
	if err != nil {
		return time.Time{}, nil, err
	}

	timestamp, err := parseTime(timeStr)
	if err != nil {
		return time.Time{}, nil, err
	}

	eventFunc, exists := eventMap[name]
	if !exists {
		return time.Time{}, nil, fmt.Errorf("unknown event %s", name)
	}

	event, err := eventFunc(fields)
	if err != nil {
		return time.Time{}, nil, err
	}

	return timestamp, event, nil
}

func parseFormat(str string) (string, string, string, error) {
	match := regexpFormat.FindStringSubmatch(str)
	if len(match) < 4 {
		return "", "", "", errors.New("invalid format")
	}
	return match[1], match[2], match[3], nil
}

func parseTime(str string) (time.Time, error) {
	match := regexpTime.FindStringSubmatch(str)
	if len(match) < 8 {
		return time.Time{}, errors.New("invalid time format")
	}

	day, err := strconv.Atoi(match[2])
	if err != nil {
		return time.Time{}, errors.New("invalid day")
	}

	month, err := strconv.Atoi(match[1])
	if err != nil {
		return time.Time{}, errors.New("invalid month")
	}

	year, err := strconv.Atoi(match[3])
	if err != nil {
		return time.Time{}, errors.New("invalid year")
	}

	hour, err := strconv.Atoi(match[4])
	if err != nil {
		return time.Time{}, errors.New("invalid hour")
	}

	minute, err := strconv.Atoi(match[5])
	if err != nil {
		return time.Time{}, errors.New("invalid minute")
	}

	second, err := strconv.Atoi(match[6])
	if err != nil {
		return time.Time{}, errors.New("invalid second")
	}

	ms, err := strconv.Atoi(match[7])
	if err != nil {
		return time.Time{}, errors.New("invalid second")
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, ms, time.UTC), nil
}

func parseVersion(str string) (Event, error) {
	match := regexpVersion.FindStringSubmatch(str)
	if len(match) < 6 {
		return nil, errors.New("invalid version format")
	}

	log, err := strconv.Atoi(match[1])
	if err != nil || log < 0 {
		return nil, errors.New("invalid version")
	}

	major, err := strconv.Atoi(match[3])
	if err != nil || log < 0 {
		return nil, errors.New("invalid version")
	}

	minor, err := strconv.Atoi(match[4])
	if err != nil || log < 0 {
		return nil, errors.New("invalid version")
	}

	patch, err := strconv.Atoi(match[5])
	if err != nil || log < 0 {
		return nil, errors.New("invalid version")
	}

	project, err := strconv.Atoi(match[6])
	if err != nil || log < 0 {
		return nil, errors.New("invalid version")
	}

	advanced, err := strconv.ParseBool(match[2])
	if err != nil {
		return nil, errors.New("invalid advanced flag")
	}

	return EventVersion{
		Log:      uint(log),
		Major:    uint(major),
		Minor:    uint(minor),
		Patch:    uint(patch),
		Project:  uint(project),
		Advanced: advanced,
	}, nil
}

func parseZoneChange(str string) (Event, error) {
	match := regexpZoneChange.FindStringSubmatch(str)
	fmt.Println(match)
	if len(match) < 3 {
		return nil, errors.New("invalid zone change format")
	}

	instance, err := strconv.Atoi(match[1])
	if err != nil || instance < 0 {
		return nil, errors.New("invalid instance")
	}

	zone := match[2]

	difficulty, err := strconv.Atoi(match[3])
	if err != nil || difficulty < 0 {
		return nil, errors.New("invalid difficulty")
	}

	return EventZoneChange{
		Instance:   uint(instance),
		Zone:       zone,
		Difficulty: uint(difficulty),
	}, nil
}
