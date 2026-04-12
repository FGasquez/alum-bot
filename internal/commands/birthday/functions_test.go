package birthday

import (
	"os"
	"testing"
	"time"

	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/spf13/viper"
)

func TestGetNextBirthday(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// We calculate strings dynamically from "today" so tests don't break dynamically
	futureDate := today.AddDate(0, 0, 2).Format("01-02")
	pastDate := today.AddDate(0, 0, -5).Format("01-02")

	birthdays := []types.Birthday{
		{Name: "Bob (Passed)", Username: "bobby", Date: pastDate},
		{Name: "Alice (Future)", Username: "<@123>", Date: futureDate},
	}

	next, err := GetNextBirthday(birthdays)
	if err != nil {
		t.Fatalf("Failed getting next birthday: %v", err)
	}

	if next.Name != "Alice (Future)" {
		t.Errorf("Expected Alice to be the next birthday, got %s", next.Name)
	}

	if next.DaysLeft != 2 {
		t.Errorf("Expected 2 days left, got %d", next.DaysLeft)
	}

	// Add a birthday that triggers exactly next year
	expectedNextYearObj := today.AddDate(1, 0, 0)
	
	birthdays = []types.Birthday{
		{Name: "AlreadyPassedThisYear", Username: "u1", Date: pastDate},
	}

	nextOnlyOne, err := GetNextBirthday(birthdays)
	
	if err != nil {
		t.Fatalf("Failed getting next single instance: %v", err)
	}

	// We don't assert days precisely here because leap years modify exact Sub(today).Hours()
	// But we assert the NextOccurrence accurately rolled over a year:
	if nextOnlyOne.NextOccurrence.Year() != expectedNextYearObj.Year() {
		t.Errorf("Expected next occurrence year to be %d, got %d", expectedNextYearObj.Year(), nextOnlyOne.NextOccurrence.Year())
	}
}

func TestGetDaysLeft(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	futureDate := today.AddDate(0, 0, 10).Format("01-02")
	pastDate := today.AddDate(0, 0, -10).Format("01-02")

	birthdays := []types.Birthday{
		{Name: "Bob (Passed)", Username: "bobby", Date: pastDate},
		{Name: "Alice (Future)", Username: "<@123>", Date: futureDate},
	}

	// Test user with future birthday this year
	parsedFuture, err := GetDaysLeft("<@123>", birthdays)
	if err != nil {
		t.Fatalf("Failed getting days left for <@123>: %v", err)
	}

	if parsedFuture.Name != "Alice (Future)" {
		t.Errorf("Expected Alice, got %s", parsedFuture.Name)
	}

	if parsedFuture.DaysLeft != 10 {
		t.Errorf("Expected 10 days left, got %d", parsedFuture.DaysLeft)
	}

	// Test user with past birthday this year (should roll over to next year)
	parsedPast, err := GetDaysLeft("bobby", birthdays)
	if err != nil {
		t.Fatalf("Failed getting days left for bobby: %v", err)
	}

	if parsedPast.Name != "Bob (Passed)" {
		t.Errorf("Expected Bob, got %s", parsedPast.Name)
	}

	expectedNextYearObj := today.AddDate(1, 0, 0)
	if parsedPast.NextOccurrence.Year() != expectedNextYearObj.Year() {
		t.Errorf("Expected next occurrence year to be %d, got %d", expectedNextYearObj.Year(), parsedPast.NextOccurrence.Year())
	}

	// Test non-existent user
	_, err = GetDaysLeft("<@404>", birthdays)
	if err == nil {
		t.Errorf("Expected error for non-existent user, got nil")
	}
}

func TestGetBirthdays(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "birthdays-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	yamlContent := `
- name: Bob
  username: bobby
  date: 01-01
- name: Alice
  date: 02-02
`
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	viper.Set("birthdays-file", tmpFile.Name())
	defer viper.Set("birthdays-file", "")

	birthdays, err := GetBirthdays()
	if err != nil {
		t.Fatalf("Expected nil err, got %v", err)
	}

	if len(birthdays) != 2 {
		t.Fatalf("Expected 2 birthdays, got %d", len(birthdays))
	}
	if birthdays[0].Name != "Bob" || birthdays[0].Username != "bobby" || birthdays[0].Date != "01-01" {
		t.Errorf("Unexpected first birthday parse: %+v", birthdays[0])
	}
	if birthdays[1].Name != "Alice" || birthdays[1].Username != "" || birthdays[1].Date != "02-02" {
		t.Errorf("Unexpected second birthday parse: %+v", birthdays[1])
	}
}
