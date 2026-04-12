package helpers

import (
	"testing"
	"time"
)

func TestFormatDateToSpanishUnparsed(t *testing.T) {
	formatted, day, month, year := FormatDateToSpanishUnparsed("2023-05-15")
	if formatted != "Lunes, 15 de mayo" {
		t.Errorf("Expected 'Lunes, 15 de mayo', got '%s'", formatted)
	}
	if day != "lunes" || month != "mayo" || year != "2023" {
		t.Errorf("Unexpected returns from unparsed: %s, %s, %s", day, month, year)
	}

	// invalid date
	f, d, m, y := FormatDateToSpanishUnparsed("not-a-date")
	if f != "" || d != "" || m != "" || y != "" {
		t.Errorf("Expected empty returns on invalid date")
	}
}

func TestFormatDateToSpanish(t *testing.T) {
	// Let's test 1st of month to trigger "1ro"
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	formatted, day, month, year := FormatDateToSpanish(date)
	
	if formatted != "Domingo, 1ro de enero" {
		t.Errorf("Expected 'Domingo, 1ro de enero', got '%s'", formatted)
	}
	if day != "domingo" || month != "enero" || year != "2023" {
		t.Errorf("Unexpected parts: %s, %s, %s", day, month, year)
	}
}

func TestMonthsToSpanish(t *testing.T) {
	if MonthsToSpanish(1) != "enero" {
		t.Errorf("Expected enero")
	}
	if MonthsToSpanish(12) != "diciembre" {
		t.Errorf("Expected diciembre")
	}
}

func TestFormatDate(t *testing.T) {
	res := FormatDate("2023-12-25")
	if res != "Lunes 25 de Diciembre" {
		t.Errorf("Expected 'Lunes 25 de Diciembre', got '%s'", res)
	}

	resErr := FormatDate("invalid")
	if resErr != "invalid" {
		t.Errorf("Expected fallback 'invalid'")
	}
}
