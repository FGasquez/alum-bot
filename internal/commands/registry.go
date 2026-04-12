package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	birthdayCmd "github.com/FGasquez/alum-bot/internal/commands/birthday"
	holidaysCmd "github.com/FGasquez/alum-bot/internal/commands/holiday"
	"github.com/FGasquez/alum-bot/internal/config"
)

// BotCommand encapsulates a Discord application command, its handler, and an optional condition for registration.
type BotCommand struct {
	AppCommand *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
	IsEnabled  func() bool
}

// allCommands lists all available commands in the bot and their requirements.
var allCommands = []BotCommand{
	{
		AppCommand: &holidaysCmd.HolidaysCommands,
		Handler:    holidaysCmd.HolidaysCommandHandlers,
		IsEnabled:  nil, // Always enabled
	},
	{
		AppCommand: &holidaysCmd.HowManyDaysToHoliday,
		Handler:    holidaysCmd.HowManyDaysToHolidayHandlers,
		IsEnabled:  nil,
	},
	{
		AppCommand: &holidaysCmd.HolidaysOfMonth,
		Handler:    holidaysCmd.HolidaysOfMonthHandlers,
		IsEnabled:  nil,
	},
	{
		AppCommand: &holidaysCmd.HolidaysLargeCommands,
		Handler:    holidaysCmd.HolidayLargeCommandHandlers,
		IsEnabled:  nil,
	},
	{
		AppCommand: &birthdayCmd.NextBdayCommand,
		Handler:    birthdayCmd.NextBdayCommandHandler,
		IsEnabled: func() bool {
			return config.GetBirthdaysPath() != ""
		},
	},
	{
		AppCommand: &birthdayCmd.BdayLeftCommand,
		Handler:    birthdayCmd.BdayLeftCommandHandler,
		IsEnabled: func() bool {
			return config.GetBirthdaysPath() != ""
		},
	},
}

// GetActiveCommands evaluates each command's requirements and returns the ones that should be registered.
func GetActiveCommands() ([]*discordgo.ApplicationCommand, map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	activeCommands := []*discordgo.ApplicationCommand{}
	handlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

	for _, cmd := range allCommands {
		if cmd.IsEnabled != nil && !cmd.IsEnabled() {
			logrus.Infof("Command %s is disabled by configuration, skipping registration.", cmd.AppCommand.Name)
			continue
		}
		activeCommands = append(activeCommands, cmd.AppCommand)
		handlers[cmd.AppCommand.Name] = cmd.Handler
	}

	return activeCommands, handlers
}
