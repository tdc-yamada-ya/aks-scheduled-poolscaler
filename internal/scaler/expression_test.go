package scaler

import (
	"testing"
	"time"
)

func TestExpressionMatch(t *testing.T) {
	tests := []struct {
		e      Expression
		t      time.Time
		expect bool
	}{
		{"* * * * * * *", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), false},
		{"* * * * * *", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), true},
		{"0 0 1 1 2021 5", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), true},
		{"0 0 1 1 2021 5", time.Date(2021, time.February, 2, 1, 1, 0, 0, time.UTC), false},
		{"0 0 1 Jan 2021 Fri", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), true},
		{"0 0 1 Jan 2021 Fri", time.Date(2021, time.February, 2, 1, 1, 0, 0, time.UTC), false},
		{"* 18-23 * * * *", time.Date(2021, time.January, 1, 17, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), false},
		{"* 18-23 * * * *", time.Date(2021, time.January, 1, 18, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* 18-23 * * * *", time.Date(2021, time.January, 1, 23, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* 0-8 * * * *", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* 0-8 * * * *", time.Date(2021, time.January, 1, 8, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* 0-8 * * * *", time.Date(2021, time.January, 1, 9, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), false},
		{"* * * * * Sun,Sat", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), false},
		{"* * * * * Sun,Sat", time.Date(2021, time.January, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* * * * * Sun,Sat", time.Date(2021, time.January, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* * * * * Sun,Sat", time.Date(2021, time.January, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), false},
		{"* * 1 * * *", time.Date(2021, time.January, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), true},
		{"* * 1 * * *", time.Date(2021, time.January, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), false},
	}

	for _, tt := range tests {
		r := tt.e.Match(tt.t)
		if r != tt.expect {
			t.Errorf("got %v, want %v", r, tt.expect)
		}
	}
}

func TestPatternMatch(t *testing.T) {
	tests := []struct {
		pattern pattern
		n       int
		expect  bool
	}{
		{"*", 0, true},
		{"*", 1, true},
		{"a", 0, false},
		{"10", 10, true},
		{"10", 11, false},
		{"10,11", 10, true},
		{"10,11", 12, false},
		{"10-", 10, true},
		{"10-", 9, false},
		{"-10", 10, true},
		{"-10", 11, false},
		{"10-20", 10, true},
		{"10-20", 9, false},
		{"10-20", 20, true},
		{"10-20", 21, false},
	}

	for _, tt := range tests {
		r := tt.pattern.Match(tt.n, convert)
		if r != tt.expect {
			t.Errorf("got %v, want %v", r, tt.expect)
		}
	}
}
