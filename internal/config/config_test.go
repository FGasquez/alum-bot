package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigGetters(t *testing.T) {
	// Setup test values via Viper
	viper.Set("token", "test-token")
	viper.Set("test-guilds", []string{"guild1", "guild2"})
	viper.Set("messages-file", "messages.yaml")
	viper.Set("birthdays-file", "birthdays.yaml")
	viper.Set("advertise-today-birthday", true)
	viper.Set("advertise-channel", "123456")

	// Cleanup at the end
	defer func() {
		viper.Set("token", "")
		viper.Set("test-guilds", nil)
		viper.Set("messages-file", "")
		viper.Set("birthdays-file", "")
		viper.Set("advertise-today-birthday", false)
		viper.Set("advertise-channel", "")
	}()

	if GetToken() != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", GetToken())
	}

	guilds := GetTestGuilds()
	if !reflect.DeepEqual(guilds, []string{"guild1", "guild2"}) {
		t.Errorf("Expected guilds [guild1, guild2], got %v", guilds)
	}

	if GetMessagesPath() != "messages.yaml" {
		t.Errorf("Expected 'messages.yaml', got '%s'", GetMessagesPath())
	}

	if GetBirthdaysPath() != "birthdays.yaml" {
		t.Errorf("Expected 'birthdays.yaml', got '%s'", GetBirthdaysPath())
	}

	if !GetAdvertiseTodayBirthday() {
		t.Errorf("Expected advertise-today-birthday to be true")
	}

	if GetAdvertiseChannel() != "123456" {
		t.Errorf("Expected advertise-channel '123456', got '%s'", GetAdvertiseChannel())
	}
}

func TestGetTimezone(t *testing.T) {
	// Test default behavior (empty string)
	viper.Set("timezone", "")
	loc := GetTimezone()
	if loc.String() != "America/Argentina/Buenos_Aires" {
		t.Errorf("Expected default timezone 'America/Argentina/Buenos_Aires', got '%s'", loc.String())
	}

	// Test valid given timezone
	viper.Set("timezone", "Europe/London")
	loc = GetTimezone()
	if loc.String() != "Europe/London" {
		t.Errorf("Expected timezone 'Europe/London', got '%s'", loc.String())
	}

	// Test invalid timezone fallback (UTC)
	viper.Set("timezone", "Invalid/Timezone")
	loc = GetTimezone()
	if loc.String() != "UTC" {
		t.Errorf("Expected fallback timezone 'UTC', got '%s'", loc.String())
	}

	// Cleanup
	viper.Set("timezone", "")
}
