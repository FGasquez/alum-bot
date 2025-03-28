package holidays

import (
	"fmt"
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const HolidaysCommandName = "next-holiday"

var HolidaysCommands = discordgo.ApplicationCommand{
	Name:        HolidaysCommandName,
	Description: "Get the next holiday",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "skip-today",
			Description: "skip today in the calculation",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "skip-weekend",
			Description: "skip weekend in the calculation",
			Required:    false,
		},
	},
}

var HolidaysCommandHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var skipToday bool = false
	var skipWeekend bool = false

	params := helpers.GetParams(i.ApplicationCommandData().Options)
	if _, ok := params["skip-today"]; ok {
		skipToday = params["skip-today"].(bool)
	}
	if _, ok := params["skip-weekend"]; ok {
		skipWeekend = params["skip-weekend"].(bool)
	}

	nextHoliday, isToday := NextHoliday(time.Now(), skipWeekend, skipToday)
	logrus.Infof("Next holiday: %s, date: %s", nextHoliday.Name, nextHoliday.Date)
	if isToday {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Es hoy! 🎉",
			},
		})
		return
	}

	// Parse nextHoliday.Date into a time.Time object
	parsedDate, err := time.Parse("2006-01-02", nextHoliday.Date)
	if err != nil {
		logrus.Errorf("Failed to parse holiday date: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve the next holiday. Please try again later.",
			},
		})
		return
	}

	day, month, year := helpers.FormatDateToSpanish(parsedDate)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🎉 El próximo feriado es **%s** el **%s %s %s**. 🎉", nextHoliday.Name, day, month, year),
		},
	})
}
