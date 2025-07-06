package holidays

import (
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/FGasquez/alum-bot/internal/types"
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
	var skipWeekend bool = true

	params := helpers.GetParams(i.ApplicationCommandData().Options)
	if _, ok := params["skip-today"]; ok {
		skipToday = params["skip-today"].(bool)
	}
	if _, ok := params["skip-weekend"]; ok {
		skipWeekend = params["skip-weekend"].(bool)
	}

	daysLeftToHoliday, nextHoliday, isToday := DaysLeft(skipWeekend, skipToday)

	logrus.Info("##################### skipWeekend: ", skipWeekend)

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
	dateFormatted, day, month, _ := helpers.FormatDateToSpanish(parsedDate)

	tmpValues := types.TemplateValues{
		HolidayName:   nextHoliday.Name,
		DaysLeft:      daysLeftToHoliday,
		FormattedDate: dateFormatted,
		NamedDate: types.NamedDate{
			Day:   day,
			Month: month,
		},
		RawDate: types.RawDate{
			Day:   parsedDate.Day(),
			Month: int(parsedDate.Month()),
			Year:  parsedDate.Year(),
		},
		FullDate:  nextHoliday.Date,
		Adjacents: nextHoliday.Adjacent,
		IsToday:   isToday,
	}

	message := messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.NextHoliday), tmpValues)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
