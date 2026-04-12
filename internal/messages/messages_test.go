package messages

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestParseMessagesFromFile(t *testing.T) {
	// Test nonexistent file
	messages := ParseMessagesFromFile("nonexistent.yaml")
	if len(messages) == 0 {
		t.Errorf("Expected fallback to default messages, got empty map")
	}
	if messages[MessageKeys.NextHoliday] != defaultMessages[MessageKeys.NextHoliday] {
		t.Errorf("Expected fallback to exact default message")
	}

	// Test valid yaml file
	tmpFile, err := os.CreateTemp("", "messages-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	yamlContent := `
nextHoliday: "Custom next holiday {{ .HolidayName }}"
`
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	customMessages := ParseMessagesFromFile(tmpFile.Name())
	if customMessages[MessageKeys.NextHoliday] != "Custom next holiday {{ .HolidayName }}" {
		t.Errorf("Expected custom message from file, got '%s'", customMessages[MessageKeys.NextHoliday])
	}
}

func TestGetMessage(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "messages-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	yamlContent := `
daysLeft: "Testing custom days left"
`
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Setting up the viper config to point to our temp file
	viper.Set("messages-file", tmpFile.Name())
	defer viper.Set("messages-file", "")

	msg := GetMessage(MessageKeys.DaysLeft)
	if msg != "Testing custom days left" {
		t.Errorf("Expected 'Testing custom days left', got '%s'", msg)
	}

	// Test fallback when key doesn't exist in custom file
	fallbackMsg := GetMessage(MessageKeys.NextLargeHoliday)
	if fallbackMsg != defaultMessages[MessageKeys.NextLargeHoliday] {
		t.Errorf("Expected fallback message, got '%s'", fallbackMsg)
	}
}

func TestTemplateMessage(t *testing.T) {
	type TestData struct {
		DaysLeft int
		Username string
	}

	msg := "Hello {{ .Username }}, you have {{ .DaysLeft }} days left!"
	data := TestData{DaysLeft: 5, Username: "Fede"}

	result := TemplateMessage(msg, data)
	expected := "Hello Fede, you have 5 days left!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test faulty template
	faultyMsg := "Hello {{ .Nonexistent"
	resultFaulty := TemplateMessage(faultyMsg, data)
	if resultFaulty != MessageKeys.FailedToParseHolidayDate {
		t.Errorf("Expected error message return placeholder on bad parse, got '%s'", resultFaulty)
	}
}
