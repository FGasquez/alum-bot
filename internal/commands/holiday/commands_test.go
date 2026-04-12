package holidays

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestHowManyDaysToHolidayFlags(t *testing.T) {
	expectedFlags := map[string]discordgo.ApplicationCommandOptionType{
		"skip-today":   discordgo.ApplicationCommandOptionBoolean,
		"skip-weekend": discordgo.ApplicationCommandOptionBoolean,
	}

	foundFlags := make(map[string]bool)

	for _, opt := range HowManyDaysToHoliday.Options {
		if expectedType, ok := expectedFlags[opt.Name]; ok {
			if opt.Type != expectedType {
				t.Errorf("Expected flag %s to be of type %v, got %v", opt.Name, expectedType, opt.Type)
			}
			foundFlags[opt.Name] = true
		}
	}

	for flag := range expectedFlags {
		if !foundFlags[flag] {
			t.Errorf("Expected flag %s not found in HowManyDaysToHoliday command", flag)
		}
	}
}

func TestHolidaysCommandsFlags(t *testing.T) {
	expectedFlags := map[string]discordgo.ApplicationCommandOptionType{
		"skip-today":   discordgo.ApplicationCommandOptionBoolean,
		"skip-weekend": discordgo.ApplicationCommandOptionBoolean,
	}

	foundFlags := make(map[string]bool)

	for _, opt := range HolidaysCommands.Options {
		if expectedType, ok := expectedFlags[opt.Name]; ok {
			if opt.Type != expectedType {
				t.Errorf("Expected flag %s to be of type %v, got %v", opt.Name, expectedType, opt.Type)
			}
			foundFlags[opt.Name] = true
		}
	}

	for flag := range expectedFlags {
		if !foundFlags[flag] {
			t.Errorf("Expected flag %s not found in HolidaysCommands command", flag)
		}
	}
}

func TestHolidaysOfMonthFlags(t *testing.T) {
	expectedFlags := map[string]discordgo.ApplicationCommandOptionType{
		"month": discordgo.ApplicationCommandOptionInteger,
	}

	foundFlags := make(map[string]bool)

	for _, opt := range HolidaysOfMonth.Options {
		if expectedType, ok := expectedFlags[opt.Name]; ok {
			if opt.Type != expectedType {
				t.Errorf("Expected flag %s to be of type %v, got %v", opt.Name, expectedType, opt.Type)
			}
			foundFlags[opt.Name] = true
		}
        
        if opt.Name == "month" {
            if !opt.Required {
                t.Errorf("Expected flag month to be required")
            }
            if len(opt.Choices) != 12 {
                t.Errorf("Expected month choice to have 12 months, got %d", len(opt.Choices))
            }
        }
	}

	for flag := range expectedFlags {
		if !foundFlags[flag] {
			t.Errorf("Expected flag %s not found in HolidaysOfMonth command", flag)
		}
	}
}
