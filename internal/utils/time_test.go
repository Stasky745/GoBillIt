package utils

import (
	"testing"
	"time"
)

func TestFormatDateString(t *testing.T) {
	d := time.Date(2004, time.April, 29, 0, 0, 0, 0, time.Now().Location())
	tests := []struct {
		date     time.Time
		pattern  string
		expected string
	}{
		// Valid cases
		{d, "DD-M,YYYY", "29-4,2004"},
		{d, "m D, YYYY", "April 29, 2004"},
		{d, "d, DD m YYYY", "Thursday, 29 April 2004"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			got := FormatDate(tt.date, tt.pattern)
			if got != tt.expected {
				t.Errorf("formatDateString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
