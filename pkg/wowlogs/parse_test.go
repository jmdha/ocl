package wowlogs

import (
	"testing"
	"time"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		event EventVersion
		err   error
	}{
		{
			"Blank",
			"1/1/2000 0:0:0.0  COMBAT_LOG_VERSION,0,ADVANCED_LOG_ENABLED,0,BUILD_VERSION,,PROJECT_ID,0",
			EventVersion{
				Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"Filled",
			"2/3/2000 5:6:7.8  COMBAT_LOG_VERSION,9,ADVANCED_LOG_ENABLED,1,BUILD_VERSION,10.11.12,PROJECT_ID,13",
			EventVersion{
				Time:     time.Date(2000, time.February, 3, 5, 6, 7, 800_000_000, time.UTC),
				Log:      9,
				Version:  "10.11.12",
				Project:  13,
				Advanced: true,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := Parse(tt.input)
			if err != tt.err {
				t.Errorf("expected %v found %v", tt.err, err)
			}
			if event != tt.event {
				t.Errorf("expected %v found %v", tt.event, event)
			}
		})
	}
}

func TestParseZoneChange(t *testing.T) {
	tests := []struct {
		name  string
		input string
		event EventZoneChange
		err   error
	}{
		{
			"Blank",
			"1/1/2000 0:0:0.0  ZONE_CHANGE,0,\"\",0",
			EventZoneChange{
				Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"Filled",
			"2/3/2000 5:6:7.8  ZONE_CHANGE,9,\"abc\",10",
			EventZoneChange{
				Time:       time.Date(2000, time.February, 3, 5, 6, 7, 800_000_000, time.UTC),
				Instance:   9,
				Zone:       "abc",
				Difficulty: 10,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := Parse(tt.input)
			if err != tt.err {
				t.Errorf("expected %v found %v", tt.err, err)
			}
			if event != tt.event {
				t.Errorf("expected %v found %v", tt.event, event)
			}
		})
	}
}
