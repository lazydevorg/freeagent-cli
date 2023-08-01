package timeslip

import (
	"testing"
	"time"
)

func Test_weekRange(t *testing.T) {
	tests := []struct {
		name string
		date string
		from string
		to   string
	}{
		{"Monday", "2023-8-1", "2023-7-31", "2023-8-6"},
		{"Tuesday", "2023-8-2", "2023-7-31", "2023-8-6"},
		{"Wednesday", "2023-8-3", "2023-7-31", "2023-8-6"},
		{"Thursday", "2023-8-4", "2023-7-31", "2023-8-6"},
		{"Friday", "2023-8-5", "2023-7-31", "2023-8-6"},
		{"Saturday", "2023-8-6", "2023-7-31", "2023-8-6"},
		{"Sunday", "2023-8-7", "2023-8-7", "2023-8-13"},
		{"Previous week", "2023-7-30", "2023-7-24", "2023-7-30"},
		{"Next week", "2023-8-7", "2023-8-7", "2023-8-13"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := time.Parse("2006-1-2", tt.date)
			if err != nil {
				t.Errorf("weekRange() error: %v", err)
			}
			from, to := weekRange(day)
			from_string := from.Format("2006-1-2")
			to_string := to.Format("2006-1-2")
			if from_string != tt.from {
				t.Errorf("weekRange(%s) from got = %v, want %v", tt.date, from, tt.from)
			}
			if to_string != tt.to {
				t.Errorf("weekRange(%s) to got = %v, want %v", tt.date, to, tt.to)
			}
		})
	}
}
