package wowlogs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		time  time.Time
		event EventVersion
		err   error
	}{
		{
			"Blank",
			"1/1/1 0:0:0.0  COMBAT_LOG_VERSION,0,ADVANCED_LOG_ENABLED,0,BUILD_VERSION,0.0.0,PROJECT_ID,0",
			time.Time{},
			EventVersion{},
			nil,
		},
		{
			"Filled",
			"2/3/4 5:6:7.8  COMBAT_LOG_VERSION,9,ADVANCED_LOG_ENABLED,1,BUILD_VERSION,10.11.12,PROJECT_ID,13",
			time.Date(4, time.February, 3, 5, 6, 7, 8, time.UTC),
			EventVersion{
				Log:      9,
				Major:    10,
				Minor:    11,
				Patch:    12,
				Project:  13,
				Advanced: true,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp, event, err := Parse(tt.input)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.event, event)
				assert.Equal(t, tt.time, timestamp)
			}
		})
	}
}

func TestParseZoneChange(t *testing.T) {
	tests := []struct {
		name  string
		input string
		time  time.Time
		event EventZoneChange
		err   error
	}{
		{
			"Blank",
			"1/1/1 0:0:0.0  ZONE_CHANGE,0,\"\",0",
			time.Time{},
			EventZoneChange{},
			nil,
		},
		{
			"Filled",
			"2/3/4 5:6:7.8  ZONE_CHANGE,9,\"abc\",10",
			time.Date(4, time.February, 3, 5, 6, 7, 8, time.UTC),
			EventZoneChange{
				Instance:   9,
				Zone:       "abc",
				Difficulty: 10,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp, event, err := Parse(tt.input)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.Equal(t, tt.event, event)
				assert.Equal(t, tt.time, timestamp)
			}
		})
	}
}
