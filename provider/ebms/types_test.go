package ebms

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOperationTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    OperationTime
		expectedErr string
	}{{
		name:  "OnlyAM",
		value: "07:10 - 09:30",
		expected: OperationTime{
			From: time.Date(0, 0, 0, 7, 10, 0, 0, time.Local),
			To:   time.Date(0, 0, 0, 9, 30, 0, 0, time.Local),
		},
	}, {
		name:  "OnlyPM",
		value: "13:10 - 20:59",
		expected: OperationTime{
			From: time.Date(0, 0, 0, 13, 10, 0, 0, time.Local),
			To:   time.Date(0, 0, 0, 20, 59, 0, 0, time.Local),
		},
	}, {
		name:  "MixedAMPM",
		value: "01:10 - 20:59",
		expected: OperationTime{
			From: time.Date(0, 0, 0, 01, 10, 0, 0, time.Local),
			To:   time.Date(0, 0, 0, 20, 59, 0, 0, time.Local),
		},
	}, {
		name:        "InvalidHour",
		value:       "25:10 - 20:59",
		expectedErr: "parsing time \"25:10\": hour out of range",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var operationTime OperationTime
			err := json.Unmarshal(fmt.Appendf(nil, `"%s"`, test.value), &operationTime)

			if test.expectedErr == "" {
				assert.NoError(t, err)

				assert.Equal(t, test.expected.From.Hour(), operationTime.From.Hour())
				assert.Equal(t, test.expected.From.Minute(), operationTime.From.Minute())
				assert.Equal(t, test.expected.To.Hour(), operationTime.To.Hour())
				assert.Equal(t, test.expected.To.Minute(), operationTime.To.Minute())
			} else {
				assert.EqualError(t, err, test.expectedErr)
			}
		})
	}
}

func TestOperationTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		value    OperationTime
		expected string
	}{{
		name: "OnlyAM",
		value: OperationTime{
			From: time.Date(0, 0, 0, 7, 10, 0, 0, time.Local),
			To:   time.Date(0, 0, 0, 9, 30, 0, 0, time.Local),
		},
		expected: `"07:10 - 09:30"`,
	}, {
		name: "OnlyPM",
		value: OperationTime{
			From: time.Date(0, 0, 0, 13, 10, 0, 0, time.Local),
			To:   time.Date(0, 0, 0, 20, 59, 0, 0, time.Local),
		},
		expected: `"13:10 - 20:59"`,
	}, {
		name: "MixedAMPM",
		value: OperationTime{
			From: time.Date(0, 0, 0, 01, 10, 0, 0, time.Local),
			To:   time.Date(0, 0, 0, 20, 59, 0, 0, time.Local),
		},
		expected: `"01:10 - 20:59"`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := json.Marshal(test.value)
			assert.NoError(t, err)
			assert.Equal(t, []byte(test.expected), value)
		})
	}
}

func TestAppliedDates_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    AppliedDates
		expectedErr string
	}{{
		name:  "Weekdays",
		value: "T2, T3, T4, T5, T6",
		expected: AppliedDates{
			Dates: []AppliedDate{
				{Weekday: time.Monday},
				{Weekday: time.Tuesday},
				{Weekday: time.Wednesday},
				{Weekday: time.Thursday},
				{Weekday: time.Friday},
			},
		},
	}, {
		name:  "Single",
		value: "T2",
		expected: AppliedDates{
			Dates: []AppliedDate{
				{Weekday: time.Monday},
			},
		},
	}, {
		name:        "Invalid",
		value:       "T9",
		expectedErr: "invalid weekday",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var appliedDates AppliedDates
			err := json.Unmarshal(fmt.Appendf(nil, `"%s"`, test.value), &appliedDates)

			if test.expectedErr == "" {
				assert.NoError(t, err)

				assert.Equal(t, test.expected, appliedDates)
			} else {
				assert.EqualError(t, err, test.expectedErr)
			}
		})
	}
}
