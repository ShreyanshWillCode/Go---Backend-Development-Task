package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name	string
		dob	time.Time
		now	time.Time
		wantAge	int
	}{
		{
			name:		"exactly 35 years — birthday already passed this year",
			dob:		time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			now:		time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC),
			wantAge:	35,
		},
		{
			name:		"birthday is today — should still count as a full year",
			dob:		time.Date(1990, 6, 12, 0, 0, 0, 0, time.UTC),
			now:		time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC),
			wantAge:	35,
		},
		{
			name:		"birthday is tomorrow — not yet 35",
			dob:		time.Date(1990, 6, 13, 0, 0, 0, 0, time.UTC),
			now:		time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC),
			wantAge:	34,
		},
		{
			name:		"leap year birthday (Feb 29) — measured on Feb 28",
			dob:		time.Date(1992, 2, 29, 0, 0, 0, 0, time.UTC),
			now:		time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			wantAge:	32,
		},
		{
			name:		"very young user — born last year",
			dob:		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			now:		time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC),
			wantAge:	1,
		},
		{
			name:		"newborn — same day",
			dob:		time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC),
			now:		time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC),
			wantAge:	0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := CalculateAge(tc.dob, tc.now)
			if got != tc.wantAge {
				t.Errorf("CalculateAge(%v, %v) = %d; want %d", tc.dob, tc.now, got, tc.wantAge)
			}
		})
	}
}
