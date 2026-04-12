package birthday

import (
	"strings"
	"testing"
	"time"

	"github.com/FGasquez/alum-bot/internal/types"
)

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{123456789, "123456789"},
		{-5, "-5"},
		{-100, "-100"},
	}

	for _, test := range tests {
		result := itoa(test.input)
		if result != test.expected {
			t.Errorf("itoa(%d) expected %s but got %s", test.input, test.expected, result)
		}
	}
}

func TestNewBirthdayTemplateData(t *testing.T) {
	parsed := &ParsedBirthday{
		Birthday: types.Birthday{
			Name:     "Test Name",
			Username: "TestUser",
			Date:     "01-01",
		},
		NextOccurrence: time.Now(),
		DaysLeft:       10,
	}

	data := NewBirthdayTemplateData(parsed, nil, "")
	if data.Name != "Test Name" {
		t.Errorf("Expected Name 'Test Name', got '%s'", data.Name)
	}
	if data.Username != "TestUser" {
		t.Errorf("Expected Username 'TestUser', got '%s'", data.Username)
	}
	if data.Date != "01-01" {
		t.Errorf("Expected Date '01-01', got '%s'", data.Date)
	}
	if data.DaysLeft != 10 {
		t.Errorf("Expected DaysLeft 10, got %d", data.DaysLeft)
	}

	parsedEmptyUser := &ParsedBirthday{
		Birthday: types.Birthday{
			Name: "No Username",
		},
	}
	dataEmpty := NewBirthdayTemplateData(parsedEmptyUser, nil, "")
	if dataEmpty.Username != "**No Username**" {
		t.Errorf("Expected Username '**No Username**', got '%s'", dataEmpty.Username)
	}
}

func TestFormatNextBirthdayMessage(t *testing.T) {
	parsed10Days := &ParsedBirthday{
		Birthday: types.Birthday{
			Name:     "Alice",
			Username: "alice",
			Date:     "01-01",
		},
		DaysLeft: 10,
	}

	msg10 := FormatNextBirthdayMessage(parsed10Days, nil, "")
	if !strings.Contains(msg10, "(10 days left)") {
		t.Errorf("Expected message to contain '(10 days left)', got '%s'", msg10)
	}
	if !strings.Contains(msg10, "alice") {
		t.Errorf("Expected message to contain 'alice', got '%s'", msg10)
	}

	parsedTomorrow := &ParsedBirthday{
		Birthday: types.Birthday{
			Name:     "Bob",
			Username: "bob",
			Date:     "01-02",
		},
		DaysLeft: 1,
	}
	msg1 := FormatNextBirthdayMessage(parsedTomorrow, nil, "")
	if !strings.Contains(msg1, "(It's tomorrow! 🎂)") {
		t.Errorf("Expected message to contain tomorrow marker, got '%s'", msg1)
	}

	parsedToday := &ParsedBirthday{
		Birthday: types.Birthday{
			Name:     "Charlie",
			Username: "charlie",
			Date:     "01-03",
		},
		DaysLeft: 0,
	}
	msg0 := FormatNextBirthdayMessage(parsedToday, nil, "")
	if !strings.Contains(msg0, "Today is charlie's birthday") && !strings.Contains(msg0, "charlie") {
		t.Errorf("Expected today message, got '%s'", msg0)
	}
}

func TestFormatTodayBirthdayMessage(t *testing.T) {
	parsed := &ParsedBirthday{
		Birthday: types.Birthday{
			Name:     "Dave",
			Username: "dave",
			Date:     "01-04",
		},
	}

	msg := FormatTodayBirthdayMessage(parsed, nil, "")
	if !strings.Contains(msg, "dave") {
		t.Errorf("Expected today message to contain 'dave', got '%s'", msg)
	}
}
