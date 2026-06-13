package service

import (
	"testing"
	"time"
)

func TestCalculateAgeAfterBirthday(t *testing.T) {
	dob := time.Date(1990, time.May, 10, 0, 0, 0, 0, time.UTC)
	today := time.Date(2026, time.June, 13, 0, 0, 0, 0, time.UTC)

	if got := CalculateAge(dob, today); got != 36 {
		t.Fatalf("CalculateAge() = %d, want 36", got)
	}
}

func TestCalculateAgeBeforeBirthday(t *testing.T) {
	dob := time.Date(1990, time.December, 10, 0, 0, 0, 0, time.UTC)
	today := time.Date(2026, time.June, 13, 0, 0, 0, 0, time.UTC)

	if got := CalculateAge(dob, today); got != 35 {
		t.Fatalf("CalculateAge() = %d, want 35", got)
	}
}

func TestCalculateAgeOnBirthday(t *testing.T) {
	dob := time.Date(1990, time.June, 13, 0, 0, 0, 0, time.UTC)
	today := time.Date(2026, time.June, 13, 0, 0, 0, 0, time.UTC)

	if got := CalculateAge(dob, today); got != 36 {
		t.Fatalf("CalculateAge() = %d, want 36", got)
	}
}
